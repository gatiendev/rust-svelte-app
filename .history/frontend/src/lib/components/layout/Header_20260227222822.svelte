<script lang="ts">
  import { theme, toggleTheme } from '$lib/stores/theme';
  import { auth } from '$lib/stores/auth';
  import { Moon, Sun, LogIn, LogOut, User, Menu } from 'lucide-svelte';
  import Button from '../ui/Button.svelte';

  let mobileMenuOpen = false;
</script>

<header class="sticky top-0 z-30 backdrop-blur-md bg-white/80 dark:bg-gray-900/80 border-b border-gray-200 dark:border-gray-800">
  <nav class="container mx-auto px-4 py-3 flex items-center justify-between">
    <!-- Logo -->
    <a href="/" class="text-2xl font-bold bg-gradient-to-r from-primary-600 to-primary-400 bg-clip-text text-transparent">
      SvelteKit Dark
    </a>

    <!-- Desktop Navigation -->
    <div class="hidden md:flex items-center space-x-8">
      <a href="/" class="text-gray-700 dark:text-gray-300 hover:text-primary-600 dark:hover:text-primary-400 transition-colors font-medium">
        Home
      </a>
      {#if !$auth.user}
        <a href="/auth/login" class="text-gray-700 dark:text-gray-300 hover:text-primary-600 dark:hover:text-primary-400 transition-colors font-medium">
          Login
        </a>
        <a href="/auth/register" class="text-gray-700 dark:text-gray-300 hover:text-primary-600 dark:hover:text-primary-400 transition-colors font-medium">
          Register
        </a>
      {:else}
        <a href="/profile" class="text-gray-700 dark:text-gray-300 hover:text-primary-600 dark:hover:text-primary-400 transition-colors font-medium">
          Profile
        </a>
        <button
          on:click={() => auth.logout()}
          class="text-gray-700 dark:text-gray-300 hover:text-primary-600 dark:hover:text-primary-400 transition-colors font-medium"
        >
          Logout
        </button>
      {/if}
    </div>

    <!-- Right side controls -->
    <div class="flex items-center gap-2">
      <Button variant="ghost" size="sm" on:click={toggleTheme} class="!p-2" aria-label="Toggle theme">
  {#if $theme === 'dark'}
    <Sun class="w-5 h-5 text-yellow-500" />
  {:else}
    <Moon class="w-5 h-5 text-gray-700 dark:text-gray-300" />
  {/if}
</Button>

      <!-- Mobile menu button -->
      <Button variant="ghost" size="sm" on:click={() => mobileMenuOpen = !mobileMenuOpen} class="md:hidden !p-2">
        <Menu class="w-5 h-5 text-gray-700 dark:text-gray-300" />
      </Button>
    </div>
  </nav>

  <!-- Mobile menu -->
  {#if mobileMenuOpen}
    <div class="md:hidden border-t border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 p-4">
      <div class="flex flex-col space-y-3">
        <a href="/" class="text-gray-700 dark:text-gray-300 hover:text-primary-600 py-2" on:click={() => mobileMenuOpen = false}>Home</a>
        {#if !$auth.user}
          <a href="/auth/login" class="text-gray-700 dark:text-gray-300 hover:text-primary-600 py-2" on:click={() => mobileMenuOpen = false}>Login</a>
          <a href="/auth/register" class="text-gray-700 dark:text-gray-300 hover:text-primary-600 py-2" on:click={() => mobileMenuOpen = false}>Register</a>
        {:else}
          <a href="/profile" class="text-gray-700 dark:text-gray-300 hover:text-primary-600 py-2" on:click={() => mobileMenuOpen = false}>Profile</a>
          <button
            on:click={() => { auth.logout(); mobileMenuOpen = false; }}
            class="text-left text-gray-700 dark:text-gray-300 hover:text-primary-600 py-2"
          >
            Logout
          </button>
        {/if}
      </div>
    </div>
  {/if}
</header>