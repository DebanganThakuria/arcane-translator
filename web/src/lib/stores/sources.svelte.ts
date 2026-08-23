import { getSources } from '$api/novels';
import type { Language, SourceSite } from '$lib/types';

/**
 * Short Latin codes rather than flag emoji: they render identically on every
 * platform, stay legible at 11px, and keep the interface in one script.
 */
const MARKS: Record<string, string> = {
	chinese: 'CN',
	korean: 'KR',
	japanese: 'JP'
};

let sites = $state<SourceSite[]>([]);
let byId = $state<Record<string, Language>>({});
let loaded = $state(false);
let inflight: Promise<void> | null = null;

export const sources = {
	get all() {
		return sites;
	},

	get loaded() {
		return loaded;
	},

	/** Idempotent: concurrent callers share one request. */
	load(): Promise<void> {
		if (loaded) return Promise.resolve();
		inflight ??= getSources()
			.then((fetched) => {
				sites = fetched;
				byId = Object.fromEntries(fetched.map((site) => [site.id, site.language]));
			})
			.catch(() => {
				// The sources list is decoration on most screens. Fail quiet and let
				// the language fallback below handle labelling.
				sites = [];
				byId = {};
			})
			.finally(() => {
				loaded = true;
				inflight = null;
			});
		return inflight;
	},

	byId(id: string): SourceSite | undefined {
		return sites.find((site) => site.id === id);
	},

	/** Source id first, then a substring guess for ids we have not fetched. */
	language(sourceId: string): Language {
		const known = byId[sourceId];
		if (known) return known;
		const haystack = sourceId.toLowerCase();
		if (/chinese|china|\bcn\b/.test(haystack)) return 'Chinese';
		if (/korean|korea|\bkr\b/.test(haystack)) return 'Korean';
		if (/japanese|japan|\bjp\b/.test(haystack)) return 'Japanese';
		return 'Other';
	},

	mark(sourceId: string): string {
		return MARKS[this.language(sourceId).toLowerCase()] ?? 'WEB';
	}
};
