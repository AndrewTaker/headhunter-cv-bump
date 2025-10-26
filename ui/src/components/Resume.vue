<script setup lang="ts">
import { reactive, computed } from 'vue';
import type { Resume } from '../types/types';

function toLocalDateTime(date: Date): string {
    const fDate: string = new Intl.DateTimeFormat('ru-RU', {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit"
    }).format(date)
    return fDate

}

const resume = defineProps<Resume>()
const updatedAt = computed(() => {
    const date: Date = new Date(resume.updated_at)
    return toLocalDateTime(date)
})
const createdAt = computed(() => {
    const date: Date = new Date(resume.created_at)
    return toLocalDateTime(date)
})

</script>

<template>
    <div v-if="resume">
        <h2>{{ resume.title }}</h2>
        <p><b>Updated:</b> {{ updatedAt }} <b>Created:</b> {{ createdAt }}</p>
        <p>
            Scheduled: {{ resume.is_scheduled === 1 ? 'Yes' : 'No' }}
            [<a :href="resume.alternate_url" target="_blank">Link</a>]
        </p>
    </div>
    <div v-else>
        <p>No resume data available.</p>
    </div>
</template>
