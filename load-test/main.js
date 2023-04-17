import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new grpc.Client();
client.load(['sgroups'], 'api.proto');

export default () => {
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

  client.close();
  sleep(1);
};