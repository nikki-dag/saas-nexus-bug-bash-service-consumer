package main

import (
	"context"
	"github.com/nikki-dag/saas-nexus-bug-bash-api/service"
	"log"
	"os"
	"time"

	"go.temporal.io/sdk/client"

	"github.com/temporalio/saas-nexus-bug-bash-service-consumer/app"
)

func main() {
	clientOptions, err := app.ParseClientOptionFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}
	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()
	runWorkflow(c, app.EchoCallerWorkflow, "Nexus Bug Bash 👋")
	runWorkflow(c, app.HelloCallerWorkflow, "Nexus", service.ES)
}

func runWorkflow(c client.Client, workflow interface{}, args ...interface{}) {
	ctx := context.Background()
	workflowOptions := client.StartWorkflowOptions{
		ID:        "nexus_bug_bash_caller_workflow_" + time.Now().Format("20060102150405"),
		TaskQueue: app.TaskQueue,
	}

	wr, err := c.ExecuteWorkflow(ctx, workflowOptions, workflow, args)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", wr.GetID(), "RunID", wr.GetRunID())

	// Synchronously wait for the workflow completion.
	var result string
	err = wr.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)
}
