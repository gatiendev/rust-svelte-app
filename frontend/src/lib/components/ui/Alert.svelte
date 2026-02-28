<script lang="ts">
  import { CheckCircle, AlertCircle, Info, AlertTriangle, X } from 'lucide-svelte';

  export let type: 'success' | 'error' | 'info' | 'warning' = 'info';
  export let dismissible = false;
  export let title = '';

  const icons = {
    success: CheckCircle,
    error: AlertCircle,
    info: Info,
    warning: AlertTriangle
  };

  const colors = {
    success: {
      bg: 'bg-green-50 dark:bg-green-900/20',
      border: 'border-green-200 dark:border-green-800',
      text: 'text-green-800 dark:text-green-300',
      icon: 'text-green-600 dark:text-green-400'
    },
    error: {
      bg: 'bg-red-50 dark:bg-red-900/20',
      border: 'border-red-200 dark:border-red-800',
      text: 'text-red-800 dark:text-red-300',
      icon: 'text-red-600 dark:text-red-400'
    },
    info: {
      bg: 'bg-blue-50 dark:bg-blue-900/20',
      border: 'border-blue-200 dark:border-blue-800',
      text: 'text-blue-800 dark:text-blue-300',
      icon: 'text-blue-600 dark:text-blue-400'
    },
    warning: {
      bg: 'bg-yellow-50 dark:bg-yellow-900/20',
      border: 'border-yellow-200 dark:border-yellow-800',
      text: 'text-yellow-800 dark:text-yellow-300',
      icon: 'text-yellow-600 dark:text-yellow-400'
    }
  };

  const Icon = icons[type];
  let visible = true;

  function dismiss() {
    visible = false;
  }
</script>

{#if visible}
  <div class="rounded-xl border {colors[type].border} {colors[type].bg} p-4 flex items-start gap-3" role="alert">
    <div class="flex-shrink-0 {colors[type].icon}">
      <Icon class="w-5 h-5" />
    </div>
    <div class="flex-1">
      {#if title}
        <h3 class="text-sm font-medium {colors[type].text}">{title}</h3>
      {/if}
      <div class="text-sm {colors[type].text} opacity-90">
        <slot />
      </div>
    </div>
    {#if dismissible}
      <button
        on:click={dismiss}
        class="flex-shrink-0 {colors[type].text} hover:opacity-75 transition-opacity"
        aria-label="Dismiss"
      >
        <X class="w-4 h-4" />
      </button>
    {/if}
  </div>
{/if}