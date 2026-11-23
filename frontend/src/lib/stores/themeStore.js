import { writable } from 'svelte/store';
import { browser } from '$app/env';

const THEME_KEY = 'cermat_theme';

function resolveIsDark(theme) {
	if (!browser) return false;
	if (theme === 'dark') return true;
	if (theme === 'light') return false;
	return window.matchMedia('(prefers-color-scheme: dark)').matches;
}

function applyTheme(theme) {
	if (!browser) return;
	document.documentElement.classList.toggle('dark', resolveIsDark(theme));
}

function createThemeStore() {
	const initial = browser ? localStorage.getItem(THEME_KEY) || 'light' : 'light';
	let currentTheme = initial;
	applyTheme(initial);

	if (browser) {
		const media = window.matchMedia('(prefers-color-scheme: dark)');
		media.addEventListener('change', () => {
			if (currentTheme === 'system') {
				applyTheme('system');
			}
		});
	}

	const { subscribe, set } = writable(initial);

	return {
		subscribe,
		setTheme(theme) {
			currentTheme = theme;
			if (browser) {
				localStorage.setItem(THEME_KEY, theme);
				applyTheme(theme);
			}
			set(theme);
		}
	};
}

export const themeStore = createThemeStore();
