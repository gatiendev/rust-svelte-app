<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { auth } from '$lib/stores/auth';
  import { Rocket, Shield, Zap, ArrowRight } from 'lucide-svelte';
  import { fade, fly } from 'svelte/transition';
  import { quintOut } from 'svelte/easing';
</script>

<div class="container mx-auto px-4 py-16 md:py-24">
  <!-- Hero Section -->
  <section class="text-center mb-20" in:fade={{ duration: 600 }}>
    <h1 class="text-5xl md:text-6xl font-extrabold mb-6 bg-gradient-to-r from-primary-600 to-primary-400 bg-clip-text text-transparent">
      Build Faster with SvelteKit
    </h1>
    <p class="text-xl text-gray-600 dark:text-gray-400 max-w-2xl mx-auto leading-relaxed">
      A modern starter template with dark mode, modular components, and authentication UI – designed for professional projects.
    </p>
    <div class="mt-10 flex flex-col sm:flex-row justify-center gap-4">
      <Button variant="primary" size="lg" href="/auth/register" class="group">
        Get Started
        <ArrowRight class="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" />
      </Button>
      <Button variant="outline" size="lg" href="https://github.com" target="_blank">
        GitHub
      </Button>
    </div>
  </section>

  <!-- Features Grid -->
  <section class="grid md:grid-cols-3 gap-8">
    {#each [
      { icon: Rocket, title: 'Fast & Light', desc: 'Built on SvelteKit for optimal performance and developer experience.', color: 'text-primary-500' },
      { icon: Shield, title: 'Modular Components', desc: 'Reusable, testable UI components with Tailwind styling.', color: 'text-primary-500' },
      { icon: Zap, title: 'Auth Ready', desc: 'Mock authentication with login/register forms – easy to connect to real backend.', color: 'text-primary-500' }
    ] as feature, i}
      <div in:fly={{ y: 30, duration: 500, delay: i * 100, easing: quintOut }}>
        <Card hover class="h-full flex flex-col items-center text-center p-8">
          <div class="w-16 h-16 rounded-full bg-primary-50 dark:bg-primary-900/20 flex items-center justify-center mb-5">
            <svelte:component this={feature.icon} class="w-8 h-8 {feature.color}" />
          </div>
          <h2 class="text-2xl font-bold mb-3 text-gray-900 dark:text-white">{feature.title}</h2>
          <p class="text-gray-600 dark:text-gray-400 leading-relaxed">{feature.desc}</p>
        </Card>
      </div>
    {/each}
  </section>

  <!-- Call to Action -->
  {#if !$auth.user}
    <section class="mt-20 text-center" in:fade={{ duration: 600, delay: 300 }}>
      <Card class="p-10 max-w-2xl mx-auto border-2 border-primary-200 dark:border-primary-800">
        <h2 class="text-3xl font-bold mb-4 text-gray-900 dark:text-white">Ready to start your project?</h2>
        <p class="mb-8 text-gray-600 dark:text-gray-400 text-lg">
          Clone this template and build something amazing today.
        </p>
        <Button variant="primary" size="lg" href="/auth/register">
          Sign Up Now
        </Button>
      </Card>
    </section>
  {/if}
</div>