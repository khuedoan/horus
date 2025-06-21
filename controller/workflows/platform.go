package workflows

import (
	"os"
	"time"

	"cloudlab/controller/activities"

	"go.temporal.io/sdk/workflow"
)

type PlatformInput struct {
	Url      string
	Revision string
	Registry string
	Cluster  string
}

func Platform(ctx workflow.Context, input PlatformInput) error {
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

	var pushResult *activities.PushResult
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 1 * time.Minute,
		}),
		activities.PushManifests,
		workspace+"/platform/"+input.Cluster,
		input.Registry+"/platform:"+input.Cluster,
	).Get(ctx, &pushResult); err != nil {
		return err
	}

	return nil
}
