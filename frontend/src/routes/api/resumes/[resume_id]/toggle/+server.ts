import type { RequestHandler } from './$types';
import { backendFetch } from '$lib/server/backend';

export const POST: RequestHandler = async (event) => {
    const resume_id = event.params.resume_id;
    const res = await backendFetch(event, `/resumes/${encodeURIComponent(resume_id)}/toggle`, {
        method: 'POST',
    });
    const referer = event.request.headers.get('referer') || '/';
    return Response.redirect(referer, 303);
};

function pickResponseHeaders(h: Headers) {
    const out = new Headers();
    const contentType = h.get('content-type');
    if (contentType) out.set('content-type', contentType);
    return out;
}

