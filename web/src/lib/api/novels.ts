import { del, request } from './client';
import type {
	Chapter,
	ChapterPage,
	Novel,
	NovelFilter,
	NovelQuery,
	NovelStats,
	PaginatedNovels,
	SourceSite
} from '$lib/types';

export const EMPTY_PAGE: PaginatedNovels = {
	novels: [],
	total_count: 0,
	current_page: 1,
	total_pages: 1
};

export function getStats(): Promise<NovelStats> {
	return request<NovelStats>('/stats/novels');
}

/**
 * Lists novels. Filters, sorting and paging are all applied by the server, so
 * `total_count` reflects the filters rather than the whole library.
 */
export function getNovels(page = 1, limit = 20, query: NovelQuery = {}): Promise<PaginatedNovels> {
	const params = new URLSearchParams({ page: String(page), limit: String(limit) });

	for (const [key, value] of Object.entries(query)) {
		if (value) params.set(key, value);
	}

	return request<PaginatedNovels>(`/novels?${params}`);
}

export function getNovelsByFilter(
	filter: NovelFilter,
	value: string,
	page = 1,
	limit = 20
): Promise<PaginatedNovels> {
	const query = new URLSearchParams({
		filter,
		value,
		page: String(page),
		limit: String(limit)
	});
	return request<PaginatedNovels>(`/novels?${query}`);
}

export function getRecentlyRead(count = 6): Promise<PaginatedNovels> {
	return getNovelsByFilter('recently_read', String(count));
}

export function getNovel(id: string): Promise<Novel> {
	return request<Novel>(`/novels/${encodeURIComponent(id)}`);
}

/**
 * One page of a novel's chapters. Novels here run to several hundred chapters,
 * so the list is paged and filtered server-side rather than fetched whole.
 */
export function getChapterPage(
	novelId: string,
	page = 1,
	limit = 50,
	options: { q?: string; dir?: 'asc' | 'desc' } = {}
): Promise<ChapterPage> {
	const params = new URLSearchParams({ page: String(page), limit: String(limit) });
	if (options.q) params.set('q', options.q);
	if (options.dir) params.set('dir', options.dir);

	return request<ChapterPage>(`/novels/${encodeURIComponent(novelId)}/chapters?${params}`);
}

/** The whole chapter list, for callers that genuinely need every row. */
export function getChapters(novelId: string): Promise<Chapter[]> {
	return request<Chapter[]>(`/novels/${encodeURIComponent(novelId)}/chapters`);
}

export function getChapterByNumber(novelId: string, number: number): Promise<Chapter> {
	return request<Chapter>(`/novels/${encodeURIComponent(novelId)}/chapters/num/${number}`);
}

export function searchNovels(query: string): Promise<Novel[]> {
	return request<Novel[]>(`/search/novels/${encodeURIComponent(query)}`);
}

export function deleteNovel(id: string): Promise<void> {
	return del(`/novels/${encodeURIComponent(id)}`);
}

export function deleteChapter(novelId: string, chapterId: string): Promise<void> {
	return del(`/novels/${encodeURIComponent(novelId)}/chapters/${encodeURIComponent(chapterId)}`);
}

export function getSources(): Promise<SourceSite[]> {
	return request<SourceSite[]>('/sources');
}
