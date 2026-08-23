import { goto } from '$app/navigation';
import { page } from '$app/state';
import { EMPTY_PAGE } from '$api/novels';
import { toast } from '$stores/toast.svelte';
import type { PaginatedNovels } from '$lib/types';

export const PER_PAGE_OPTIONS = [12, 24, 48, 100] as const;
const DEFAULT_PER_PAGE = 24;

type Fetcher = (page: number, limit: number) => Promise<PaginatedNovels>;

function readInt(key: string, fallback: number): number {
	const raw = Number(page.url.searchParams.get(key));
	return Number.isInteger(raw) && raw > 0 ? raw : fallback;
}

/**
 * Shared paging for the library, language and genre shelves.
 *
 * The URL is the source of truth, so a page is shareable, the back button
 * steps through history, and a reload lands where the reader left off.
 */
export function createPagination(fetcher: Fetcher) {
	let data = $state<PaginatedNovels>(EMPTY_PAGE);
	let loading = $state(true);
	let refreshing = $state(false);
	let failed = $state(false);
	let requestId = 0;

	const current = $derived(readInt('page', 1));
	const limit = $derived(readInt('per', DEFAULT_PER_PAGE));

	async function load() {
		const ticket = ++requestId;
		loading = true;
		failed = false;
		try {
			const result = await fetcher(current, limit);
			// A slower earlier request must not overwrite a newer result.
			if (ticket !== requestId) return;
			data = result;
		} catch (error) {
			if (ticket !== requestId) return;
			failed = true;
			data = EMPTY_PAGE;
			toast.fromError('Could not load novels', error);
		} finally {
			if (ticket === requestId) loading = false;
		}
	}

	function navigate(nextPage: number, nextLimit: number) {
		const params = new URLSearchParams(page.url.searchParams);

		if (nextPage > 1) params.set('page', String(nextPage));
		else params.delete('page');

		if (nextLimit !== DEFAULT_PER_PAGE) params.set('per', String(nextLimit));
		else params.delete('per');

		const query = params.toString();
		goto(`${page.url.pathname}${query ? `?${query}` : ''}`, { noScroll: true, keepFocus: true });
	}

	return {
		get data() {
			return data;
		},
		get loading() {
			return loading;
		},
		get refreshing() {
			return refreshing;
		},
		get failed() {
			return failed;
		},
		get limit() {
			return limit;
		},
		get current() {
			return current;
		},

		load,

		goTo(nextPage: number) {
			if (nextPage >= 1 && nextPage <= data.total_pages && nextPage !== current) {
				navigate(nextPage, limit);
			}
		},

		async refresh() {
			refreshing = true;
			await load();
			refreshing = false;
		},

		setLimit(nextLimit: number) {
			if (nextLimit !== limit) navigate(1, nextLimit);
		}
	};
}

export type Pagination = ReturnType<typeof createPagination>;

/**
 * Page numbers with gaps, e.g. [1, -1, 7, 8, 9, -1, 24]. -1 renders an ellipsis.
 */
export function pageWindow(current: number, total: number, span = 1): number[] {
	if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1);

	const pages = new Set<number>([1, total]);
	for (let p = current - span; p <= current + span; p++) {
		if (p > 1 && p < total) pages.add(p);
	}

	const sorted = [...pages].sort((a, b) => a - b);
	const out: number[] = [];
	for (let i = 0; i < sorted.length; i++) {
		if (i > 0 && sorted[i] - sorted[i - 1] > 1) out.push(-1);
		out.push(sorted[i]);
	}
	return out;
}
