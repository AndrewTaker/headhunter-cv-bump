import type { RequestHandler } from './$types';
import { BACKEND_ORIGIN } from '$lib/server/backend';

export const GET: RequestHandler = async () => {
    return new Response(null, {
        status: 302,
        headers: {
            location: `${BACKEND_ORIGIN}/auth/logout`
        }
    });
};

