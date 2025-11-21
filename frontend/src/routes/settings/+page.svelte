<script lang="ts">
import { user } from '$lib/stores/user';
import { goto } from '$app/navigation';

async function deleteAccount() {
    if (!confirm('Действительно удалить все данные?')) return;

    const res = await fetch('/api/cleanup', { method: 'DELETE', credentials: 'include' });
    if (res.ok) {
        user.set(null);
        goto('/', { replaceState: true });
    } else {
        const json = await res.json().catch(() => ({ error: 'delete failed' }));
        alert('Delete failed: ' + (json.error ?? res.status));
    }
}
</script>

<h1>Найстройки аккаунта</h1>
<button on:click={deleteAccount}>Удалить мои данные</button>

