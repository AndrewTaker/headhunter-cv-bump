import { computed, type Ref } from 'vue';
import { useApiFetch } from './useApiFetch';
import type { Resume, ResumeResponseMany } from '@/types/types';

export function useResumes() {
    const endpoint = `/resumes`;
    const { data, isLoading, error, fetchApi } = useApiFetch<ResumeResponseMany>(endpoint);

    const resumes: Ref<Resume[]> = computed(() => {
        if (!data.value) return [];
        return data.value.resumes.map(r => ({
            ...r,
            created_at: new Date(r.created_at as unknown as string),
            updated_at: new Date(r.updated_at as unknown as string),
        }));
    });

    return {
        resumes,
        isLoading,
        error,
        fetchResumes: fetchApi,
    };
}
