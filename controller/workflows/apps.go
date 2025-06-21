package workflows

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloudlab/controller/activities"

	"go.temporal.io/sdk/workflow"
)

type AppsInput struct {
	Url      string
	Revision string
	Registry string
	Cluster  string
}

func Apps(ctx workflow.Context, input PlatformInput) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Platform workflow started", "platform", input)

	var workspace string
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 1 * time.Minute,
		}),
		activities.Clone,
		input.Url,
		input.Revision,
	).Get(ctx, &workspace); err != nil {
		return err
	}

	defer os.RemoveAll(workspace)

	appsDir := workspace + "/apps"

	// TODO this should be a separate activity
	var matchedPaths []string
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 10 * time.Second,
		}),
		activities.DiscoverApps,
		appsDir,
		input.Cluster,
	).Get(ctx, &matchedPaths); err != nil {
		return err
	}
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 1 * time.Minute,
	})

	var futures []workflow.Future
	var results []activities.PushResult

	for _, yamlPath := range matchedPaths {
		parts := strings.Split(filepath.ToSlash(yamlPath), "/")
		if len(parts) < 4 {
			logger.Warn("Skipping invalid path", "path", yamlPath)
			continue
		}

		namespace := parts[len(parts)-3]
		app := parts[len(parts)-2]

		logger.Info("Dispatching PushRenderedHelm", "path", yamlPath, "namespace", namespace, "app", app)

		fut := workflow.ExecuteActivity(
			ctx,
			activities.PushRenderedApp,
			appsDir,
			namespace,
			app,
			input.Cluster,
			input.Registry,
		)
		futures = append(futures, fut)
	}

	for _, fut := range futures {
		var result activities.PushResult
		if err := fut.Get(ctx, &result); err != nil {
			return err
		}
		results = append(results, result)
	}

	logger.Info("Finished pushing all matching apps", "count", len(results))
	return nil
}
