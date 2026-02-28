import { defineConfig } from 'vite';
import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  plugins: [
    tailwindcss({
      config: {
        darkMode: 'class', // essential for class-based dark mode
      },
    }),
    sveltekit(),
  ],
});