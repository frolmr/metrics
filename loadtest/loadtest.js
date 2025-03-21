import http from 'k6/http';
import { check, sleep } from 'k6';
import { hmac } from 'k6/crypto';
import pako from 'https://cdnjs.cloudflare.com/ajax/libs/pako/2.0.4/pako.min.js';

const BASE_URL = 'http://localhost:8080';
const SHARED_KEY = __ENV.SECRET_KEY;

const gaugeMetrics = [
    { id: 'test_g_1', type: 'gauge', value: 1.0 },
    { id: 'test_g_2', type: 'gauge', value: 1.1 },
    { id: 'test_g_3', type: 'gauge', value: 1.2 },
    { id: 'test_g_4', type: 'gauge', value: 1.3 },
];

const counterMetrics = [
    { id: 'test_c_1', type: 'counter', delta: 1 },
    { id: 'test_c_2', type: 'counter', delta: 2 },
    { id: 'test_c_3', type: 'counter', delta: 3 },
    { id: 'test_c_4', type: 'counter', delta: 4 },
];

export const options = {
    vus: 10,
    duration: '1m',
};

function generateSignature(payload) {
    const signature = hmac('sha256', SHARED_KEY, payload, 'hex');
    return signature;
}

function sendSignedCompressedRequest(method, url, body, headers = {}) {
    const bodyString = JSON.stringify(body);
    const signature = generateSignature(bodyString);
    const compressedBody = pako.gzip(bodyString);

    headers['Content-Encoding'] = 'gzip';
    headers['Accept-Encoding'] = 'gzip';
    headers['Content-Type'] = 'application/json';
    headers['HashSHA256'] = signature;

    return http.request(method, url, compressedBody, { headers: headers });
}

export default function () {
    let res = http.get(`${BASE_URL}/`);
    check(res, {
        'GET / status was 200': (r) => r.status === 200,
    });

    gaugeMetrics.forEach((metric) => {
        let updatePayload = { id: metric.id, type: metric.type, value: metric.value + 0.1 };
        res = http.post(`${BASE_URL}/update/`, JSON.stringify(updatePayload), {
            headers: { 'Content-Type': 'application/json' },
        });
        check(res, {
            [`POST /update/ for ${metric.id} status was 200`]: (r) => r.status === 200,
        });
        res = sendSignedCompressedRequest('POST', `${BASE_URL}/update/`, updatePayload);
        check(res, {
            [`POST compressed /update/ for ${metric.id} status was 200`]: (r) => r.status === 200,
        });

        res = http.post(`${BASE_URL}/update/${metric.type}/${metric.id}/${metric.value + 0.2}`);
        check(res, {
            [`POST /${metric.type}/${metric.id}/${metric.value} status was 200`]: (r) => r.status === 200,
        });

        let valuePayload = { id: metric.id, type: metric.type };
        res = http.post(`${BASE_URL}/value/`, JSON.stringify(valuePayload), {
            headers: { 'Content-Type': 'application/json' },
        });
        check(res, {
            [`POST /value for ${metric.id} status was 200`]: (r) => r.status === 200,
        });
        res = sendSignedCompressedRequest('POST', `${BASE_URL}/value/`, valuePayload);
        check(res, {
            [`POST compressed /value for ${metric.id} status was 200`]: (r) => r.status === 200,
        });

        res = http.get(`${BASE_URL}/value/${metric.type}/${metric.id}`);
        check(res, {
            [`GET /value/${metric.type}/${metric.id} status was 200`]: (r) => r.status === 200,
        });
    });

    counterMetrics.forEach((metric) => {
        let updatePayload = { id: metric.id, type: metric.type, delta: metric.delta + 1 };
        res = http.post(`${BASE_URL}/update/`, JSON.stringify(updatePayload), {
            headers: { 'Content-Type': 'application/json' },
        });
        check(res, {
            [`POST /update/ for ${metric.id} status was 200`]: (r) => r.status === 200,
        });
        res = sendSignedCompressedRequest('POST', `${BASE_URL}/update/`, updatePayload);
        check(res, {
            [`POST compressed /update/ for ${metric.id} status was 200`]: (r) => r.status === 200,
        });

        res = http.post(`${BASE_URL}/update/${metric.type}/${metric.id}/${metric.delta + 2}`);
        check(res, {
            [`POST /${metric.type}/${metric.id}/{delta} status was 200`]: (r) => r.status === 200,
        });

        let valuePayload = { id: metric.id, type: metric.type };
        res = http.post(`${BASE_URL}/value/`, JSON.stringify(valuePayload), {
            headers: { 'Content-Type': 'application/json' },
        });
        check(res, {
            [`POST /value for ${metric.id} status was 200`]: (r) => r.status === 200,
        });
        res = sendSignedCompressedRequest('POST', `${BASE_URL}/value/`, valuePayload);
        check(res, {
            [`POST compressed /value for ${metric.id} status was 200`]: (r) => r.status === 200,
        });

        res = http.get(`${BASE_URL}/value/${metric.type}/${metric.id}`);
        check(res, {
            [`GET /value/${metric.type}/${metric.id} status was 200`]: (r) => r.status === 200,
        });
    });

    res = http.get(`${BASE_URL}/ping`);
    check(res, {
        'GET /ping status was 200': (r) => r.status === 200,
    });

    res = http.post(`${BASE_URL}/updates/`, JSON.stringify(gaugeMetrics), {
        headers: { 'Content-Type': 'application/json' },
    });
    check(res, {
        'POST /updates status was 200': (r) => r.status === 200,
    });

    res = http.post(`${BASE_URL}/updates/`, JSON.stringify(counterMetrics), {
        headers: { 'Content-Type': 'application/json' },
    });
    check(res, {
        'POST /updates status was 200': (r) => r.status === 200,
    });

    res = sendSignedCompressedRequest('POST', `${BASE_URL}/updates/`, gaugeMetrics);
    check(res, {
        'Compressed POST /updates status was 200': (r) => r.status === 200,
    });

    res = sendSignedCompressedRequest('POST', `${BASE_URL}/updates/`, counterMetrics);
    check(res, {
        'Compressed POST /updates status was 200': (r) => r.status === 200,
    });

    sleep(1);
}
