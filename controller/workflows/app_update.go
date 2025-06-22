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
	Registry  string
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

	// Step 3: Git add changes
	appFilePath := filepath.Join(appsDir, input.Namespace, input.App, fmt.Sprintf("%s.yaml", input.Cluster))
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 30 * time.Second,
		}),
		activities.GitAdd,
		appFilePath,
	).Get(ctx, nil); err != nil {
		logger.Error("Failed to add changes to git", "error", err)
		return fmt.Errorf("failed to add changes to git: %w", err)
	}

	// Step 4: Git commit changes
	commitMessage := fmt.Sprintf("chore(%s/%s): update %s version", input.Namespace, input.App, input.Cluster)
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 30 * time.Second,
		}),
		activities.GitCommit,
		workspace,
		commitMessage,
	).Get(ctx, nil); err != nil {
		logger.Error("Failed to commit changes to git", "error", err)
		return fmt.Errorf("failed to commit changes to git: %w", err)
	}

	// Step 5: Git push changes
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 1 * time.Minute,
		}),
		activities.GitPush,
		workspace,
	).Get(ctx, nil); err != nil {
		logger.Error("Failed to push changes to git", "error", err)
		return fmt.Errorf("failed to push changes to git: %w", err)
	}

	// Step 6: Push rendered app to registry
	var pushResult *activities.PushResult
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 2 * time.Minute,
		}),
		activities.PushRenderedApp,
		appsDir,
		input.Namespace,
		input.App,
		input.Cluster,
		input.Registry,
	).Get(ctx, &pushResult); err != nil {
		logger.Error("Failed to push rendered app to registry", "error", err)
		return fmt.Errorf("failed to push rendered app to registry: %w", err)
	}

	logger.Info("AppUpdate workflow completed successfully",
		"namespace", input.Namespace,
		"app", input.App,
		"cluster", input.Cluster,
		"updated_images", len(input.NewImages),
		"rendered_app_digest", pushResult.Digest)

	return nil
}
