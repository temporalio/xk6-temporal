import temporal from 'k6/x/temporal';
import { scenario } from 'k6/execution';

export const options = {
    scenarios: {
      start_complete: {
        executor: 'constant-vus',
        duration: '5m',
        vus: 200,
      },
    },
};

export default () => {
    const client = temporal.newClient()

    const handle = client.startWorkflow(
        {
            task_queue: 'benchmark',
            id: 'wf-' + scenario.iterationInTest,
        },
        'ExecuteActivity',
        {"Count": 1, "Activity": "Echo", "Input": {"Message": "test"}},
    )

    // Wait until the workflow has completed.
    handle.result()

    client.close()
};