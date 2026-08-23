import { browser } from '$app/environment';

export type ReaderTheme = 'sepia' | 'light' | 'dark';
export type ReaderFont = 'verdana' | 'georgia' | 'system' | 'mono';
/** `scroll` is one continuous column; `paged` is a two-page spread. */
export type ReaderMode = 'scroll' | 'paged';

export interface ReaderPrefs {
	theme: ReaderTheme;
	font: ReaderFont;
	mode: ReaderMode;
	fontSize: number;
	lineHeight: number;
	measure: number;
	justify: boolean;
}

export const FONT_LABELS: Record<ReaderFont, string> = {
	verdana: 'Verdana',
	georgia: 'Georgia',
	system: 'System sans',
	mono: 'Monospace'
};

export const FONT_STACKS: Record<ReaderFont, string> = {
	verdana: 'var(--font-read-verdana)',
	georgia: 'var(--font-read-georgia)',
	system: 'var(--font-read-system)',
	mono: 'var(--font-read-mono)'
};

export const LIMITS = {
	fontSize: { min: 14, max: 30, step: 1 },
	lineHeight: { min: 1.3, max: 2.4, step: 0.1 },
	measure: { min: 46, max: 130, step: 2 }
} as const;

/** Sepia and Verdana are the defaults the reader asked to keep. */
const DEFAULTS: ReaderPrefs = {
	theme: 'sepia',
	font: 'verdana',
	mode: 'scroll',
	fontSize: 18,
	lineHeight: 1.7,
	// Wide enough that a chapter is a few screens rather than a dozen. The
	// slider goes wider still for large displays.
	measure: 86,
	justify: false
};

const KEY = 'arcane:reader';

/**
 * Bumped when a default changes in a way existing readers should pick up.
 * Anything they explicitly chose is preserved; only untouched settings move.
 */
const PREFS_VERSION = 2;

function clamp(value: number, { min, max }: { min: number; max: number }) {
	return Math.min(Math.max(value, min), max);
}

function initial(): ReaderPrefs {
	if (!browser) return { ...DEFAULTS };
	try {
		const raw = localStorage.getItem(KEY);
		if (!raw) return { ...DEFAULTS };
		const saved = JSON.parse(raw) as Partial<ReaderPrefs> & { version?: number };

		// Readers stored before the measure default widened kept the old 68ch
		// column, which meant a lot of scrolling on a large screen.
		const measure =
			saved.version === PREFS_VERSION ? (saved.measure ?? DEFAULTS.measure) : DEFAULTS.measure;

		return {
			theme: saved.theme ?? DEFAULTS.theme,
			font: saved.font ?? DEFAULTS.font,
			mode: saved.mode ?? DEFAULTS.mode,
			fontSize: clamp(saved.fontSize ?? DEFAULTS.fontSize, LIMITS.fontSize),
			lineHeight: clamp(saved.lineHeight ?? DEFAULTS.lineHeight, LIMITS.lineHeight),
			measure: clamp(measure, LIMITS.measure),
			justify: saved.justify ?? DEFAULTS.justify
		};
	} catch {
		return { ...DEFAULTS };
	}
}

let prefs = $state<ReaderPrefs>(initial());

// Stamp the current version straight away, so a reader who then narrows the
// column keeps their choice rather than having it reset next time.
if (browser) persist();

function persist() {
	if (browser) localStorage.setItem(KEY, JSON.stringify({ ...prefs, version: PREFS_VERSION }));
}

export const reader = {
	get prefs() {
		return prefs;
	},

	set<K extends keyof ReaderPrefs>(key: K, value: ReaderPrefs[K]) {
		prefs[key] = value;
		persist();
	},

	reset() {
		prefs = { ...DEFAULTS };
		persist();
	}
};
