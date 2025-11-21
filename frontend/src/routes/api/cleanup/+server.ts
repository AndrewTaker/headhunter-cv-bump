// src/routes/api/cleanup/+server.ts
import type { RequestHandler } from './$types';
import { backendFetch } from '$lib/server/backend';

export const DELETE: RequestHandler = async (event) => {
    const res = await backendFetch(event, '/cleanup', { method: 'DELETE' });

    const body = await res.arrayBuffer();
    return new Response(body, {
        status: res.status,
        headers: pickResponseHeaders(res.headers)
    });
};

function pickResponseHeaders(h: Headers) {
    const out = new Headers();
    const ct = h.get('content-type');
    if (ct) out.set('content-type', ct);
    return out;
}

