<script lang="ts">
  export let padding = true;
  export let shadow = 'md'; // 'none', 'sm', 'md', 'lg'
  export let border = false;
  export let hover = false; // if true, adds hover scale and shadow
  export let interactive = false; // if true, adds cursor-pointer

  const shadowClasses = {
    none: '',
    sm: 'shadow-sm',
    md: 'shadow-md',
    lg: 'shadow-lg'
  };

  const hoverClasses = hover ? 'transition-all duration-200 hover:shadow-xl hover:-translate-y-1' : '';
  const interactiveClasses = interactive ? 'cursor-pointer' : '';
</script>

<div
  class="bg-white dark:bg-gray-800 rounded-xl {padding ? 'p-6' : ''} {shadowClasses[shadow]} {border ? 'border border-gray-200 dark:border-gray-700' : ''} {hoverClasses} {interactiveClasses}"
>
  <slot />
</div>