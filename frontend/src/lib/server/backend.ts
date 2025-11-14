// src/lib/server/backend.ts
import type { RequestEvent } from '@sveltejs/kit';

export const BACKEND_ORIGIN = process.env.BACKEND_ORIGIN ?? 'http://localhost:44444';

/**
 * Proxy a request to the Go backend.
 * - Forwards incoming request cookies to backend.
 * - Forwards method, headers (except host) and body.
 * - Returns backend response.
 */
export async function backendFetch(event: RequestEvent, path: string, opts: RequestInit = {}) {
    const url = new URL(path, BACKEND_ORIGIN).toString();

    const incomingCookie = event.request.headers.get('cookie') ?? '';

    const headers = new Headers(opts.headers ?? {});
    // preserve content-type if provided by caller, otherwise let it pass through
    if (!headers.has('cookie') && incomingCookie) headers.set('cookie', incomingCookie);

    // ensure Host not forwarded
    headers.delete('host');

    const init: RequestInit = {
        method: opts.method ?? event.request.method ?? 'GET',
        headers,
        body: opts.body ?? (await cloneRequestBodyIfNeeded(event)),
        redirect: 'manual',
        // note: fetch on the server will not include browser cookies automatically; we forward them above
    };

    const res = await fetch(url, init);
    return res;
}

async function cloneRequestBodyIfNeeded(event: RequestEvent) {
    const method = event.request.method?.toUpperCase();
    if (method === 'GET' || method === 'HEAD') return undefined;
    // clone body from original request
    const b = await event.request.arrayBuffer();
    return b.length ? new Uint8Array(b) : undefined;
}

