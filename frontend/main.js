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
        fillAuthSection(true, user)
        fillUserSection(user)

        const resumes = await genericFetch("GET", "/resumes")
        fillResumesSection(resumes.resumes)
    } else {
        fillAuthSection(false)
    }
}

function fillAuthSection(isLoggedIn, user = null) {
    const loginUrl = `${BASE_URL}/auth/login`
    const logoutUrl = `${BASE_URL}/auth/logout`

    const nav = document.getElementById("navigation")
    const identity = document.createElement("ul")
    const links = document.createElement("ul")
    const auth = document.createElement("a")

    const linksItems = []

    if (isLoggedIn) {
        auth.href = logoutUrl
        auth.textContent = "Выйти"
    } else {
        auth.href = loginUrl
        auth.textContent = "Войти"
    }

    linksItems.push(auth)
    links.append(linksItems)
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
    const resumeSection = document.getElementById("resumes")

    if (resumes) {
        resumes.forEach((resume) => {
            const article = document.createElement("article")

            const header = document.createElement("header")
            const hgroup = document.createElement("hgroup")
            const heading = document.createElement("h2")
            const small = document.createElement("small")

            const footer = document.createElement("div")
            const ca = document.createElement("p")
            const ua = document.createElement("p")
            const scheduled = document.createElement("p")

            heading.textContent = resume.title
            small.textContent = resume.id

            const setUpScheduleAnchor = document.createElement("a")
            setUpScheduleAnchor.href = `${BASE_URL}/resumes/${resume.id}/toggle`
            setUpScheduleAnchor.textContent = "toggle"
            scheduled.textContent = resume.is_scheduled == 0 ? "Не стоит в очереди" : "Стоит в очереди"

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
            footer.appendChild(scheduled)
            footer.appendChild(setUpScheduleAnchor)
            article.appendChild(header)
            article.appendChild(footer)

            resumeSection.appendChild(article)
        })
    } else {
        const p = document.createElement("p")
        p.textContent = "В базе нет резюме"
        const a = document.createElement("a")
        a.href =
            resumeSection.appendChild(p)
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
