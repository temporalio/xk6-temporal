package temporal

import (
	"context"

	"go.k6.io/k6/js/modules"

	sdkclient "go.temporal.io/sdk/client"

	"github.com/temporalio/xk6-temporal/client"
	"github.com/temporalio/xk6-temporal/logger"
	"github.com/temporalio/xk6-temporal/metrics"
	"github.com/temporalio/xk6-temporal/worker"
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
	customMetrics metrics.CustomMetrics
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
		customMetrics: metrics.RegisterMetrics(vu.InitEnv().Registry),
	}
}

// Exports implements the modules.Instance interface and returns the exports
// of the JS module.
func (temporal *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{Default: temporal}
}

// NewClient returns a new Temporal Client.
func (m *ModuleInstance) NewClient(options client.Options) (*client.Client, error) {
	options.MetricsHandler = metrics.NewClientMetricsHandler(
		m.vu.Context(),
		m.vu.State().Samples,
		m.vu.State().Tags.Clone(),
		m.customMetrics,
	)

	return client.NewClient(options)
}

// NewWorker returns a new Temporal Worker with the example workflows registered.
func (m *ModuleInstance) NewWorker(clientOptions client.Options, options worker.Options) (worker.Worker, error) {
	clientOptions.MetricsHandler = metrics.NewClientMetricsHandler(
		context.Background(),
		m.vu.State().Samples,
		m.vu.State().Tags.Clone(),
		m.customMetrics,
	)
	// Not sure what to do with these logs, they may be useful.
	clientOptions.Logger = logger.NewNopLogger()

	c, err := sdkclient.Dial(clientOptions)
	if err != nil {
		return nil, err
	}

	return worker.NewWorker(c, options)
}
