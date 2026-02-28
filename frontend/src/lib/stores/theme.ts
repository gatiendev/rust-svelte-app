import { writable } from 'svelte/store';
import { browser } from '$app/environment';

type Theme = 'light' | 'dark';

function createThemeStore() {
    // Read initial value from localStorage, default to 'dark'
    const stored = browser ? (localStorage.getItem('theme') as Theme) : 'dark';
    const initial = stored === 'light' ? 'light' : 'dark';

    const { subscribe, set, update } = writable<Theme>(initial);

    // This will run whenever the store changes
    if (browser) {
        subscribe((value) => {
            console.log('Theme changed to:', value);
            localStorage.setItem('theme', value);
            document.documentElement.classList.toggle('dark', value === 'dark');
        });
    }

    return {
        subscribe,
        toggle: () => {
            console.log('toggle function called at:', Date.now());
            update(t => t === 'light' ? 'dark' : 'light');
        },
        set,
    };
}

export const theme = createThemeStore();
export const toggleTheme = () => theme.toggle();