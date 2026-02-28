<script lang="ts">
  export let padding = true;
  export let shadow = true;
  export let border = false;
</script>

<div
  class="bg-white dark:bg-gray-800 rounded-lg {padding ? 'p-6' : ''} {shadow ? 'shadow-md' : ''} {border ? 'border border-gray-200 dark:border-gray-700' : ''}"
>
  <slot />
</div>