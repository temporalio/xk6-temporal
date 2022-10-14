package client

import (
	"context"
	"errors"
	"time"

	"go.temporal.io/api/workflowservice/v1"
	sdkclient "go.temporal.io/sdk/client"

	"github.com/temporalio/xk6-temporal/logger"
)

// Client is the exported module instance.
type Client struct {
	sdkClient sdkclient.Client
}

type Options = sdkclient.Options

type WorkflowHandle struct {
	client *Client
	run    sdkclient.WorkflowRun
	ID     string
	RunID  string
}

func NewClient(options Options) (*Client, error) {
	// Not sure what to do with these logs, they may be useful.
	options.Logger = logger.NewNopLogger()

	c, err := sdkclient.Dial(options)
	if err != nil {
		return nil, err
	}

	return &Client{sdkClient: c}, nil
}

func (r WorkflowHandle) Result() (interface{}, error) {
	var result interface{}

	err := r.run.Get(context.Background(), &result)
	for errors.Is(err, context.DeadlineExceeded) {
		err = r.run.Get(context.Background(), &result)
	}

	return result, err
}

func (r WorkflowHandle) Signal(name string, arg interface{}) error {
	return r.client.sdkClient.SignalWorkflow(
		context.Background(),
		r.ID,
		r.RunID,
		name,
		arg,
	)
}

func (r WorkflowHandle) Cancel() error {
	return r.client.sdkClient.CancelWorkflow(
		context.Background(),
		r.ID,
		r.RunID,
	)
}

func (r WorkflowHandle) Terminate(reason string) error {
	return r.client.sdkClient.TerminateWorkflow(
		context.Background(),
		r.ID,
		r.RunID,
		reason,
	)
}

func (c *Client) GetSDKClient() sdkclient.Client {
	return c.sdkClient
}

func (c *Client) Close() {
	c.sdkClient.Close()
}

func (c *Client) GetWorkflowHandle(workflowID string, runID string) WorkflowHandle {
	run := c.sdkClient.GetWorkflow(context.Background(), workflowID, runID)
	return WorkflowHandle{client: c, run: run, ID: workflowID, RunID: runID}
}

func (c *Client) StartWorkflow(options sdkclient.StartWorkflowOptions, workflowType string, workflowArgs ...interface{}) (WorkflowHandle, error) {
	run, err := c.sdkClient.ExecuteWorkflow(
		context.Background(),
		options,
		workflowType,
		workflowArgs...,
	)

	if err != nil {
		return WorkflowHandle{}, err
	}

	return WorkflowHandle{client: c, run: run, ID: run.GetID(), RunID: run.GetRunID()}, err
}

func (c *Client) SignalWithStartWorkflow(workflowID string, signalName string, signalArg interface{}, options sdkclient.StartWorkflowOptions, workflowType string, workflowArgs ...interface{}) (WorkflowHandle, error) {
	run, err := c.sdkClient.SignalWithStartWorkflow(
		context.Background(),
		workflowID,
		signalName,
		signalArg,
		options,
		workflowType,
		workflowArgs...,
	)

	if err != nil {
		return WorkflowHandle{}, err
	}

	return WorkflowHandle{client: c, run: run, ID: run.GetID(), RunID: run.GetRunID()}, err
}

func (c *Client) WaitForAllWorkflowToComplete(namespace string) error {
	request := workflowservice.ListOpenWorkflowExecutionsRequest{
		Namespace:       namespace,
		MaximumPageSize: 1,
	}

	for {
		response, err := c.sdkClient.ListOpenWorkflow(context.Background(), &request)
		if err != nil {
			return err
		}

		if len(response.Executions) == 0 {
			break
		}

		time.Sleep(time.Second)
	}

	return nil
}
