package activities

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go.temporal.io/sdk/activity"
)

func generateRepoPath(url string, revision string) string {
	hash := sha256.Sum256([]byte(url + ":" + revision))
	return filepath.Join("/tmp", "cloudlab-repos", fmt.Sprintf("%x", hash)[:16])
}

func hasCorrectRevision(ctx context.Context, path, revision string) bool {
	if _, err := os.Stat(filepath.Join(path, ".git")); os.IsNotExist(err) {
		return false
	}

	cmd := exec.CommandContext(ctx, "git", "rev-parse", revision)
	cmd.Dir = path
	return cmd.Run() == nil
}

func Clone(ctx context.Context, url string, revision string) (string, error) {
	logger := activity.GetLogger(ctx)
	path := generateRepoPath(url, revision)

	if hasCorrectRevision(ctx, path, revision) {
		logger.Info("Repository already available", "path", path)
		return path, nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", fmt.Errorf("failed to create parent directory: %w", err)
	}
	os.RemoveAll(path)

	logger.Info("Cloning repository", "url", url, "revision", revision)

	cmd := exec.CommandContext(ctx, "git", "clone", "--depth", "1", "--branch", revision, url, path)
	if err := cmd.Run(); err != nil {
		os.RemoveAll(path)
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}

	return path, nil
}

func ChangedModules(ctx context.Context, repoPath string, oldRevision string) ([]string, error) {
	logger := activity.GetLogger(ctx)

	// Since we now clone with depth 1, we need to fetch the oldRevision before we can diff against it
	logger.Info("Fetching old revision for comparison", "oldRevision", oldRevision)
	fetchCmd := exec.CommandContext(ctx, "git", "fetch", "origin", oldRevision)
	fetchCmd.Dir = repoPath
	if err := fetchCmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to fetch old revision %s: %w", oldRevision, err)
	}

	cmd := exec.CommandContext(ctx, "git", "diff", "--name-only", oldRevision, "HEAD")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run git diff: %w", err)
	}

	seen := make(map[string]struct{})
	var modules []string

	for _, file := range strings.Fields(string(output)) {
		if file == "" {
			continue
		}

		for dir := filepath.Dir(file); dir != "." && dir != "/"; dir = filepath.Dir(dir) {
			if _, err := os.Stat(filepath.Join(repoPath, dir, "terragrunt.hcl")); err == nil {
				// Remove infra/stack prefix to get module path
				if parts := strings.Split(filepath.ToSlash(dir), "/"); len(parts) >= 3 && parts[0] == "infra" {
					if module := strings.Join(parts[2:], "/"); module != "" {
						if _, exists := seen[module]; !exists {
							modules = append(modules, module)
							seen[module] = struct{}{}
						}
					}
				}
				break
			}
		}
	}

	return modules, nil
}

func GitAdd(ctx context.Context, path string) error {
	logger := activity.GetLogger(ctx)

	dir := filepath.Dir(path)
	relPath := filepath.Base(path)

	cmd := exec.Command("git", "-C", dir, "add", relPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Error("git add failed", "error", err)
		return fmt.Errorf("git add failed: %w", err)
	}

	return nil
}

func GitCommit(ctx context.Context, dir string, message string) error {
	logger := activity.GetLogger(ctx)

	cmd := exec.Command("git", "-C", dir, "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Error("git commit failed", "error", err)
		return fmt.Errorf("git commit failed: %w", err)
	}

	return nil
}

func GitPush(ctx context.Context, dir string) error {
	logger := activity.GetLogger(ctx)

	cmd := exec.Command("git", "-C", dir, "push")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Error("git push failed", "error", err)
		return fmt.Errorf("git push failed: %w", err)
	}

	return nil
}
