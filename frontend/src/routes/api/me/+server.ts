// src/routes/api/me/+server.ts
import type { RequestHandler } from './$types';
import { backendFetch } from '$lib/server/backend';

export const GET: RequestHandler = async (event) => {
    const res = await backendFetch(event, '/me', { method: 'GET' });

    const body = await res.arrayBuffer();
    return new Response(body, {
        status: res.status,
        headers: pickResponseHeaders(res.headers)
    });
};

function pickResponseHeaders(h: Headers) {
    const out = new Headers();
    const contentType = h.get('content-type');
    if (contentType) out.set('content-type', contentType);
    return out;
}

