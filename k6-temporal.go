package temporal

import (
	"context"
	"time"

	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/temporal", new(RootModule))
}

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create `k6/x/temporal` module instances for each VU.
type RootModule struct{}

// ModuleInstance represents an instance of the module for every VU.
type ModuleInstance struct {
	vu            modules.VU
	customMetrics CustomMetrics
}

// Client is the exported module instance.
type Client struct {
	sdkclient client.Client
}

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &ModuleInstance{}
)

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ModuleInstance{
		vu:            vu,
		customMetrics: RegisterMetrics(vu.InitEnv().Registry),
	}
}

// Exports implements the modules.Instance interface and returns the exports
// of the JS module.
func (temporal *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{Default: temporal}
}

// NewClient returns a new Temporal Client.
func (m *ModuleInstance) NewClient(options client.Options) (*Client, error) {
	options.Logger = NewNopLogger()
	options.MetricsHandler = NewMetricsHandler(m.vu, map[string]string{}, m.customMetrics)

	c, err := client.Dial(options)
	if err != nil {
		return nil, err
	}

	return &Client{sdkclient: c}, nil
}

type (
	WorkflowHandle struct {
		client *Client
		run    client.WorkflowRun
		ID     string
		RunID  string
	}
)

func (r WorkflowHandle) Result() (interface{}, error) {
	var result interface{}

	err := r.run.Get(context.Background(), &result)

	return result, err
}

func (r WorkflowHandle) Signal(name string, arg interface{}) error {
	return r.client.sdkclient.SignalWorkflow(
		context.Background(),
		r.ID,
		r.RunID,
		name,
		arg,
	)
}

func (r WorkflowHandle) Cancel() error {
	return r.client.sdkclient.CancelWorkflow(
		context.Background(),
		r.ID,
		r.RunID,
	)
}

func (r WorkflowHandle) Terminate(reason string) error {
	return r.client.sdkclient.TerminateWorkflow(
		context.Background(),
		r.ID,
		r.RunID,
		reason,
	)
}

func (c *Client) Close() {
	c.sdkclient.Close()
}

func (c *Client) GetWorkflowHandle(workflowID string, runID string) WorkflowHandle {
	run := c.sdkclient.GetWorkflow(context.Background(), workflowID, runID)
	return WorkflowHandle{client: c, run: run, ID: workflowID, RunID: runID}
}

func (c *Client) StartWorkflow(options client.StartWorkflowOptions, workflowType string, workflowArgs ...interface{}) (WorkflowHandle, error) {
	run, err := c.sdkclient.ExecuteWorkflow(
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

func (c *Client) SignalWithStartWorkflow(workflowID string, signalName string, signalArg interface{}, options client.StartWorkflowOptions, workflowType string, workflowArgs ...interface{}) (WorkflowHandle, error) {
	run, err := c.sdkclient.SignalWithStartWorkflow(
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
		response, err := c.sdkclient.ListOpenWorkflow(context.Background(), &request)
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
