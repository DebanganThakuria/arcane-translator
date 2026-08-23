import { browser } from '$app/environment';

export type ChromeTheme = 'light' | 'dark';

const KEY = 'arcane:theme';

function initial(): ChromeTheme {
	if (!browser) return 'light';
	const saved = localStorage.getItem(KEY);
	if (saved === 'light' || saved === 'dark') return saved;
	return matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

let current = $state<ChromeTheme>(initial());

export const theme = {
	get value() {
		return current;
	},

	set(next: ChromeTheme) {
		current = next;
		if (!browser) return;
		localStorage.setItem(KEY, next);
		document.documentElement.dataset.theme = next;
	},

	toggle() {
		this.set(current === 'dark' ? 'light' : 'dark');
	}
};
