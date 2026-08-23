import { browser } from '$app/environment';

export interface ChapterProgress {
	lastChapter: number;
	lastReadAt: number;
	/** Percentage scrolled through the chapter, 0 to 100. */
	progress: number;
	chapterTitle?: string;
}

type ProgressMap = Record<string, ChapterProgress>;

const KEY = 'novel_reading_progress';

function load(): ProgressMap {
	if (!browser) return {};
	try {
		return JSON.parse(localStorage.getItem(KEY) ?? '{}') as ProgressMap;
	} catch {
		return {};
	}
}

let entries = $state<ProgressMap>(load());

function persist() {
	if (!browser) return;
	try {
		localStorage.setItem(KEY, JSON.stringify(entries));
	} catch {
		// Storage full or blocked. Progress stays in memory for this session.
	}
}

export const progress = {
	get(novelId: string): ChapterProgress | null {
		return entries[novelId] ?? null;
	},

	has(novelId: string): boolean {
		const entry = entries[novelId];
		return Boolean(entry && entry.lastChapter > 0);
	},

	save(novelId: string, chapterNumber: number, percent = 0, chapterTitle?: string) {
		entries[novelId] = {
			lastChapter: chapterNumber,
			lastReadAt: Date.now(),
			progress: Math.min(Math.max(percent, 0), 100),
			chapterTitle: chapterTitle ?? entries[novelId]?.chapterTitle
		};
		persist();
	},

	clear(novelId: string) {
		delete entries[novelId];
		persist();
	}
};
