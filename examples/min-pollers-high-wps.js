import temporal from 'k6/x/temporal';
import { scenario } from 'k6/execution';

export const options = {
    thresholds: {
      temporal_workflow_task_schedule_to_start_latency: [
        { threshold: 'max<10000', abortOnFail: true }
      ],
      temporal_activity_schedule_to_start_latency: [
        { threshold: 'max<10000', abortOnFail: true }
      ],
    },
    scenarios: {
      min_pollers_high_wps: {
        executor: 'shared-iterations',
        iterations: '1000',
        vus: 100,
      },
    },
};

export function setup() {
  // Experiment by reducing the poller counts below.
  
  // If the poller counts are too low and there are not pollers on each task
  // queue partition the worker can block long polling on an empty task queue
  // partition. This will show in the metrics as high
  // temporal_workflow_task_schedule_to_start_latency and
  // temporal_activity_schedule_to_start_latency.

  temporal.newWorker(
    { host_port: __ENV.TEMPORAL_GRPC_ENDPOINT },
    {
      max_concurrent_workflow_task_pollers: 8,
      max_concurrent_activity_task_pollers: 8,
    }
  ).start()
}

export default () => {
    const client = temporal.newClient({ host_port: __ENV.TEMPORAL_GRPC_ENDPOINT })

    const handle = client.startWorkflow(
        {
            task_queue: 'benchmark',
            id: 'wf-' + scenario.iterationInTest,
        },
        'SingleActivityWorkflow',
        'bob',
    )

    // Wait until the workflow has completed.
    handle.result()

    client.close()
};