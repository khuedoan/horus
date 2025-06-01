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
		Depth:         1,
	})
	if err != nil {
		return "", err
	}

	return path, nil
}
