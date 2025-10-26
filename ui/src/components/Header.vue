<script setup lang="ts">
import type { User } from "../types/types"
import { useAuth } from '../composables/useAuth';

const user = defineProps<User>()
const {
    logIn,
    logOut
} = useAuth();

</script>
<template>
    <header>
        <nav>
            <span v-if="user">{{ user.first_name }} {{ user.last_name }} ({{ user.id }})</span>
            <span v-else>Guest</span>

            <button v-if="user" @click="logOut" class="btn-logout">
                Log Out
            </button>

            <section v-else class="login-view">
                <h2>🔒 Please Log In</h2>
                <p>Access to your resume data requires authentication via OAuth2.</p>
                <button @click="logIn" class="btn-login">
                    Start OAuth2 Login
                </button>
            </section>

        </nav>
    </header>
</template>
