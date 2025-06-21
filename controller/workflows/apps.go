package workflows

import (
	"time"

	"cloudlab/controller/activities"

	"go.temporal.io/sdk/temporal"
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

	var path string
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 1 * time.Minute,
			RetryPolicy: &temporal.RetryPolicy{
				MaximumAttempts: 3,
			},
		}),
		activities.Clone,
		input.Url,
		input.Revision,
	).Get(ctx, &path); err != nil {
		return err
	}

	var pushResult *activities.PushResult
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 5 * time.Minute,
			RetryPolicy: &temporal.RetryPolicy{
				MaximumAttempts: 2,
			},
		}),
		activities.PushRenderedHelm,
		// TODO loop through this
		path+"/apps",
		"khuedoan",
		"blog",
		input.Cluster,
		input.Registry,
	).Get(ctx, &pushResult); err != nil {
		return err
	}

	return nil
}
