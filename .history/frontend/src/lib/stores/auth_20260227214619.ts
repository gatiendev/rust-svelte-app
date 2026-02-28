import { writable } from 'svelte/store';
import { goto } from '$app/navigation';

export interface User {
    email: string;
    name?: string;
}

interface AuthState {
    user: User | null;
    loading: boolean;
    error: string | null;
}

function createAuthStore() {
    const { subscribe, set, update } = writable<AuthState>({
        user: null,
        loading: false,
        error: null
    });

    // Mock login – accepts any email/password that passes basic validation
    async function login(email: string, password: string) {
        update((state) => ({ ...state, loading: true, error: null }));
        try {
            // Simulate API delay
            await new Promise((resolve) => setTimeout(resolve, 500));

            // Simple validation
            if (!email.includes('@') || password.length < 6) {
                throw new Error('Invalid email or password (mock: email must contain @, password min 6 chars)');
            }

            const user: User = { email, name: email.split('@')[0] };
            update((state) => ({ ...state, user, loading: false }));
            goto('/'); // Redirect to home after login
        } catch (error) {
            update((state) => ({ ...state, error: (error as Error).message, loading: false }));
        }
    }

    async function register(email: string, password: string, name?: string) {
        update((state) => ({ ...state, loading: true, error: null }));
        try {
            await new Promise((resolve) => setTimeout(resolve, 500));
            if (!email.includes('@') || password.length < 6) {
                throw new Error('Invalid email or password (mock: email must contain @, password min 6 chars)');
            }
            const user: User = { email, name: name || email.split('@')[0] };
            update((state) => ({ ...state, user, loading: false }));
            goto('/');
        } catch (error) {
            update((state) => ({ ...state, error: (error as Error).message, loading: false }));
        }
    }

    function logout() {
        set({ user: null, loading: false, error: null });
        goto('/');
    }

    return {
        subscribe,
        login,
        register,
        logout
    };
}

export const auth = createAuthStore();