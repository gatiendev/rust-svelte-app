<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { fade, fly } from 'svelte/transition';
  import { quintOut } from 'svelte/easing';
  import Button from './Button.svelte';
  import { X } from 'lucide-svelte';

  export let open = false;
  export let title = 'Modal';
  export let closeOnOutsideClick = true;
  export let size: 'sm' | 'md' | 'lg' = 'md';

  const dispatch = createEventDispatcher();

  const sizeClasses = {
    sm: 'max-w-md',
    md: 'max-w-lg',
    lg: 'max-w-2xl'
  };

  function close() {
    open = false;
    dispatch('close');
  }

  function handleBackdropClick(e: MouseEvent) {
    if (closeOnOutsideClick && e.target === e.currentTarget) {
      close();
    }
  }
</script>

{#if open}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm"
    on:click={handleBackdropClick}
    transition:fade={{ duration: 200 }}
  >
    <div
      class="bg-white dark:bg-gray-800 rounded-xl shadow-2xl {sizeClasses[size]} w-full max-h-[90vh] overflow-auto"
      transition:fly={{ y: 20, duration: 300, easing: quintOut }}
    >
      <div class="flex justify-between items-center p-5 border-b border-gray-200 dark:border-gray-700">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white">{title}</h2>
        <Button variant="ghost" size="sm" on:click={close} aria-label="Close modal" class="!p-1.5">
          <X class="w-5 h-5" />
        </Button>
      </div>
      <div class="p-5">
        <slot />
      </div>
    </div>
  </div>
{/if}