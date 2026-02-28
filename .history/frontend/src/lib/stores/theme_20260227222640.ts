import { writable } from 'svelte/store';
import { browser } from '$app/environment';

type Theme = 'light' | 'dark';

// Get initial theme from localStorage or default to 'dark'
const storedTheme = browser ? (localStorage.getItem('theme') as Theme) : 'dark';
const initialTheme: Theme = storedTheme === 'light' ? 'light' : 'dark';

export const theme = writable<Theme>(initialTheme);

// Subscribe to changes and update localStorage + html class
if (browser) {
    theme.subscribe((value) => {
        console.log('Theme changed to:', value); // Add this
        localStorage.setItem('theme', value);
        const root = document.documentElement;
        if (value === 'dark') {
            root.classList.add('dark');
        } else {
            root.classList.remove('dark');
        }
    });
}

// Toggle function for convenience
export function toggleTheme() {
    theme.update((t) => (t === 'light' ? 'dark' : 'light'));
}