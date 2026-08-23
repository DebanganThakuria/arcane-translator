export type NovelStatus = 'Ongoing' | 'Completed' | 'Unknown';
export type Language = 'Chinese' | 'Korean' | 'Japanese' | 'Other';

/** Unix seconds, as the Go backend serialises them. */
export type UnixSeconds = number;

export interface Novel {
	id: string;
	title: string;
	original_title?: string;
	cover?: string;
	source: string;
	url: string;
	summary: string;
	author?: string;
	status?: NovelStatus;
	genres?: string[];
	chapters_count: number;
	first_chapter_url?: string;
	last_read_chapter_number: number;
	last_read_timestamp: UnixSeconds;
	last_updated: UnixSeconds;
	date_added: UnixSeconds;
}

export interface PaginatedNovels {
	novels: Novel[];
	total_count: number;
	current_page: number;
	total_pages: number;
}

export interface Chapter {
	id: string;
	novel_id: string;
	number: number;
	title: string;
	original_tile?: string;
	content: string;
	date_translated: UnixSeconds;
	word_count?: number;
	url?: string;
	next_chapter_url?: string;
}

export interface SourceSite {
	id: string;
	name: string;
	url: string;
	language: Language;
	icon?: string;
}

export interface NovelStats {
	novel_count: number;
	chapter_count: number;
}

export type NovelFilter = 'language' | 'genre' | 'recently_read' | 'recently_updated';

/** Sort keys the backend accepts for the novel list. */
export type NovelSort = 'last_read' | 'last_updated' | 'date_added' | 'title' | 'chapters';

/** Server-side novel filters. Everything is optional and composes with AND. */
export interface NovelQuery {
	q?: string;
	language?: string;
	genre?: string;
	status?: string;
	sort?: NovelSort;
	dir?: 'asc' | 'desc';
}

/**
 * One page of a novel's chapters. `first_number` and `last_number` describe the
 * whole novel, not the page, so "start from the beginning" needs no extra call.
 */
export interface ChapterPage {
	chapters: Chapter[];
	total_count: number;
	current_page: number;
	total_pages: number;
	first_number: number;
	last_number: number;
}
