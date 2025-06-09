package activities

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"go.temporal.io/sdk/activity"
)

func Clone(ctx context.Context, url string, revision string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Cloning", "url", url)

	path, err := os.MkdirTemp("", "infra-")
	if err != nil {
		return "", err
	}

	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.ReferenceName(revision),
	})
	if err != nil {
		return "", err
	}

	return path, nil
}

func changedFiles(ctx context.Context, path string, oldRevision string) ([]string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Getting changed files", "path", path, "oldRevision", oldRevision)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	// Resolve old revision
	oldHash, err := repo.ResolveRevision(plumbing.Revision(oldRevision))
	if err != nil {
		return nil, err
	}
	oldCommit, err := repo.CommitObject(*oldHash)
	if err != nil {
		return nil, err
	}

	// Resolve HEAD
	headRef, err := repo.Head()
	if err != nil {
		return nil, err
	}
	newCommit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return nil, err
	}

	// Get trees and diff
	oldTree, err := oldCommit.Tree()
	if err != nil {
		return nil, err
	}
	newTree, err := newCommit.Tree()
	if err != nil {
		return nil, err
	}

	changes, err := oldTree.Diff(newTree)
	if err != nil {
		return nil, err
	}

	// Like `git diff --name-only`
	seen := map[string]struct{}{}
	var files []string
	for _, change := range changes {
		if change.From.Name != "" {
			if _, ok := seen[change.From.Name]; !ok {
				files = append(files, change.From.Name)
				seen[change.From.Name] = struct{}{}
			}
		}
		if change.To.Name != "" {
			if _, ok := seen[change.To.Name]; !ok {
				files = append(files, change.To.Name)
				seen[change.To.Name] = struct{}{}
			}
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
