import { computed, type Ref } from 'vue';
import { useApiFetch } from './useApiFetch';
import type { Resume } from '../types/types';

export function useResume(resumeId: string) {
    const endpoint = `/resumes/${resumeId}`;
    const { data, isLoading, error, fetchApi } = useApiFetch<Resume>(endpoint);

    const resume: Ref<Resume | null> = computed(() => {
        if (!data.value) return null;

        return {
            ...data.value,
            created_at: new Date(data.value.created_at as unknown as string),
            updated_at: new Date(data.value.updated_at as unknown as string),
        };
    });

    return {
        resume,
        isLoading,
        error,
        fetchResume: fetchApi,
    };
}
