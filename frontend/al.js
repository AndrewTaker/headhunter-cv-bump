const BASE_URL = "http://localhost:44444";

async function genericFetch(method, url, data = null) {
    const options = {
        method: method,
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include'
    };

    if (data) { options.body = JSON.stringify(data); }

    try {
        const response = await fetch(BASE_URL + url, options);

        if (response.status === 401 || response.status === 403) {
            console.warn("Authentication failed or session expired.");
            return null;
        }

        if (!response.ok) {
            throw new Error(`Server error! Status: ${response.status}`);
        }
        return await response.json();
    } catch (error) {
        console.error(`${method} failed:`, error.message);
        return null;
    }
}

document.addEventListener('alpine:init', () => {
    Alpine.data('mainApp', () => ({
        user: null,
        resumes: [],
        loading: true,
        error: null,

        get isLoggedIn() {
            return this.user !== null;
        },

        async init() {
            this.loading = true;
            this.error = null;
            try {
                await this.loadInitialData();
            } catch (e) {
                this.error = "Не удалось загрузить начальные данные. Проверьте подключение к API.";
                console.error(e);
            }
            this.loading = false;
        },

        async loadInitialData() {
            const userData = await genericFetch("GET", "/me");
            this.user = userData;

            if (this.isLoggedIn) {
                const resumesData = await genericFetch("GET", "/resumes");
                if (resumesData && Array.isArray(resumesData.resumes)) {
                    this.resumes = resumesData.resumes;
                }
            } else {
                this.resumes = [];
            }
        },

        login() {
            window.location.href = `${BASE_URL}/auth/login`;
        },

        logout() {
            window.location.href = `${BASE_URL}/auth/logout`;
        },

        formatDate(dateString) {
            const date = new Date(dateString);
            const options = {
                year: "numeric", month: "numeric", day: "numeric",
                hour: "2-digit", minute: "2-digit", hour12: false,
            };
            return date.toLocaleString('ru-RU', options);
        }
    }));
});
