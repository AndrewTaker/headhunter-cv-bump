<script setup>
import { ref, onMounted } from 'vue';
import { useApi } from './composables/useApi.js'
import  Resume  from './components/Resume.vue'

const baseUrl = import.meta.env.VITE_API_URL
const url = `${baseUrl}/resumes`
const { data, error } = useApi(url)

</script>

<template>
    <Resume
        v-for="resume in data"
        :key="resume.ID"
        :id="resume.ID"
        :title="resume.Title"
    />
</template>

<style scoped></style>
