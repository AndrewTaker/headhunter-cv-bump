// src/routes/login/+server.ts
import type { RequestHandler } from './$types';
import { BACKEND_ORIGIN } from '$lib/server/backend';

export const GET: RequestHandler = async () => {
    // Redirect the browser to the backend login entrypoint which starts OAuth flow
    return new Response(null, {
        status: 302,
        headers: {
            location: `${BACKEND_ORIGIN}/auth/login`
        }
    });
};

