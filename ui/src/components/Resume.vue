<script setup>
import { defineProps } from "vue"

const props = defineProps(["id", "title"])
</script>

<template>
    <h2>{{ props.title }}</h2>
    <small>{{ props.id }}</small>
</template>

<style></style>
