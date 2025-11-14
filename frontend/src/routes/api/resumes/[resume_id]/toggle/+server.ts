// src/routes/api/resumes/[resume_id]/toggle/+server.ts
import type { RequestHandler } from './$types';
import { backendFetch } from '$lib/server/backend';

export const POST: RequestHandler = async (event) => {
    const resume_id = event.params.resume_id;
    const res = await backendFetch(event, `/resumes/${encodeURIComponent(resume_id)}/toggle`, {
        method: 'POST',
    });
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

