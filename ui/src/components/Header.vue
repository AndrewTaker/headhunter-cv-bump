<script setup>
const props = defineProps(["user"])
</script>

<template>
    <div v-if="user">
        <h2 >{{ props.user.first_name }}</h2>
        <small>{{ props.user.id }}</small>
    </div>
    <h2 v-else>Not authorized</h2>

    <ul>
        <li v-if="user"><a href="http://localhost:44444/auth/logout">logout</a></li>
        <li v-else><a href="http://localhost:44444/auth/login">login</a></li>
    </ul>
</template>

<style></style>
