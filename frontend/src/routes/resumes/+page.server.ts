// src/routes/resumes/+page.server.ts
import type { PageServerLoad } from './$types';
import { backendFetch } from '$lib/server/backend';

export const load: PageServerLoad = async (event) => {
    const res = await backendFetch(event, '/resumes', { method: 'GET' });
    console.log('res status', res.status);
    if (res.status !== 200) return { resumes: [] };
    const body = await res.json();
    console.log('res body', body);
    return { resumes: body.resumes ?? [] };
};

