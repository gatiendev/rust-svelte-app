<script lang="ts">
  import '../app.css';
  import Header from '$lib/components/layout/Header.svelte';
  import Footer from '$lib/components/layout/Footer.svelte';
  import { theme } from '$lib/stores/theme';
</script>

<!-- Apply dark class to html is done in theme store, but we can also set it here for SSR? -->
<svelte:head>
  <!-- Optionally set theme class on server via onMount? Not needed because store handles it on client -->
</svelte:head>

<Header />
<main class="min-h-screen bg-gray-50 dark:bg-gray-950 text-gray-900 dark:text-gray-100">
  <slot />
</main>
<Footer />