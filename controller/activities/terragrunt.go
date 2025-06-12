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

func TerragruntPrune(ctx context.Context, graph *Graph, changedFiles []string) (*Graph, error) {
	return PruneGraph(ctx, graph, changedFiles)
}

func TerragruntApply(ctx context.Context, repoUrl string, revision string, modulePath string, stack string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Running terragrunt apply", "module", modulePath, "stack", stack, "repo", repoUrl, "revision", revision)

	// Ensure repository is available (clone if necessary)
	repoPath, err := Clone(ctx, repoUrl, revision)
	if err != nil {
		return fmt.Errorf("failed to ensure repository is available: %w", err)
	}

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
