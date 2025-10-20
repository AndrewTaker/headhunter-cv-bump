<script setup>
import { ref, onMounted } from 'vue';

const apiUrl = import.meta.env.VITE_API_URL
const data = ref(null);

</script>

<template>
    {{ apiUrl }}
</template>

<style scoped></style>
