const BASE_URL = "http://localhost:44444"

function login() {
    window.location.href = `${BASE_URL}/auth/login`
}

function logout() {
    window.location.href = `${BASE_URL}/auth/logout`
}

async function loadPage() {
    const user = await genericFetch("GET", "/me")
    if (user) {
        fillAuthSection(true)
        fillUserSection(user)

        const resumes = await genericFetch("GET", "/resumes")
        fillResumesSection(resumes.resumes)
    } else {
        fillAuthSection(false)
    }
}

function fillAuthSection(isLoggedIn) {
    const auth = document.getElementById("auth")
    const button = document.createElement("a")
    button.role = "button"

    if (isLoggedIn) {
        button.href = `${BASE_URL}/auth/logout`
        button.textContent = "logout"
    } else {
        button.href = `${BASE_URL}/auth/login`
        button.textContent = "login"
    }

    auth.appendChild(button)
}

function fillUserSection(user) {
    if (user) {
        const profile = document.getElementById("profile")

        const hgroup = document.createElement("hgroup")
        const heading = document.createElement("h1")
        const small = document.createElement("small")

        const h1Content = document.createTextNode(`${user.last_name} ${user.first_name} ${user.middle_name}`)
        const smallContent = document.createTextNode(user.id)
        heading.appendChild(h1Content)
        small.appendChild(smallContent)

        hgroup.appendChild(heading)
        hgroup.appendChild(small)
        profile.appendChild(hgroup)
    }
}

function fillResumesSection(resumes) {
    if (resumes) {
        const resumeSection = document.getElementById("resumes")

        resumes.forEach((resume) => {
            const article = document.createElement("article")

            const header = document.createElement("header")
            const hgroup = document.createElement("hgroup")
            const heading = document.createElement("h2")
            const small = document.createElement("small")

            const footer = document.createElement("div")
            const ca = document.createElement("p")
            const ua = document.createElement("p")

            heading.textContent = resume.title
            small.textContent = resume.id

            const createdAt = new Date(resume.created_at)
            const updatedAt = new Date(resume.updated_at)
            const localFormat = {
                year: "numeric",
                month: "numeric",
                day: "numeric",
                hour: "2-digit",
                minute: "2-digit",
                hour12: false,
            }
            ca.textContent = `Создано: ${createdAt.toLocaleString(undefined, localFormat)}`
            ua.textContent = `Обновлено: ${updatedAt.toLocaleString(undefined, localFormat)}`

            hgroup.appendChild(heading)
            hgroup.appendChild(small)
            header.appendChild(hgroup)
            footer.appendChild(ca)
            footer.appendChild(ua)
            article.appendChild(header)
            article.appendChild(footer)

            resumeSection.appendChild(article)
        })
    }
}


async function genericFetch(method, url, headers = null, data = null) {
    const options = { method: method };

    if (headers) { options.headers = headers }
    if (data) { options.body = JSON.stringify(data) }

    try {
        const response = await fetch(url, options);

        if (!response.ok) {
            throw new Error(`Server error! Status: ${response.status}`);
        }

        const responseData = await response.json();
        return responseData;

    } catch (error) {
        console.error(`${method} failed:`, error.message);
        return null;
    }
}

(async () => { await loadPage() })();
