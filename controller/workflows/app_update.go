package workflows

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"cloudlab/controller/activities"

	"go.temporal.io/sdk/workflow"
)

type AppUpdateInput struct {
	Url       string
	Revision  string
	Namespace string
	App       string
	Cluster   string
	NewImages []activities.Image
}

// AppUpdate workflow clones a repository, updates app versions, and syncs changes back to git
func AppUpdate(ctx workflow.Context, input AppUpdateInput) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("AppUpdate workflow started", "input", input)

	var workspace string
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 2 * time.Minute,
		}),
		activities.Clone,
		input.Url,
		input.Revision,
	).Get(ctx, &workspace); err != nil {
		logger.Error("Failed to clone repository", "error", err)
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	logger.Info("Repository cloned successfully", "workspace", workspace)

	defer func() {
		if err := os.RemoveAll(workspace); err != nil {
			logger.Error("Failed to cleanup workspace", "workspace", workspace, "error", err)
		}
	}()

	appsDir := filepath.Join(workspace, "apps")
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 30 * time.Second,
		}),
		activities.UpdateAppVersion,
		appsDir,
		input.Namespace,
		input.App,
		input.Cluster,
		input.NewImages,
	).Get(ctx, nil); err != nil {
		logger.Error("failed to update app version", "error", err)
		return fmt.Errorf("failed to update app version: %w", err)
	}

	logger.Info("App version updated successfully",
		"namespace", input.Namespace,
		"app", input.App,
		"cluster", input.Cluster)

	// Step 3: Sync changes back to git
	appFilePath := filepath.Join(appsDir, input.Namespace, input.App, fmt.Sprintf("%s.yaml", input.Cluster))
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 1 * time.Minute,
		}),
		activities.GitSync,
		appFilePath,
	).Get(ctx, nil); err != nil {
		logger.Error("Failed to sync changes to git", "error", err)
		return fmt.Errorf("failed to sync changes to git: %w", err)
	}

	logger.Info("AppUpdate workflow completed successfully",
		"namespace", input.Namespace,
		"app", input.App,
		"cluster", input.Cluster,
		"updated_images", len(input.NewImages))

	return nil
}
