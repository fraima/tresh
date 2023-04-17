import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new grpc.Client();
client.load(['protos/api/sgroups'], 'api.proto');

export default () => {
  client.connect('51.250.78.167:80', {
    plaintext: false
  });

  const payload = {
    "sgFrom": [
      "d3b8a0f2ca"
    ],
    "sgTo": [
      "a3e7cbe612"
    ]
  }
  const response = client.invoke('hbf.v1.sgroups.SecGroupService/FindRules', payload);

  check(response, {
    'status is OK': (r) => r && r.status === grpc.StatusOK,
  });

  console.log(JSON.stringify(response.message));

  client.close();
  sleep(1);
};