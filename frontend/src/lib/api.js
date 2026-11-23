import { get } from 'svelte/store';
import { authStore } from '$lib/stores/authStore';

const API_BASE = '/api';

export async function apiFetch(path, options = {}) {
	const state = get(authStore);
	const headers = {
		'Content-Type': 'application/json',
		...(options.headers || {})
	};

	if (state?.token) {
		headers['Authorization'] = `Bearer ${state.token}`;
	}

	const res = await fetch(`${API_BASE}${path}`, {
		...options,
		headers
	});

	if (res.status === 401 && !path.startsWith('/auth/')) {
		authStore.logout();
		if (typeof window !== 'undefined') {
			window.location.href = '/login';
		}
		throw new Error('Sesi Anda telah berakhir. Silakan login kembali.');
	}

	const data = await res.json().catch(() => ({}));
	if (!res.ok) {
		throw new Error(data.error || `Request gagal (${res.status})`);
	}

	return data;
}

export const api = {
	get: (path) => apiFetch(path, { method: 'GET' }),
	post: (path, body) => apiFetch(path, { method: 'POST', body: JSON.stringify(body) }),
	put: (path, body) => apiFetch(path, { method: 'PUT', body: JSON.stringify(body) }),
	patch: (path, body) => apiFetch(path, { method: 'PATCH', body: JSON.stringify(body) }),
	delete: (path) => apiFetch(path, { method: 'DELETE' })
};
