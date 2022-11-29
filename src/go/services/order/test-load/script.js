import http from 'k6/http';

export const BASE_URI = 'http://app:8080'
export const VIRTUAL_USERS = 100
export const ITERATIONS = 100

export const options = {
    scenarios: {
        getOrders_postOrder_getOrder: {
            executor: 'per-vu-iterations',
            exec: 'getOrder',
            vus: VIRTUAL_USERS,
            iterations: ITERATIONS,
        },
    },
};

export function getOrders() {
    http.get(BASE_URI + '/api/orders');
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

    const response = http.post(BASE_URI + '/api/orders', payload);

    return response.json()
}

export function getOrder() {
    let orderId = postOrder().order_id

    const response = http.get(BASE_URI + '/api/orders/' + orderId);

    return response.json()
}