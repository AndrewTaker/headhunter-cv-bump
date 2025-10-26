import { ref, type Ref } from 'vue';
import type { ErrorResponse } from '../types/types';

const BASE_API_URL = 'http://localhost:44444';

export function useApiFetch<T>(endpoint: string, options: RequestInit = {}) {
    const data: Ref<T | null> = ref(null);
    const isLoading: Ref<boolean> = ref(false);
    const error: Ref<string | null> = ref(null);

    const fullUrl = `${BASE_API_URL}${endpoint}`;

    const fetchApi = async () => {
        isLoading.value = true;
        error.value = null;
        data.value = null;

        try {
            const response = await fetch(fullUrl, {
                ...options,
                credentials: 'include',
            });

            if (response.status === 401) {
                throw new Error('Unauthorized. Session cookie "sess" is invalid or missing.');
            }

            if (!response.ok) {
                const errorData: ErrorResponse = await response.json();
                throw new Error(errorData.error || `HTTP error! Status: ${response.status}`);
            }

            if (response.status !== 204) {
                data.value = await response.json() as T;
            } else {
                data.value = null;
            }
        } catch (e) {
            console.error('API Fetch Failed:', e);
            error.value = e instanceof Error ? e.message : 'An unknown API error occurred.';
        } finally {
            isLoading.value = false;
        }
    };

    return { data, isLoading, error, fetchApi };
}
