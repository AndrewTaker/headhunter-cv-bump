import type { PageServerLoad } from './$types';
import { backendFetch } from '$lib/server/backend';

export const load: PageServerLoad = async (event) => {
    const res = await backendFetch(event, '/resumes', { method: 'GET' });
    if (res.status !== 200) return { resumes: [] };
    const body = await res.json();
    return { resumes: body.resumes ?? [] };
};

