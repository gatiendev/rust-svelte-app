<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { fade } from 'svelte/transition';
  import Button from './Button.svelte';
  import { X } from 'lucide-svelte';

  export let open = false;
  export let title = 'Modal';
  export let closeOnOutsideClick = true;

  const dispatch = createEventDispatcher();

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
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 dark:bg-black/70"
    on:click={handleBackdropClick}
    transition:fade={{ duration: 200 }}
  >
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full max-h-[90vh] overflow-auto">
      <div class="flex justify-between items-center p-4 border-b border-gray-200 dark:border-gray-700">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{title}</h2>
        <Button variant="outline" size="sm" on:click={close} aria-label="Close modal">
          <X class="w-4 h-4" />
        </Button>
      </div>
      <div class="p-4">
        <slot />
      </div>
    </div>
  </div>
{/if}