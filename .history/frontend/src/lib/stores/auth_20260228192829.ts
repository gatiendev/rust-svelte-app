import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// Types
export interface User {
    id: string;
    username: string;
    created_at: string;
}

interface AuthState {
    user: User | null;
    loading: boolean;
    error: string | null;
}

// API base URL from environment
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000';

// Helper to handle fetch responses
async function handleResponse<T>(response: Response): Promise<T> {
    if (!response.ok) {
        // Try to parse error message from backend (which sends { "message": "..." })
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || `Request failed with status ${response.status}`);
    }
    // For 204 No Content, return empty object
    if (response.status === 204) {
        return {} as T;
    }
    return response.json();
}

function createAuthStore() {
    const { subscribe, set, update } = writable<AuthState>({
        user: null,
        loading: false,
        error: null,
    });

    // Load user from session on initial load (if access token exists)
    // We can call /profile to check authentication
    if (browser) {
        // Optionally fetch user on mount
        fetch(`${API_URL}/profile`, {
            credentials: 'include', // important: include cookies
        })
            .then(res => {
                if (res.ok) return res.json();
                throw new Error('Not authenticated');
            })
            .then(userData => {
                update(state => ({ ...state, user: userData }));
            })
            .catch(() => {
                // No active session – do nothing
            });
    }

    return {
        subscribe,

        async login(email: string, password: string) {
            update(state => ({ ...state, loading: true, error: null }));
            try {
                const response = await fetch(`${API_URL}/login`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username: email, password }),
                    credentials: 'include', // required to receive cookies
                });

                if (!response.ok) {
                    const errorData = await response.json().catch(() => ({}));
                    throw new Error(errorData.message || 'Login failed');
                }

                // After successful login, fetch user profile
                const profileRes = await fetch(`${API_URL}/profile`, {
                    credentials: 'include',
                });
                if (!profileRes.ok) throw new Error('Failed to fetch user profile');
                const user = await profileRes.json();

                update(state => ({ ...state, user, loading: false }));
            } catch (err: any) {
                update(state => ({ ...state, error: err.message, loading: false }));
            }
        },

        async register(email: string, password: string, name?: string) {
            update(state => ({ ...state, loading: true, error: null }));
            try {
                const response = await fetch(`${API_URL}/register`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username: email, password }),
                    credentials: 'include',
                });

                if (!response.ok) {
                    const errorData = await response.json().catch(() => ({}));
                    throw new Error(errorData.message || 'Registration failed');
                }

                // Registration does not automatically log in – you may want to call login afterwards
                // Optionally, you can automatically log in:
                await this.login(email, password);
            } catch (err: any) {
                update(state => ({ ...state, error: err.message, loading: false }));
            }
        },

        async logout() {
            update(state => ({ ...state, loading: true }));
            try {
                await fetch(`${API_URL}/logout`, {
                    method: 'POST',
                    credentials: 'include',
                });
                // Clear user regardless of response (server should have cleared cookies)
                update(state => ({ ...state, user: null, loading: false, error: null }));
            } catch (err: any) {
                update(state => ({ ...state, error: err.message, loading: false }));
            }
        },
    };
}

export const auth = createAuthStore();