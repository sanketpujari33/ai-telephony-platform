import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '10s', target: 10 }, 
    { duration: '30s', target: 10 }, 
    { duration: '10s', target: 0 },  
  ],
  thresholds: {
    http_req_duration: ['p(95)<300'], // Telephony should be faster
    http_req_failed: ['rate<0.01'],   
  },
};

const BASE_URL = __ENV.TELEPHONY_URL || 'http://localhost:8085';

export default function () {
  // Telephony typically exposes health or webhooks
  const res = http.get(`${BASE_URL}/health`);
  
  check(res, {
    'is status 200 or 404': (r) => r.status === 200 || r.status === 404,
  });

  sleep(1);
}
