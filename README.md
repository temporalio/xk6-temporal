# xk6-temporal

k6 Extension for testing/benchmarking Temporal.

Note: This project is still a spike. The API may change at anytime as we learn from experience.

## Usage

```
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
        maxVUs: 10,
      },
    },
};

let client = temporal.newClient({ HostPort: "127.0.0.1:7233" })

export function teardown () {
    console.log("Waiting for workflows to complete")

    client.waitForAllWorkflowToComplete()
}

export default () => {
    client.startWorkflow(
        {
            namespace: 'default',
            task_queue: 'benchtest',
            id: 'wf-' + scenario.iterationInTest,
        },
        'BenchTestWorkflow',
        'john',
    )
};
```