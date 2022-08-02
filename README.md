# xk6-temporal

k6 extension for testing/benchmarking Temporal.

Note: This project is still a spike. The API may change at anytime as we learn from experience.

We recommend that you use this k6 extension alongside our benchmark workers which provide some pre-written workflow and activity workers that you can make use of for benchmarking. You can of course bring your own workflow and activity workers if you want to benchmark a closer simulation of your specific workload.

Our benchmark workers are available at: https://github.com/temporalio/benchmark-workers

## Usage

This extension is available compiled into k6 as docker image for use in Docker or Kubernetes setups.

You can pull the latest image from: `ghcr.io/temporalio/xk6-temporal:main`.

In future we will provide releases with appropriate image tags to make benchmarks more easily repeatable.

Before you run a benchmark make sure your workers are deployed and scaled as required. If you would like to use our pre-written benchmark workers rather than your own you can find details of how to do that on the (benchmark-workers site)[https://github.com/temporalio/benchmark-workers].

To run one of our example benchmark scripts against Temporal in a Kubernetes cluster you can use:

```
kubectl run k6-benchmark -ti \
    --image ghcr.io/temporalio/xk6-temporal:main \
    --image-pull-policy Always \
    --env TEMPORAL_GRPC_ENDPOINT=temporal-frontend.temporal:7233 \
    --
    k6 run ./examples/start-complete.js
```

You will see the benchmark progress and the final statistics on screen.

If you have prometheus setup we recommend that you also send k6 metrics to your prometheus instance to more easily tie k6 results to changes in Temporal metrics. To do this you can run the benchmark like so:

```
kubectl run k6-benchmark -ti \
    --image ghcr.io/temporalio/xk6-temporal:main \
    --image-pull-policy Always \
    --env TEMPORAL_GRPC_ENDPOINT=temporal-frontend.temporal:7233 \
    --env K6_PROMETHEUS_REMOTE_URL=http://temporal-prometheus.temporal:9090/api/v1/write \
    --
    k6 run -o output-prometheus-remote ./examples/start-complete.js
```

Note: Your prometheus instance will need to have remote write enabled for the metrics to be receieved, this is often not enabled by default.