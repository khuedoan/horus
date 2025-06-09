package workflows

import (
	"fmt"
	"strings"
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

// TODO create trigger
// For now do that manually on the UI
// Task queue: cloudlab
// Workflow: Infra
// Input json/plain:
//
//	{
//	  "url": "https://github.com/khuedoan/cloudlab",
//	  "revision": "infra-rewrite",
//	  "oldRevision": "7796870a3c17105d7a13c5b6c990fa895de64952",
//	  "stack": "local"
//	}
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
	err = workflow.ExecuteActivity(ctx, activities.TerragruntGraphShaking, graph, changedModules).Get(ctx, &prunedGraph)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return nil, err
	}

	logger.Info("Infra workflow completed graph pruning.", "nodes", len(prunedGraph.Nodes), "edges", len(prunedGraph.Edges))

	// Get dependency levels for parallel execution
	dependencyLevels := prunedGraph.TopologicalSort()

	// Execute terragrunt apply for each level in dependency order
	for levelIndex, level := range dependencyLevels {
		logger.Info("Starting terragrunt apply for dependency level", "level", levelIndex, "modules", level)

		// Create futures for parallel execution within this level
		var futures []workflow.Future
		for _, moduleName := range level {
			// Create activity options with custom activity ID that includes module name
			// Replace slashes with hyphens for ActivityID compatibility
			safeModuleName := strings.ReplaceAll(moduleName, "/", "-")
			moduleActivityOptions := workflow.ActivityOptions{
				StartToCloseTimeout: 10 * time.Minute,
				ActivityID:          fmt.Sprintf("TerragruntApply-%s", safeModuleName),
			}
			moduleCtx := workflow.WithActivityOptions(ctx, moduleActivityOptions)

			future := workflow.ExecuteActivity(moduleCtx, activities.TerragruntApply, path, moduleName, input.Stack)
			futures = append(futures, future)
		}

		// Wait for all modules in this level to complete
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

	logger.Info("Infra workflow completed successfully.", "totalLevels", len(dependencyLevels), "appliedModules", len(prunedGraph.Nodes))

	return prunedGraph, nil
}
