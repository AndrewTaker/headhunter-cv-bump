<script setup lang="ts">
import { onMounted, watch } from 'vue';
import { useAuth } from './composables/useAuth';
import { useResumes } from './composables/useResumes';
import Header from './components/Header.vue'
import Resume from './components/Resume.vue'

const {
    isAuthenticated,
    isProfileLoading,
    authError,
    user,
    logIn,
    logOut
} = useAuth();

const {
    resumes,
    isLoading: isResumesLoading,
    error: resumesError,
    fetchResumes
} = useResumes();

onMounted(() => {
    if (isAuthenticated.value) {
        fetchResumes();
    }
});

watch(isAuthenticated, (isAuth) => {
    if (isAuth) {
        fetchResumes();
    }
});
</script>

<template>
    <div id="app-container">
        <Header v-bind="user"/>

        <main>
            <div v-if="isProfileLoading" class="status-box loading">
                Authenticating session... Please wait.
            </div>

            <div v-else-if="authError" class="status-box error">
                🚨 Authentication Error: {{ authError }}. Please try logging in.
            </div>

            <section v-else-if="isAuthenticated" class="dashboard">
                <h2>📊 Your Dashboard</h2>

                <div v-if="isResumesLoading" class="status-box loading">
                    Fetching your resumes...
                </div>

                <div v-else-if="resumesError" class="status-box error">
                    ⚠️ **Error fetching resumes:** {{ resumesError }}
                </div>

                <div
                    v-else-if="resumes.length"
                    v-for="resume in resumes" :key="resume.id"
                    class="resume-list"
                >
                    <Resume v-bind="resume"/>
                </div>

                <div v-else class="status-box empty">
                    You have no resumes yet.
                </div>
            </section>

        </main>
    </div>
</template>

<style scoped>
#app-container {
    max-width: 800px;
    margin: 0 auto;
    font-family: sans-serif;
    padding: 20px;
}
header {
    background: #f0f0f0;
    padding: 15px;
    margin-bottom: 20px;
    border-radius: 4px;
}
nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
}
main {
    padding: 20px;
    border: 1px solid #ddd;
    border-radius: 4px;
}

.status-box {
    padding: 10px;
    border-radius: 4px;
    margin-bottom: 15px;
}
.loading {
    background-color: #e6f7ff;
    color: #1890ff;
}
.error {
    background-color: #fff1f0;
    color: #ff4d4f;
    border: 1px solid #ff4d4f;
    font-weight: bold;
}
.empty {
    background-color: #fafafa;
    color: #595959;
}

/* Buttons */
.btn-login {
    padding: 10px 20px;
    background-color: #007bff;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1.1em;
}
.btn-logout {
    padding: 8px 15px;
    background-color: #dc3545;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
}
.resume-list {
    list-style: disc;
    padding-left: 20px;
}
.resume-list li {
    margin-bottom: 5px;
}
</style>
