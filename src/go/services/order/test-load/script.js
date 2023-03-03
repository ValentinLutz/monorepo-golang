import http from 'k6/http';
import encoding from 'k6/encoding';
import {check} from "k6";

export const BASE_URI = 'https://app:8443'

export const options = {
    insecureSkipTLSVerify: true,
    scenarios: {
        getOrders: {
            executor: 'constant-arrival-rate',
            exec: 'getOrders',
            timeUnit: '1s',
            rate: 100,
            duration: '5m',
            preAllocatedVUs: 1,
            maxVUs: 100
        },
        postOrder_getOrder: {
            executor: 'constant-arrival-rate',
            exec: 'postOrder_getOrder',
            timeUnit: '1s',
            rate: 100,
            duration: '5m',
            preAllocatedVUs: 1,
            maxVUs: 100
        },
    },
};

const credentials = `test:test`;
const encodedCredentials = encoding.b64encode(credentials);

export function getOrders() {
    const response = http.get(BASE_URI + '/api/orders', {
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

    const response = http.post(BASE_URI + '/api/orders', payload, {
        headers: {
            Authorization: `Basic ${encodedCredentials}`,
        },
    });

    check(response, {
        'postOrder is status 201': (r) => r.status === 201,
    });

    return response.json().order_id
}

export function postOrder_getOrder() {
    let orderId = postOrder()

    const response = http.get(BASE_URI + '/api/orders/' + orderId, {
        headers: {
            Authorization: `Basic ${encodedCredentials}`,
        },
    });

    check(response, {
        'getOrder is status 200': (r) => r.status === 200,
    });
}