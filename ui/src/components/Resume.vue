<script setup>
const props = defineProps(["id", "title", "created_at", "updated_at", "is_scheduled"])
</script>

<template>
    <h2>{{ props.title }}</h2>
    <small>{{ props.id }}</small>
    <p>{{ props.created_at }} / {{ props.updated_at }}</p>
    <p>{{ props.is_scheduled }}</p>
</template>

<style></style>
