package activities

import (
	"context"
	"os"

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

func ChangedFiles(ctx context.Context, path string, oldRevision string) ([]string, error) {
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
