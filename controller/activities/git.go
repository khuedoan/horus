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

// generateRepoPath creates a deterministic path for the repository based on URL and revision
func generateRepoPath(url string, revision string) string {
	// Create a hash of the URL and revision for a deterministic path
	hash := sha256.Sum256([]byte(url + ":" + revision))
	hashStr := fmt.Sprintf("%x", hash)[:16] // Use first 16 chars of hash

	// Use /tmp/cloudlab-repos/ as base directory
	return filepath.Join("/tmp", "cloudlab-repos", hashStr)
}

// checkRepoStatus checks if repository exists and returns the current commit hash
func checkRepoStatus(ctx context.Context, path string, revision string) (exists bool, currentHash string) {
	// Check if .git directory exists
	gitDir := filepath.Join(path, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return false, ""
	}

	// Get current commit hash
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "HEAD")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return false, ""
	}

	currentHash = strings.TrimSpace(string(output))
	return true, currentHash
}

// isCommitAvailable checks if the desired commit/revision is available in the repository
func isCommitAvailable(ctx context.Context, path string, revision string) bool {
	// Try to resolve the revision to a commit hash
	cmd := exec.CommandContext(ctx, "git", "rev-parse", revision)
	cmd.Dir = path
	_, err := cmd.Output()
	return err == nil
}

func Clone(ctx context.Context, url string, revision string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Ensuring repository availability", "url", url, "revision", revision)

	// Create deterministic path based on URL and revision to enable reuse
	path := generateRepoPath(url, revision)

	// Check if repository already exists and has the correct revision
	if repoExists, currentHash := checkRepoStatus(ctx, path, revision); repoExists {
		if currentHash == revision || isCommitAvailable(ctx, path, revision) {
			logger.Info("Repository already available with correct revision", "path", path)
			return path, nil
		}
		logger.Info("Repository exists but wrong revision, will update", "path", path, "current", currentHash, "desired", revision)
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Remove existing directory if it exists but is inconsistent
	if _, err := os.Stat(path); err == nil {
		logger.Info("Removing existing inconsistent repository", "path", path)
		if err := os.RemoveAll(path); err != nil {
			return "", fmt.Errorf("failed to remove existing repository: %w", err)
		}
	}

	// Clone the repository
	logger.Info("Cloning repository", "url", url, "revision", revision, "path", path)
	cmd := exec.CommandContext(ctx, "git", "clone", "--branch", revision, url, path)
	if err := cmd.Run(); err != nil {
		// Clean up the directory if clone fails
		os.RemoveAll(path)
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}

	logger.Info("Successfully cloned repository", "path", path)
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
