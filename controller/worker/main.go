package main

import (
	"log"
	"os"

	"cloudlab/controller/activities"
	"cloudlab/controller/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	temporalClient, err := client.Dial(client.Options{
		HostPort: os.Getenv("TEMPORAL_HOST"),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer temporalClient.Close()

	w := worker.New(temporalClient, "cloudlab", worker.Options{})

	w.RegisterActivity(activities.Clone)
	w.RegisterActivity(activities.ChangedModules)
	w.RegisterActivity(activities.TerragruntGraph)
	w.RegisterActivity(activities.PruneGraph)
	w.RegisterActivity(activities.TerragruntApply)
	w.RegisterActivity(activities.PushManifests)
	w.RegisterActivity(activities.PushRenderedApp)
	w.RegisterActivity(activities.DiscoverApps)
	w.RegisterActivity(activities.UpdateAppVersion)
	w.RegisterActivity(activities.GitSync)

	w.RegisterWorkflow(workflows.Infra)
	w.RegisterWorkflow(workflows.Platform)
	w.RegisterWorkflow(workflows.Apps)
	w.RegisterWorkflow(workflows.AppUpdate)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}
}
