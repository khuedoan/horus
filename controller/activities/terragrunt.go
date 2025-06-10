package activities

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

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

func TerragruntGraphShaking(ctx context.Context, graph *Graph, changedFiles []string) (*Graph, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("Pruning Terragrunt DAG graph")

	prunedGraph, err := PruneGraph(ctx, graph, changedFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to prune dependency graph: %w", err)
	}

	return prunedGraph, nil
}

func TerragruntApply(ctx context.Context, repoPath string, modulePath string, stack string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Running terragrunt apply", "module", modulePath, "stack", stack)

	// Construct the full path to the module
	fullPath := filepath.Join(repoPath, "infra", stack, modulePath)

	cmd := exec.CommandContext(ctx, "terragrunt", "apply", "--backend-bootstrap", "--auto-approve")
	cmd.Dir = fullPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run terragrunt apply for module %s: %w\nOutput: %s", modulePath, err, string(output))
	}

	logger.Info("Terragrunt apply completed", "module", modulePath, "output", string(output))
	return nil
}
