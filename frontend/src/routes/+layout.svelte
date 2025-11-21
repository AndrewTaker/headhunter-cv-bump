<script lang="ts">
export let data;
</script>

<svelte:head>
    <meta name="viewport" content="width=device-width,initial-scale=1" />
</svelte:head>
{#if data.user}
    <header class="container">
        <nav>
            <ul>
                <li><h2>{data.user.first_name}</h2></li>
            </ul>
            <ul>
                <li><a href="/" class="contrast">🏠</a></li>
                <li><a href="/settings" class="contrast">Настройки</a></li>
                <li><a href="/resumes" class="contrast">Резюме</a></li>
                <li><a href="/logout">Выйти</a></li>
            </ul>
        </nav>
    </header>
{:else}
    <header>
        <nav>
            <ul><li><a href="/login">Войти</a></li></ul>
        </nav>
    </header>
{/if}

<style>
    header { border-bottom: 1px solid white; margin-bottom: 1em; }
</style>
<!-- {@render children?.()} -->
<slot />
