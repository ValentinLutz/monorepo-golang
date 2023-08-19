import http from 'k6/http';
import encoding from 'k6/encoding';
import {check} from "k6";

export const BASE_URI = 'https://localhost:10443'

export const options = {
    insecureSkipTLSVerify: true,
    thresholds: {
        http_req_duration: ['p(90) < 100', 'p(95) < 200'],
    },
    scenarios: {
        full_scenario: {
            executor: 'constant-vus',
            exec: 'fullScenario',
            vus: 10,
            duration: '300s',
        },
    },
};

const credentials = `test:test`;
const encodedCredentials = encoding.b64encode(credentials);

export function fullScenario() {
    const order_id = postOrder().order_id;
    getOrder(order_id);
    getOrders();
}

export function getOrders() {
    const response = http.get(BASE_URI + '/orders', {
        headers: {
            Authorization: `Basic ${encodedCredentials}`,
        },
    });

    check(response, {
        'getOrders is status 200': (r) => r.status === 200,
    });
}

export function postOrder() {
    const payload = JSON.stringify({
        "items": [
            {
                "name": "avocado"
            },
            {
                "name": "blueberry"
            },
            {
                "name": "lemon"
            }
        ]
    });

    const response = http.post(BASE_URI + '/orders', payload, {
        headers: {
            Authorization: `Basic ${encodedCredentials}`,
        },
    });

    check(response, {
        'postOrder is status 201': (r) => r.status === 201,
    });

    return response.json()
}

export function getOrder(order_id) {
    const response = http.get(BASE_URI + '/orders/' + order_id, {
        headers: {
            Authorization: `Basic ${encodedCredentials}`,
        },
    });

    check(response, {
        'getOrder is status 200': (r) => r.status === 200,
    });
}