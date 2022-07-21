import temporal from 'k6/x/temporal';
import { scenario } from 'k6/execution';

export const options = {
  scenarios: {
    signal_waiters_start: {
      executor: 'shared-iterations',
      iterations: '10000',
      vus: 100,
      exec: 'starter',
    },
    signal_waiters_complete: {
      executor: 'shared-iterations',
      iterations: '10000',
      vus: 100,
      exec: 'waiter',
      startTime: '5s',
    },
  },
};

const ECHO_WORKFLOW_SHARD = 4

export function setup() {
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