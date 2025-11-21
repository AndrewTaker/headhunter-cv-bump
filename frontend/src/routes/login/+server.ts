import type { RequestHandler } from './$types';
import { BACKEND_ORIGIN } from '$lib/server/backend';

export const GET: RequestHandler = async () => {
    return new Response(null, {
        status: 302,
        headers: {
            location: "http://localhost:8080/auth/login",
        },
    });
};
