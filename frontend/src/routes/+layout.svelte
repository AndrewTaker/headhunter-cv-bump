<script lang="ts">
export let data;
</script>

<svelte:head>
    <meta name="viewport" content="width=device-width,initial-scale=1" />
</svelte:head>

<header class="container">
    <nav>
        {#if data.user}
            <ul>
                <li>
                    <hgroup>
                        <h2>{data.user.first_name}</h2>
                        <small>{data.user.id}</small>
                    </hgroup>
                </li>
            </ul>
            <ul><li><a href="/logout">Выйти</a></li></ul>
        {:else}
            <ul><li><a href="/login">Войти</a></li></ul>
        {/if}
    </nav>
</header>

<slot />

