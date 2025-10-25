import { ref, watchEffect, toValue } from 'vue'


export function useApi(url, creds = "omit") {
    const data = ref(null)
    const error = ref(null)

    const fetchData = () => {
        data.value = null
        error.value = null

        fetch(toValue(url), { credentials: creds })
            .then((res) => res.json())
            .then((json) => (data.value = json))
            .catch((err) => (error.value = err))
    }

    watchEffect(() => {
        fetchData()
    })

    return { data, error }
}
