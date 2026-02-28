import { writable } from 'svelte/store';
import { browser } from '$app/environment';

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

// Read API base URL from environment (set in .env)
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

    // On initial load, try to fetch the current user (if already logged in)
    if (browser) {
        fetch(`${API_URL}/profile`, {
            credentials: 'include',
        })
            .then(res => {
                if (res.ok) return res.json();
                throw new Error('Not authenticated');
            })
            .then(userData => {
                update(state => ({ ...state, user: userData }));
            })
            .catch(() => {
                // No active session – ignore
            });
    }

    return {
        subscribe,

        async login(email: string, password: string) {
            update(state => ({ ...state, loading: true, error: null }));
            try {
                // 1. Call login endpoint with username = email
                const loginRes = await fetch(`${API_URL}/login`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username: email, password }),
                    credentials: 'include',
                });

                if (!loginRes.ok) {
                    const errorData = await loginRes.json().catch(() => ({}));
                    throw new Error(errorData.message || 'Login failed');
                }

                // 2. After successful login, fetch user profile
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
                // Backend expects username and password (name is currently ignored)
                const registerRes = await fetch(`${API_URL}/register`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username: email, password }),
                    credentials: 'include',
                });

                if (!registerRes.ok) {
                    const errorData = await registerRes.json().catch(() => ({}));
                    throw new Error(errorData.message || 'Registration failed');
                }

                // Automatically log in after successful registration
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
                // Clear user regardless of server response
                update(state => ({ ...state, user: null, loading: false, error: null }));
            } catch (err: any) {
                update(state => ({ ...state, error: err.message, loading: false }));
            }
        },
    };
}

export const auth = createAuthStore();