<script setup>
import { ref } from 'vue';
import { watch } from 'vue';
import { useApi } from './composables/useApi.js'
import Resume from './components/Resume.vue'
import Header from './components/Header.vue'

const baseUrl = import.meta.env.VITE_API_URL
const { data, error } = useApi(`${baseUrl}/resumes`, "include");
const { data: meData, error: meError } = useApi(`${baseUrl}/me`, "include");
const statusMessage = ref('Fetching data...');

watch([data, meData, error, meError], () => {
    if (data.value !== null && meData.value !== null) {
        statusMessage.value = 'All data loaded successfully!';
    } else if (error.value || meError.value) {
        statusMessage.value = 'An error occurred during fetch.';
    }
});
</script>
<template>
    <Header :user="meData" />

    <div v-if="isLoading" class="loading-state">
        <p>Loading full page content...</p>

        <p>Status: {{ statusMessage }}</p>
    </div>

    <div v-else class="content-loaded">
        <p>Status: **{{ statusMessage }}**</p>

        <hr>

        <h3>Me Data</h3>
        <p v-if="meError">Error fetching user: {{ meError.message || meError }}</p>
        <pre v-else>{{ meData }}</pre>

        <hr>

        <h3>Resumes</h3>
        <template v-if="data">
            <Resume
                v-for="resume in data.resumes"
                :key="resume.id"
                :id="resume.id"
                :title="resume.title"
                :created_at="resume.created_at"
                :updated_at="resume.updated_at"
                :is_scheduled="resume.is_scheduled"
            />
        </template>
        <p v-else-if="error">Error loading resumes: {{ error.message || error }}</p>
        <p v-else>No resumes available.</p>
    </div>
</template>
