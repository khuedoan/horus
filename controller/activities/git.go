package activities

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go.temporal.io/sdk/activity"
)

func Clone(ctx context.Context, url string, revision string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Cloning", "url", url, "revision", revision)

	path, err := os.MkdirTemp("", "infra-")
	if err != nil {
		return "", err
	}

	// Clone the repository
	cmd := exec.CommandContext(ctx, "git", "clone", "--branch", revision, url, path)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return path, nil
}

func changedFiles(ctx context.Context, path string, oldRevision string) ([]string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Getting changed files", "path", path, "oldRevision", oldRevision)

	// Get changed files using git diff --name-only
	cmd := exec.CommandContext(ctx, "git", "diff", "--name-only", oldRevision, "HEAD")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Split output by newlines and filter empty lines
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var files []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, line)
		}
	}

	return files, nil
}

func ChangedModules(ctx context.Context, repoPath string, oldRevision string) ([]string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Getting changed modules", "path", repoPath, "oldRevision", oldRevision)

	// Get all changed files
	changedFiles, err := changedFiles(ctx, repoPath, oldRevision)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})
	modules := make([]string, 0)

	for _, file := range changedFiles {
		// Get the directory of the changed file
		dir := filepath.Dir(file)

		// Walk up the directory tree to find the closest directory containing terragrunt.hcl
		currentDir := dir
		for {
			terragruntPath := filepath.Join(repoPath, currentDir, "terragrunt.hcl")
			if _, err := os.Stat(terragruntPath); err == nil {
				// Found terragrunt.hcl, this is a module directory
				modulePath := currentDir

				// Remove infra/<env>/ prefix if present
				if strings.HasPrefix(modulePath, "infra/") {
					parts := strings.Split(filepath.ToSlash(modulePath), "/")
					if len(parts) >= 3 && parts[0] == "infra" {
						// Remove "infra" and environment (e.g., "dev", "prod")
						modulePath = strings.Join(parts[2:], "/")
					}
				}

				// Skip empty paths
				if modulePath != "" && modulePath != "." {
					// Normalize path separators to forward slashes
					modulePath = filepath.ToSlash(modulePath)

					if _, exists := seen[modulePath]; !exists {
						modules = append(modules, modulePath)
						seen[modulePath] = struct{}{}
					}
				}
				break
			}

			// Move up one directory level
			parent := filepath.Dir(currentDir)
			if parent == currentDir || parent == "." {
				// Reached the root, no terragrunt.hcl found
				break
			}
			currentDir = parent
		}
	}

	return modules, nil
}
