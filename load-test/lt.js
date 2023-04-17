import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

export const options = {
    noVUConnectionReuse: false,
    noConnectionReuse: false,
    teardownTimeout: '1s',
    setupTimeout: '1s',
    thresholds: {
      'http_req_failed': ['rate<0.01'],
    },
    scenarios: {
        scenario_1: {
            executor: 'ramping-vus',
            gracefulStop: '3m',
            gracefulRampDown: '3m',
            stages: [
                {target: 200, duration: '1'},
                {target: 200, duration: '600s'},
            ],
            exec: 'scenario_1'
        },
    },
  insecureSkipTLSVerify: true,
};

const client = new grpc.Client();
client.load(['sgroups'], 'api.proto');



export function scenario_1 () {

  client.connect('51.250.78.167:80', {
    plaintext: true,
    timeout: "5s"
  });
    const payload = {
      "sgFrom": [
        "d3b8a0f2ca"
      ],
  
    }
    const response = client.invoke('hbf.v1.sgroups.SecGroupService/FindRules', payload);
  
    check(response, {
      'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
    sleep(1)
};
