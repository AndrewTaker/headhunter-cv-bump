import { ref, type Ref } from 'vue';
import { useUser } from './useUser';

const BASE_API_URL = 'http://localhost:44444';
const LOGIN_URL = `${BASE_API_URL}/auth/login`;
const LOGOUT_URL = `${BASE_API_URL}/auth/logout`;

const isAuthenticated: Ref<boolean> = ref(false);
const authLoading: Ref<boolean> = ref(false);
const authError: Ref<string | null> = ref(null);

const { user, isLoading: isProfileLoading, error: profileError, fetchUser } = useUser();

async function checkAuthStatus() {
    authLoading.value = true;
    authError.value = null;

    await fetchUser();

    if (user.value) {
        isAuthenticated.value = true;
    } else if (profileError.value && profileError.value.includes('401')) {
        isAuthenticated.value = false;
        authError.value = null;
    } else if (profileError.value) {
        authError.value = profileError.value;
    } else {
        isAuthenticated.value = false;
    }

    authLoading.value = false;
}

function logIn() {
    window.location.href = LOGIN_URL;
}

async function logOut() {
    try {
        const response = await fetch(LOGOUT_URL, {
            method: 'GET',
            credentials: 'include',
        });

        if (response.ok) {
            isAuthenticated.value = false;
            user.value = null;
            authError.value = null;
        } else {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Logout failed.');
        }
    } catch (e) {
        console.error('Logout error:', e);
        authError.value = e instanceof Error ? e.message : 'An unknown logout error occurred.';
    }
}

checkAuthStatus();


export function useAuth() {
    return {
        isAuthenticated,
        isProfileLoading,
        authError,
        user,
        checkAuthStatus,
        logIn,
        logOut,
    };
}
