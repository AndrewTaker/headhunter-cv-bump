<script lang="ts">
export let data;
async function toggle(id: string) {
    const res = await fetch(`/api/resumes/${id}/toggle`, { method: 'POST', credentials: 'include' });
    if (!res.ok) console.error('toggle failed');
    location.reload();
}
function formatDateTime(iso: string) {
    return new Date(iso).toLocaleString('ru-RU', {
        year: 'numeric',
        month: 'numeric',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
    });
}

</script>

{#each data.resumes as r}
    <article>
        <hgroup>
            <h3>{r.title}</h3>
            <small>{r.id}</small>
        </hgroup>
        <p>Создано: {formatDateTime(r.created_at)}</p>
        <p>Обновлено: {formatDateTime(r.updated_at)}</p>
        <footer>
            <button on:click={() => toggle(r.id)}>
                {#if r.is_scheduled === 1} Поставить в очередь {:else} Убрать из очереди {/if}
            </button>
        </footer>
    </article>
{/each}

