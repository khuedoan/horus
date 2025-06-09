package activities

import (
	"context"
	"fmt"
	"os/exec"

	"go.temporal.io/sdk/activity"
)

func TerragruntGraph(ctx context.Context, path string) (*Graph, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Generating Terragrunt DAG graph", "path", path)

	cmd := exec.CommandContext(ctx, "terragrunt", "dag", "graph")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run terragrunt dag graph: %w", err)
	}

	graph, err := NewGraphFromDot(string(output))
	if err != nil {
		return nil, fmt.Errorf("failed to parse terragrunt graph output: %w", err)
	}

	return graph, nil
}

func TerragruntGraphShaking(ctx context.Context, graph *Graph, changedModules []string) (*Graph, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("Pruning Terragrunt DAG graph")

	prunedGraph, err := PruneGraph(ctx, graph, changedModules)
	if err != nil {
		return nil, fmt.Errorf("failed to prune dependency graph: %w", err)
	}

	return prunedGraph, nil
}
