package activities

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"go.temporal.io/sdk/activity"
)

func TerragruntGraph(ctx context.Context, path string) (*Graph, error) {
	safeHeartbeat(ctx, "Generating terragrunt dependency graph")

	cmd := exec.CommandContext(ctx, "terragrunt", "dag", "graph")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run terragrunt dag graph: %w", err)
	}

	safeHeartbeat(ctx, "Parsing dependency graph")
	return NewGraphFromDot(string(output))
}

func TerragruntPrune(ctx context.Context, graph *Graph, changedFiles []string) (*Graph, error) {
	return PruneGraph(ctx, graph, changedFiles)
}

func TerragruntApply(ctx context.Context, repoUrl string, revision string, modulePath string, stack string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Running terragrunt apply", "module", modulePath, "stack", stack)

	safeHeartbeat(ctx, fmt.Sprintf("Ensuring repository availability for %s", modulePath))

	repoPath, err := Clone(ctx, repoUrl, revision)
	if err != nil {
		return fmt.Errorf("failed to ensure repository is available: %w", err)
	}

	fullPath := filepath.Join(repoPath, "infra", stack, modulePath)
	safeHeartbeat(ctx, fmt.Sprintf("Starting terragrunt apply for %s", modulePath))

	cmd := exec.CommandContext(ctx, "terragrunt", "apply", "--backend-bootstrap", "--auto-approve")
	cmd.Dir = fullPath

	// Create pipes to capture output and send heartbeats
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start terragrunt apply: %w", err)
	}

	// Monitor output and send heartbeats
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// Send heartbeats while monitoring output
	heartbeatTicker := time.NewTicker(25 * time.Second) // Send heartbeat every 25s (before 30s timeout)
	defer heartbeatTicker.Stop()

	var lastOutput string
	outputScanner := bufio.NewScanner(stdout)
	errorScanner := bufio.NewScanner(stderr)

	for {
		select {
		case err := <-done:
			if err != nil {
				return fmt.Errorf("terragrunt apply failed for module %s: %w", modulePath, err)
			}
			safeHeartbeat(ctx, fmt.Sprintf("Terragrunt apply completed for %s", modulePath))
			return nil

		case <-heartbeatTicker.C:
			safeHeartbeat(ctx, fmt.Sprintf("Terragrunt apply in progress for %s - %s", modulePath, lastOutput))

		default:
			// Check for new output
			if outputScanner.Scan() {
				line := strings.TrimSpace(outputScanner.Text())
				if line != "" {
					lastOutput = line
					logger.Info("Terragrunt output", "module", modulePath, "output", line)
				}
			}
			if errorScanner.Scan() {
				line := strings.TrimSpace(errorScanner.Text())
				if line != "" {
					lastOutput = line
					logger.Info("Terragrunt error output", "module", modulePath, "error", line)
				}
			}

			// Small sleep to prevent busy waiting
			time.Sleep(100 * time.Millisecond)
		}
	}
}
