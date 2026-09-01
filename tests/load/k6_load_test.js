import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metric to track error rates
export let errorRate = new Rate('errors');

export let options = {
    // Phase 1: Ramp up to 10k concurrent virtual users (representing inbound SIP calls hitting gateway)
    // Phase 2: Hold for 2 minutes to stress the Orchestrator
    // Phase 3: Ramp down
    stages: [
        { duration: '30s', target: 5000 },
        { duration: '1m', target: 10000 },
        { duration: '2m', target: 10000 },
        { duration: '30s', target: 0 },
    ],
    thresholds: {
        // Assert that 95% of requests complete within 1.5s
        http_req_duration: ['p(95)<1500'],
        // Assert that less than 1% of requests return an error
        errors: ['rate<0.01'],
    },
};

export default function () {
    // The API Gateway endpoint we created in T12.1 for call initiation
    const url = 'http://localhost:8080/api/v1/calls';
    
    const payload = JSON.stringify({
        phone_number: '+18005559999',
        tenant_id: 'tenant_perf_test',
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
            // Mock operator token for the JWT RBAC we implemented in T12.4
            'Authorization': 'Bearer operator-test-token', 
        },
    };

    // Simulate an inbound call triggering the Call Router and Orchestrator
    let res = http.post(url, payload, params);

    let success = check(res, {
        'status is 202': (r) => r.status === 202,
        'has call_id': (r) => {
            try {
                return JSON.parse(r.body).call_id !== undefined;
            } catch (e) {
                console.log(`Failed to parse JSON. Status: ${r.status}, Body: ${r.body}`);
                return false;
            }
        },
    });

    errorRate.add(!success);

    // Simulate think time before placing the next call
    sleep(1);
}
