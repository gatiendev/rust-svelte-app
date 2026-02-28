<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { auth } from '$lib/stores/auth';
  import { Rocket, Shield, Zap } from 'lucide-svelte';
</script>

<div class="container mx-auto px-4 py-12">
  <!-- Hero Section -->
  <section class="text-center mb-16">
   <h1 class="text-5xl font-extrabold mb-4 text-blue-600 dark:text-blue-400">
  Welcome to SvelteKit Dark Starter
</h1>
    <p class="text-xl text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
      A modern starter template with dark mode, modular components, and authentication UI.
    </p>
    <div class="mt-8 flex justify-center gap-4">
      <Button variant="primary" size="lg" href="/auth/register">
        Get Started
      </Button>
      <Button variant="outline" size="lg" href="https://github.com" target="_blank">
        GitHub
      </Button>
    </div>
  </section>

  <!-- Features Grid -->
  <section class="grid md:grid-cols-3 gap-6">
    <Card>
      <div class="flex flex-col items-center text-center">
        <Rocket class="w-12 h-12 text-blue-600 dark:text-blue-400 mb-4" />
        <h2 class="text-2xl font-bold mb-2">Fast & Light</h2>
        <p class="text-gray-600 dark:text-gray-400">
          Built on SvelteKit for optimal performance and developer experience.
        </p>
      </div>
    </Card>

    <Card>
      <div class="flex flex-col items-center text-center">
        <Shield class="w-12 h-12 text-green-600 dark:text-green-400 mb-4" />
        <h2 class="text-2xl font-bold mb-2">Modular Components</h2>
        <p class="text-gray-600 dark:text-gray-400">
          Reusable, testable UI components with Tailwind styling.
        </p>
      </div>
    </Card>

    <Card>
      <div class="flex flex-col items-center text-center">
        <Zap class="w-12 h-12 text-yellow-600 dark:text-yellow-400 mb-4" />
        <h2 class="text-2xl font-bold mb-2">Auth Ready</h2>
        <p class="text-gray-600 dark:text-gray-400">
          Mock authentication with login/register forms – easy to connect to real backend.
        </p>
      </div>
    </Card>
  </section>

  <!-- Call to action if not logged in -->
  {#if !$auth.user}
    <section class="mt-16 text-center">
      <Card padding={false} border class="p-8">
        <h2 class="text-3xl font-bold mb-4">Ready to start your project?</h2>
        <p class="mb-6 text-gray-600 dark:text-gray-400">
          Clone this template and build something amazing.
        </p>
        <Button variant="primary" size="lg" href="/auth/register">
          Sign Up Now
        </Button>
      </Card>
    </section>
  {/if}
</div>