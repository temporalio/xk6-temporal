package main

import (
	"context"
	"flag"
	"log"
	"net"
	"strconv"
	"time"

	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var portFlag int
var upstreamFlag string

func init() {
	flag.IntVar(&portFlag, "port", 8081, "Port to listen on")
	flag.StringVar(&upstreamFlag, "upstream", "127.0.0.1:7233", "Upstream Temporal Server Endpoint")
}

func main() {
	flag.Parse()

	grpcClient, err := grpc.Dial(
		upstreamFlag,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			log.Printf("> %v\n", req)
			start := time.Now()
			err := invoker(ctx, method, req, reply, cc, opts...)
			log.Printf("< [%v] %v\n", time.Since(start), reply)

			return err
		}),
	)
	defer func() { _ = grpcClient.Close() }()

	workflowClient := workflowservice.NewWorkflowServiceClient(grpcClient)

	if err != nil {
		log.Fatalf("unable to create client: %v", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(portFlag))
	if err != nil {
		log.Fatalf("unable to create listener: %v", err)
	}

	server := grpc.NewServer()
	handler, err := client.NewWorkflowServiceProxyServer(
		client.WorkflowServiceProxyOptions{Client: workflowClient},
	)
	if err != nil {
		log.Fatalf("unable to create service proxy: %v", err)
	}

	workflowservice.RegisterWorkflowServiceServer(server, handler)

	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("unable to serve: %v", err)
	}
}
