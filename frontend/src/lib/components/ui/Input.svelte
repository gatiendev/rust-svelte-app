<script lang="ts">
  export let label = '';
  export let type: 'text' | 'email' | 'password' | 'number' = 'text';
  export let value = '';
  export let placeholder = '';
  export let error: string | undefined = undefined;
  export let disabled = false;
  export let id = `input-${Math.random().toString(36).substring(2, 9)}`;
  export let icon = undefined; // can pass a component like <User />
</script>

<div class="mb-4">
  {#if label}
    <label for={id} class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
      {label}
    </label>
  {/if}
  <div class="relative">
    {#if icon}
      <div class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500 dark:text-gray-400">
        <svelte:component this={icon} class="w-4 h-4" />
      </div>
    {/if}
    <input
      {id}
      {type}
      {placeholder}
      bind:value
      {disabled}
      class="w-full px-4 py-2.5 border rounded-lg bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 dark:focus:ring-primary-400 transition duration-200 {icon ? 'pl-10' : ''} {error
        ? 'border-red-500 dark:border-red-400 focus:ring-red-500'
        : 'border-gray-300 dark:border-gray-700'}"
      aria-invalid={!!error}
      aria-describedby={error ? `${id}-error` : undefined}
    />
  </div>
  {#if error}
    <p id="{id}-error" class="mt-1.5 text-sm text-red-600 dark:text-red-400">
      {error}
    </p>
  {/if}
</div>