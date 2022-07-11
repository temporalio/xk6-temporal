package worker

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

type SignalEchoRequestSignal struct {
	SenderWorkflowID string
	SenderRunID      string
	Message          string
}

type SignalEchoResponseSignal struct {
	Message string
}

func sendSignalEchoResponseSignal(ctx workflow.Context, request SignalEchoRequestSignal) error {
	f := workflow.SignalExternalWorkflow(
		ctx,
		request.SenderWorkflowID,
		request.SenderRunID,
		"SignalEchoResponseSignal",
		SignalEchoResponseSignal{
			Message: request.Message,
		},
	)
	return f.Get(ctx, nil)
}

func SignalEchoWorkflow(ctx workflow.Context) error {
	requestCh := workflow.GetSignalChannel(ctx, "SignalEchoRequestSignal")
	var request SignalEchoRequestSignal

	for i := 0; i < 1000; i++ {
		requestCh.Receive(ctx, &request)
		err := sendSignalEchoResponseSignal(ctx, request)
		if err != nil {
			return err
		}
	}

	for requestCh.ReceiveAsync(&request) {
		err := sendSignalEchoResponseSignal(ctx, request)
		if err != nil {
			return err
		}
	}

	return workflow.NewContinueAsNewError(ctx, SignalEchoWorkflow)
}

func SignalWaiterWorkflow(ctx workflow.Context, echoWorkflowID string, message string) (string, error) {
	responseCh := workflow.GetSignalChannel(ctx, "SignalEchoResponseSignal")
	var response SignalEchoResponseSignal

	info := workflow.GetInfo(ctx)

	f := workflow.SignalExternalWorkflow(
		ctx,
		echoWorkflowID,
		"",
		"SignalEchoRequestSignal",
		SignalEchoRequestSignal{
			SenderWorkflowID: info.WorkflowExecution.ID,
			SenderRunID:      info.WorkflowExecution.RunID,
			Message:          message,
		},
	)
	if err := f.Get(ctx, nil); err != nil {
		return "", err
	}

	responseCh.Receive(ctx, &response)

	return response.Message, nil
}
