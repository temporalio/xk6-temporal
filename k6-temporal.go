package temporal

import (
	"context"

	"go.temporal.io/sdk/client"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/temporal", new(RootModule))
}

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create `k6/x/sql` module instances for each VU.
type RootModule struct{}

// Temporal represents an instance of the Temporal module for every VU.
type Temporal struct{}

// Client is the exported module instance.
type Client struct {
	sdkclient client.Client
}

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &Temporal{}
)

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &Temporal{}
}

// Exports implements the modules.Instance interface and returns the exports
// of the JS module.
func (temporal *Temporal) Exports() modules.Exports {
	return modules.Exports{Default: temporal}
}

// NewClient returns a new Temporal Client.
func (*Temporal) NewClient(options client.Options) (*Client, error) {
	c, err := client.Dial(options)
	if err != nil {
		return nil, err
	}

	return &Client{sdkclient: c}, nil
}

type (
	WorkflowRun struct {
		run   client.WorkflowRun
		ID    string
		RunID string
	}

	WorkflowResult struct {
		Result interface{}
		Error  error
	}
)

func (r WorkflowRun) Get() WorkflowResult {
	var result WorkflowResult

	result.Error = r.run.Get(context.Background(), &result.Result)

	return result
}

func (c *Client) StartWorkflow(options client.StartWorkflowOptions, workflowType string, workflowArgs ...interface{}) (WorkflowRun, error) {
	run, err := c.sdkclient.ExecuteWorkflow(
		context.Background(),
		options,
		workflowType,
		workflowArgs...,
	)

	if err != nil {
		return WorkflowRun{}, err
	}

	return WorkflowRun{run: run, ID: run.GetID(), RunID: run.GetRunID()}, err
}
