package worker

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type Options = worker.Options
type Worker = worker.Worker

func NewWorker(clientOptions client.Options, options Options) (worker.Worker, error) {
	c, err := client.Dial(clientOptions)
	if err != nil {
		return nil, err
	}

	w := worker.New(c, "benchmark", options)

	w.RegisterWorkflow(SingleActivityWorkflow)
	w.RegisterActivity(SayHello)

	return w, nil
}
