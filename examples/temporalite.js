import temporal from 'k6/x/temporal';
import { scenario } from 'k6/execution';

export const options = {
  scenarios: {
    temporalite: {
      executor: 'shared-iterations',
      iterations: '100',
      vus: 20,
    },
  },
};

export function setup() {
  const server = temporal.newTemporaliteServer({})

  temporal.newWorker(
    { host_port: server.frontendHostPort() },
    {
      max_concurrent_workflow_task_pollers: 16,
      max_concurrent_activity_task_pollers: 16,
    }
  ).start()

  return { server: { grpc_endpoint: server.frontendHostPort() } }
}

export default (data) => {
  const server = data.server

  const client = temporal.newClient({ host_port: server.grpc_endpoint })

  const handle = client.startWorkflow(
    {
      task_queue: 'benchmark',
      id: 'wf-' + scenario.iterationInTest,
    },
    'SingleActivityWorkflow',
    'bob',
  )

  // Wait until the workflow has completed.
  const result = handle.result()

  client.close()
};