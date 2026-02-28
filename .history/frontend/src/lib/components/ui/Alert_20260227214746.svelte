<script lang="ts">
  import { AlertCircle, CheckCircle, Info, AlertTriangle, X } from 'lucide-svelte';
  import { createEventDispatcher } from 'svelte';

  export let type: 'success' | 'error' | 'info' | 'warning' = 'info';
  export let dismissible = false;

  const dispatch = createEventDispatcher();

  const icons = {
    success: CheckCircle,
    error: AlertCircle,
    info: Info,
    warning: AlertTriangle
  };

  const iconColors = {
    success: 'text-green-600 dark:text-green-400',
    error: 'text-red-600 dark:text-red-400',
    info: 'text-blue-600 dark:text-blue-400',
    warning: 'text-yellow-600 dark:text-yellow-400'
  };

  const bgColors = {
    success: 'bg-green-50 dark:bg-green-900/20',
    error: 'bg-red-50 dark:bg-red-900/20',
    info: 'bg-blue-50 dark:bg-blue-900/20',
    warning: 'bg-yellow-50 dark:bg-yellow-900/20'
  };

  const Icon = icons[type];
</script>

<div class="rounded-lg p-4 {bgColors[type]} flex items-start gap-3" role="alert">
  <div class="flex-shrink-0 {iconColors[type]}">
    <Icon class="w-5 h-5" />
  </div>
  <div class="flex-1 text-sm text-gray-800 dark:text-gray-200">
    <slot />
  </div>
  {#if dismissible}
    <button
      on:click={() => dispatch('dismiss')}
      class="flex-shrink-0 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
      aria-label="Dismiss"
    >
      <X class="w-4 h-4" />
    </button>
  {/if}
</div>