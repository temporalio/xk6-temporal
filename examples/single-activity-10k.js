import temporal from 'k6/x/temporal';
import { scenario } from 'k6/execution';

export const options = {
    scenarios: {
      single_activity_workflow_10k: {
        executor: 'shared-iterations',
        iterations: '10000',
        vus: 100,
        maxDuration: '30m',
      },
    },
};

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
    const result = handle.result()

    client.close()
};