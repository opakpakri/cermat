import { writable } from 'svelte/store';

const STORAGE_KEY = 'cermat_auth';

function loadInitialState() {
	if (typeof window === 'undefined') {
		return { token: null, user: null };
	}
	try {
		const raw = localStorage.getItem(STORAGE_KEY);
		if (raw) {
			return JSON.parse(raw);
		}
	} catch {
		// ignore corrupted storage
	}
	return { token: null, user: null };
}

function createAuthStore() {
	const { subscribe, set, update } = writable(loadInitialState());

	return {
		subscribe,
		login: (token, user) => {
			const next = { token, user };
			if (typeof window !== 'undefined') {
				localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
			}
			set(next);
		},
		setUser: (user) => {
			update((state) => {
				const next = { ...state, user };
				if (typeof window !== 'undefined') {
					localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
				}
				return next;
			});
		},
		logout: () => {
			if (typeof window !== 'undefined') {
				localStorage.removeItem(STORAGE_KEY);
			}
			set({ token: null, user: null });
		}
	};
}

export const authStore = createAuthStore();
