package workflows

import (
	"fmt"
	"time"

	"cloudlab/controller/activities"

	"go.temporal.io/sdk/workflow"
)

type InfraInputs struct {
	Url         string
	Revision    string
	OldRevision string
	Stack       string
}

func Infra(ctx workflow.Context, input InfraInputs) (*activities.Graph, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("Infra workflow started", "infra", input)

	var path string
	err := workflow.ExecuteActivity(ctx, activities.Clone, input.Url, input.Revision).Get(ctx, &path)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return nil, err
	}

	var (
		graph          *activities.Graph
		changedModules []string
	)

	graphFuture := workflow.ExecuteActivity(ctx, activities.TerragruntGraph, path+"/infra/"+input.Stack)
	changedModulesFuture := workflow.ExecuteActivity(ctx, activities.ChangedModules, path, input.OldRevision)

	err = graphFuture.Get(ctx, &graph)
	if err != nil {
		logger.Error("TerragruntGraph failed", "Error", err)
		return nil, err
	}

	err = changedModulesFuture.Get(ctx, &changedModules)
	if err != nil {
		logger.Error("ChangedModules failed", "Error", err)
		return nil, err
	}

	var prunedGraph *activities.Graph
	err = workflow.ExecuteActivity(ctx, activities.TerragruntPrune, graph, changedModules).Get(ctx, &prunedGraph)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return nil, err
	}

	logger.Info("Infra workflow completed graph pruning.", "nodes", prunedGraph.NodeCount(), "edges", prunedGraph.EdgeCount())

	dependencyLevels := prunedGraph.TopologicalSort()

	for levelIndex, level := range dependencyLevels {
		logger.Info("Starting terragrunt apply for dependency level", "level", levelIndex, "modules", level)

		var futures []workflow.Future
		for _, moduleName := range level {
			moduleActivityOptions := workflow.ActivityOptions{
				StartToCloseTimeout: 10 * time.Minute,
				Summary:             fmt.Sprintf("%s/%s", input.Stack, moduleName),
			}
			moduleCtx := workflow.WithActivityOptions(ctx, moduleActivityOptions)

			future := workflow.ExecuteActivity(moduleCtx, activities.TerragruntApply, path, moduleName, input.Stack)
			futures = append(futures, future)
		}

		for i, future := range futures {
			err := future.Get(ctx, nil)
			if err != nil {
				logger.Error("TerragruntApply failed", "module", level[i], "level", levelIndex, "Error", err)
				return nil, err
			}
			logger.Info("Module apply completed", "module", level[i], "level", levelIndex)
		}

		logger.Info("Completed terragrunt apply for dependency level", "level", levelIndex, "modules", level)
	}

	logger.Info("Infra workflow completed successfully.", "totalLevels", len(dependencyLevels), "appliedModules", prunedGraph.NodeCount())

	return prunedGraph, nil
}
