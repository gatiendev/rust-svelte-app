import { writable } from 'svelte/store';
import { browser } from '$app/environment';

type Theme = 'light' | 'dark';

function createThemeStore() {
    // Determine initial theme once
    const getInitialTheme = (): Theme => {
        if (!browser) return 'dark'; // SSR default
        const stored = localStorage.getItem('theme') as Theme | null;
        return stored === 'light' ? 'light' : 'dark';
    };

    const initialTheme = getInitialTheme();
    const { subscribe, set, update } = writable<Theme>(initialTheme);

    // Apply theme to DOM on every change (client-side only)
    if (browser) {
        // Apply the initial theme immediately
        const root = document.documentElement;
        if (initialTheme === 'dark') {
            root.classList.add('dark');
        } else {
            root.classList.remove('dark');
        }

        // Subscribe to future changes
        subscribe((value) => {
            console.log('Theme changed to:', value);
            localStorage.setItem('theme', value);
            if (value === 'dark') {
                root.classList.add('dark');
            } else {
                root.classList.remove('dark');
            }
        });
    }

    return {
        subscribe,
        toggle: () => update(t => (t === 'light' ? 'dark' : 'light')),
        set,
    };
}

export const theme = createThemeStore();
export const toggleTheme = () => theme.toggle();