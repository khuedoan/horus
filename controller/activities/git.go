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

	cmd := exec.CommandContext(ctx, "git", "clone", "--branch", revision, url, path)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return path, nil
}

func changedFiles(ctx context.Context, path string, oldRevision string) ([]string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Getting changed files", "path", path, "oldRevision", oldRevision)

	cmd := exec.CommandContext(ctx, "git", "diff", "--name-only", oldRevision, "HEAD")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

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

	changedFiles, err := changedFiles(ctx, repoPath, oldRevision)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})
	var modules []string

	for _, file := range changedFiles {
		dir := filepath.Dir(file)

		currentDir := dir
		for {
			terragruntPath := filepath.Join(repoPath, currentDir, "terragrunt.hcl")
			if _, err := os.Stat(terragruntPath); err == nil {
				modulePath := currentDir

				if strings.HasPrefix(modulePath, "infra/") {
					parts := strings.Split(filepath.ToSlash(modulePath), "/")
					if len(parts) >= 3 && parts[0] == "infra" {
						modulePath = strings.Join(parts[2:], "/")
					}
				}

				if modulePath != "" && modulePath != "." {
					modulePath = filepath.ToSlash(modulePath)

					if _, exists := seen[modulePath]; !exists {
						modules = append(modules, modulePath)
						seen[modulePath] = struct{}{}
					}
				}
				break
			}

			parent := filepath.Dir(currentDir)
			if parent == currentDir || parent == "." {
				break
			}
			currentDir = parent
		}
	}

	return modules, nil
}
