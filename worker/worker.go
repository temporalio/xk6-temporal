package main

import (
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort:       os.Getenv("TEMPORAL_GRPC_ENDPOINT"),
		MetricsHandler: newMetricsHandler("0.0.0.0:9000"),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "benchmark", worker.Options{})

	w.RegisterWorkflow(SingleActivityWorkflow)

	w.RegisterActivity(SayHello)

	err = w.Run(nil)
	if err != nil {
		log.Fatal(err)
	}
}
