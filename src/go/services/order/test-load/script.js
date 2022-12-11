import http from 'k6/http';
import encoding from 'k6/encoding';
import {check} from "k6";

export const BASE_URI = 'https://app:8443'
export const VIRTUAL_USERS = 100
export const ITERATIONS = 100

export const options = {
    insecureSkipTLSVerify: true,
    scenarios: {
        getOrders: {
            executor: 'per-vu-iterations',
            exec: 'getOrders',
            vus: VIRTUAL_USERS,
            iterations: ITERATIONS,
        },
        postOrder_getOrder: {
            executor: 'per-vu-iterations',
            exec: 'postOrder_getOrder',
            vus: VIRTUAL_USERS,
            iterations: ITERATIONS,
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