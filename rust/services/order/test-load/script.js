import http from 'k6/http';
import { check } from "k6";

export const BASE_URI = 'http://localhost:8080'

export const options = {
    insecureSkipTLSVerify: true,
    scenarios: {
        getOrders: {
            executor: 'per-vu-iterations',
            exec: 'getOrders_postOrder_getOrder',
            vus: 1000,
            iterations: 100,
        },
    },
};

export function getOrders() {
    const response = http.get(BASE_URI + '/api/orders', {
        headers: { 'Content-Type': 'application/json' },
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
        headers: { 'Content-Type': 'application/json' },
    });

    check(response, {
        'postOrder is status 201': (r) => r.status === 201,
    });

    return response.json().order_id
}

export function postOrder_getOrder() {
    let orderId = postOrder()

    const response = http.get(BASE_URI + '/api/orders/' + orderId, {
        headers: { 'Content-Type': 'application/json' },
    });

    check(response, {
        'getOrder is status 200': (r) => r.status === 200,
    });
}

export function getOrders_postOrder_getOrder() {
    getOrders()
    postOrder()
}