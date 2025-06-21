package workflows

import (
	"fmt"
	"time"

	"cloudlab/controller/activities"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type InfraInputs struct {
	Url         string
	Revision    string
	OldRevision string
	Stack       string
}

func Infra(ctx workflow.Context, input InfraInputs) (*activities.Graph, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Infra workflow started", "infra", input)

	// Clone activity: 30s timeout, quick retry on worker failure
	cloneCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 1 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	})

	var path string
	if err := workflow.ExecuteActivity(cloneCtx, activities.Clone, input.Url, input.Revision).Get(ctx, &path); err != nil {
		return nil, err
	}

	// Graph and analysis activities: moderate timeout
	analysisCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	})

	var graph *activities.Graph
	var prunedGraph *activities.Graph

	// Get the terragrunt graph
	if err := workflow.ExecuteActivity(analysisCtx, activities.TerragruntGraph, path+"/infra/"+input.Stack).Get(ctx, &graph); err != nil {
		return nil, err
	}

	// If oldRevision is not provided, use the full graph (no pruning)
	if input.OldRevision == "" {
		logger.Info("No oldRevision provided, using full graph", "nodes", len(graph.Nodes))
		prunedGraph = graph
	} else {
		// Determine changed modules and prune graph
		var changedModules []string
		if err := workflow.ExecuteActivity(analysisCtx, activities.ChangedModules, path, input.OldRevision).Get(ctx, &changedModules); err != nil {
			return nil, err
		}

		if err := workflow.ExecuteActivity(analysisCtx, activities.PruneGraph, graph, changedModules).Get(ctx, &prunedGraph); err != nil {
			return nil, err
		}

		logger.Info("Graph pruning completed", "nodes", len(prunedGraph.Nodes))
	}

	for levelIndex, level := range prunedGraph.TopologicalSort() {
		logger.Info("Starting terragrunt apply", "level", levelIndex, "modules", level)

		var futures []workflow.Future
		for _, module := range level {
			moduleCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
				StartToCloseTimeout: 30 * time.Minute,
				HeartbeatTimeout:    2 * time.Minute,
				Summary:             fmt.Sprintf("%s/%s", input.Stack, module),
				RetryPolicy: &temporal.RetryPolicy{
					MaximumAttempts: 2,
					NonRetryableErrorTypes: []string{
						"TerraformValidationError",
						"TerraformPlanError",
					},
				},
			})
			futures = append(futures, workflow.ExecuteActivity(moduleCtx, activities.TerragruntApply, input.Url, input.Revision, module, input.Stack))
		}

		for i, future := range futures {
			if err := future.Get(ctx, nil); err != nil {
				logger.Error("TerragruntApply failed", "module", level[i], "level", levelIndex, "error", err)
				return nil, err
			}
			logger.Info("Module apply completed", "module", level[i], "level", levelIndex)
		}
	}

	logger.Info("Infra workflow completed", "levels", len(prunedGraph.TopologicalSort()), "modules", len(prunedGraph.Nodes))
	return prunedGraph, nil
}
