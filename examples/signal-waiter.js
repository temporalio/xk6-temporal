import temporal from 'k6/x/temporal';
import { scenario } from 'k6/execution';

export const options = {
  scenarios: {
    signal_waiters_start: {
      executor: 'shared-iterations',
      iterations: '1000',
      vus: 100,
      exec: 'starter',
    },
    signal_waiters_complete: {
      executor: 'shared-iterations',
      iterations: '1000',
      vus: 100,
      exec: 'waiter',
      startTime: '3s',
    },
  },
};

// Experiment by changing the shard count.

// A good comparison is to try with ECHO_WORKFLOW_SHARD=1 and
// ECHO_WORKFLOW_SHARD=4. You should see a marked improvement in test speed with
// 4 shards versus 1 as it reduces the contention on the lock which is required
// to update the workflow history for the signal echo workflow as it receives
// and sends signals.

// Try to avoid any Temporal designs that involve a single workflow sending or
// receiving a large number of signals. Where possible this responsiblity should
// be shared among multiple workflow executions.

const ECHO_WORKFLOW_SHARD = 4

export function setup() {
  temporal.newWorker(
    { host_port: __ENV.TEMPORAL_GRPC_ENDPOINT },
    {
      max_concurrent_workflow_task_pollers: 16,
      max_concurrent_activity_task_pollers: 16,
    }
  ).start()

  const client = temporal.newClient({ host_port: __ENV.TEMPORAL_GRPC_ENDPOINT })
  
  for (let i = 0; i < ECHO_WORKFLOW_SHARD; i++) {
    client.startWorkflow(
      {
        task_queue: 'benchmark',
        id: 'signal-echo-' + i,
      },
      'SignalEchoWorkflow',
    )  
  }

  client.close()
}

export function teardown() {
  const client = temporal.newClient({ host_port: __ENV.TEMPORAL_GRPC_ENDPOINT })

  for (let i = 0; i < ECHO_WORKFLOW_SHARD; i++) {
    client.getWorkflowHandle('signal-echo-' + i, "").terminate()
  }

  client.close()
}

export function starter() {
  const client = temporal.newClient({ host_port: __ENV.TEMPORAL_GRPC_ENDPOINT })
  const wfID = 'signal-waiter-' + scenario.iterationInTest
  const echoID = 'signal-echo-' + (scenario.iterationInTest % ECHO_WORKFLOW_SHARD)

  client.startWorkflow(
    {
      task_queue: 'benchmark',
      id: wfID,
    },
    'SignalWaiterWorkflow',
    echoID,
    `ping from VU ${scenario.iterationInTest}`
  )
  client.close()
};

export function waiter() {
  const client = temporal.newClient({ host_port: __ENV.TEMPORAL_GRPC_ENDPOINT })
  const wfID = 'signal-waiter-' + scenario.iterationInTest

  client.getWorkflowHandle(wfID, "").result()

  client.close()
};