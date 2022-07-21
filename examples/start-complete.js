import temporal from 'k6/x/temporal';
import { scenario } from 'k6/execution';

export const options = {
    scenarios: {
      min_pollers_high_wps: {
        executor: 'constant-vus',
        duration: '5m',
        vus: 200,
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
        'MyWorkflow',
        'bob',
    )

    // Wait until the workflow has completed.
    handle.result()

    client.close()
};