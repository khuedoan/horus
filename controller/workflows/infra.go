package workflows

import (
	"time"

	"cloudlab/controller/activities"
	"go.temporal.io/sdk/workflow"
)

type InfraInputs struct {
	Url      string
	Revision string
	Stack    string
}

// TODO create trigger
// For now do that manually on the UI
// Task queue: cloudlab
// Workflow: Infra
// Input json/plain: {"url": "https://github.com/khuedoan/cloudlab", "revision": "infra-rewrite", "stack": "local"}
func Infra(ctx workflow.Context, input InfraInputs) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("Infra workflow started", "infra", input)

	var path string
	err := workflow.ExecuteActivity(ctx, activities.Clone, input.Url, input.Revision).Get(ctx, &path)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	var dotGraph string
	err = workflow.ExecuteActivity(ctx, activities.TerragruntGraph, path+"/infra/"+input.Stack).Get(ctx, &dotGraph)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	var result string
	err = workflow.ExecuteActivity(ctx, activities.TerragruntTreeShaking, dotGraph, []string{"cluster"}).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("Infra workflow completed.", "result", result)

	return result, nil
}
