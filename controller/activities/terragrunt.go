package activities

import (
	"context"
	"fmt"
	"os/exec"

	"go.temporal.io/sdk/activity"
)

func TerragruntGraph(ctx context.Context, path string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Generating Terragrunt DAG graph", "path", path)

	cmd := exec.CommandContext(ctx, "terragrunt", "dag", "graph")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run terragrunt dag graph: %w", err)
	}

	return string(output), nil
}

func TerragruntGraphShaking(ctx context.Context, dotGraph string, changedFiles []string) (string, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("Parsing Terragrunt DAG graph")

	graph, err := NewGraphFromDot(dotGraph)
	if err != nil {
		return "", fmt.Errorf("failed to parse dot graph: %w", err)
	}

	prunedGraph, err := PruneGraph(ctx, graph, changedFiles)
	if err != nil {
		return "", fmt.Errorf("failed to prune dependency graph: %w", err)
	}

	return prunedGraph.ToDot(), nil
}
