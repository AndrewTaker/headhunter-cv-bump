// src/routes/+layout.server.ts
import type { LayoutServerLoad } from './$types';
import { backendFetch } from '$lib/server/backend';

export const load: LayoutServerLoad = async (event) => {
    const res = await backendFetch(event, '/me', { method: 'GET' });

    if (res.status !== 200) return { user: null };

    const user = await res.json();
    return { user };
};

