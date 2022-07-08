package worker

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type Options = worker.Options
type Worker = worker.Worker

func NewWorker(client client.Client, options Options) (worker.Worker, error) {
	w := worker.New(client, "benchmark", options)

	w.RegisterWorkflow(SingleActivityWorkflow)
	w.RegisterActivity(SayHello)

	return w, nil
}
