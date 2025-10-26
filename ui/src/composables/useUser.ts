import { computed, type Ref } from 'vue';
import { useApiFetch } from './useApiFetch';
import type { User } from '@/types/types';

export function useUser() {
    const endpoint = `/me`;
    const { data, isLoading, error, fetchApi } = useApiFetch<User>(endpoint);
    const user: Ref<User | null> = computed(() => data.value);

    return {
        user,
        isLoading,
        error,
        fetchUser: fetchApi,
    };
}
