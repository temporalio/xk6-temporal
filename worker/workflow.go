package main

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func SingleActivityWorkflow(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	err := workflow.ExecuteActivity(ctx, SayHello, name).Get(ctx, &result)
	if err != nil {
		return "", err
	}

	return result, err
}
