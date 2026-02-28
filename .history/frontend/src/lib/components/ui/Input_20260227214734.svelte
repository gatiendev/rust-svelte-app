<script lang="ts">
  export let label = '';
  export let type: 'text' | 'email' | 'password' | 'number' = 'text';
  export let value = '';
  export let placeholder = '';
  export let error: string | undefined = undefined;
  export let disabled = false;
  export let id = `input-${Math.random().toString(36).substring(2, 9)}`;
</script>

<div class="mb-4">
  {#if label}
    <label for={id} class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
      {label}
    </label>
  {/if}
  <input
    {id}
    {type}
    {placeholder}
    bind:value
    {disabled}
    class="w-full px-3 py-2 border rounded-lg bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 transition {error
      ? 'border-red-500 dark:border-red-400'
      : 'border-gray-300 dark:border-gray-700'}"
    aria-invalid={!!error}
    aria-describedby={error ? `${id}-error` : undefined}
  />
  {#if error}
    <p id="{id}-error" class="mt-1 text-sm text-red-600 dark:text-red-400">
      {error}
    </p>
  {/if}
</div>