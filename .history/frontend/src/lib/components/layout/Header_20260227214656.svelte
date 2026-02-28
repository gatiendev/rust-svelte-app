<script lang="ts">
  import { theme, toggleTheme } from '$lib/stores/theme';
  import { auth } from '$lib/stores/auth';
  import { page } from '$app/stores';
  import { Moon, Sun, LogIn, LogOut, User } from 'lucide-svelte';
  import Button from '../ui/Button.svelte';

  // Derived logged-in state
  $: user = $auth.user;
</script>

<header class="bg-white dark:bg-gray-900 shadow-md sticky top-0 z-10">
  <nav class="container mx-auto px-4 py-3 flex items-center justify-between">
    <!-- Logo / App Name -->
    <a href="/" class="text-2xl font-bold text-gray-800 dark:text-white">
      SvelteKit Dark
    </a>

    <!-- Navigation Links (desktop) -->
    <div class="hidden md:flex space-x-6">
      <a href="/" class="text-gray-700 dark:text-gray-300 hover:text-primary-600 dark:hover:text-primary-400">
        Home
      </a>
      {#if !user}
        <a href="/auth/login" class="text-gray-700 dark:text-gray-300 hover:text-primary-600 dark:hover:text-primary-400">
          Login
        </a>
        <a href="/auth/register" class="text-gray-700 dark:text-gray-300 hover:text-primary-600 dark:hover:text-primary-400">
          Register
        </a>
      {:else}
        <a href="/profile" class="text-gray-700 dark:text-gray-300 hover:text-primary-600 dark:hover:text-primary-400">
          Profile
        </a>
        <button
          on:click={() => auth.logout()}
          class="text-gray-700 dark:text-gray-300 hover:text-primary-600 dark:hover:text-primary-400"
        >
          Logout
        </button>
      {/if}
    </div>

    <!-- Right side: dark mode toggle & user menu (mobile) -->
    <div class="flex items-center space-x-3">
      <!-- Dark Mode Toggle -->
      <button
        on:click={toggleTheme}
        class="p-2 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700 transition"
        aria-label="Toggle theme"
      >
        {#if $theme === 'dark'}
          <Sun class="w-5 h-5 text-yellow-500" />
        {:else}
          <Moon class="w-5 h-5 text-gray-700" />
        {/if}
      </button>

      <!-- Mobile menu button (simplified: just show user icon or login) -->
      <div class="md:hidden">
        {#if user}
          <button class="p-2 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700">
            <User class="w-5 h-5 text-gray-700 dark:text-gray-300" />
          </button>
        {:else}
          <a href="/auth/login" class="p-2 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700">
            <LogIn class="w-5 h-5 text-gray-700 dark:text-gray-300" />
          </a>
        {/if}
      </div>
    </div>
  </nav>
</header>