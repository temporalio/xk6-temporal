import temporal from 'k6/x/temporal';
import { scenario } from 'k6/execution';

export const options = {
    teardownTimeout: '10m',
    scenarios: {
      contacts: {
        executor: 'constant-arrival-rate',
        duration: '60s',
        rate: 50,
        timeUnit: '1s',
        preAllocatedVUs: 1,
        maxVUs: 1000,
      },
    },
};

export default () => {
    const client = temporal.newClient({ HostPort: "127.0.0.1:7233" })

    const handle = client.startWorkflow(
        {
            namespace: 'default',
            task_queue: 'benchtest',
            id: 'wf-' + scenario.iterationInTest,
        },
        'BenchTestWorkflow',
        'john',
    )

    // Wait until the workflow has completed.
    const result = handle.result()

    client.close()
};