import { post } from './client';
import { manualExtract } from '$stores/manualExtract.svelte';
import type { Chapter, Novel } from '$lib/types';

function assertHttpUrl(url: string) {
	if (!/^https?:\/\//i.test(url)) {
		throw new Error('URL must start with http:// or https://');
	}
}

/**
 * Every translation endpoint has the same shape: try server-side scraping
 * first, and if the source site blocks us, ask the reader to paste the page
 * source and retry with `html_content`.
 */
async function withManualFallback<T>(
	scrapeUrl: string,
	send: (html?: string) => Promise<T>
): Promise<T> {
	try {
		return await send();
	} catch (serverError) {
		// Opens the page so the reader can copy its source, then waits for them.
		window.open(scrapeUrl, '_blank', 'noopener');
		const html = await manualExtract.request(scrapeUrl).catch(() => {
			throw serverError;
		});
		return await send(html);
	}
}

export function extractNovel(url: string, source: string): Promise<Novel> {
	assertHttpUrl(url);
	return withManualFallback(url, (html_content) =>
		post<Novel>('/novels/translate', { url, source, html_content })
	);
}

export function translateChapter(novelId: string, chapterUrl?: string): Promise<Chapter> {
	if (!chapterUrl) {
		return post<Chapter>('/novels/translate/chapter', { novel_id: novelId });
	}
	assertHttpUrl(chapterUrl);
	return withManualFallback(chapterUrl, (html_content) =>
		post<Chapter>('/novels/translate/chapter', {
			novel_id: novelId,
			chapter_url: chapterUrl,
			html_content
		})
	);
}

export function setFirstChapter(novelId: string, chapterUrl: string): Promise<Chapter> {
	assertHttpUrl(chapterUrl);
	return withManualFallback(chapterUrl, (html_content) =>
		post<Chapter>('/novels/translate/first_chapter', {
			novel_id: novelId,
			chapter_url: chapterUrl,
			html_content
		})
	);
}

export function refreshNovel(novelId: string, novelUrl: string): Promise<Novel> {
	return withManualFallback(novelUrl, (html_content) =>
		post<Novel>('/novels/refresh', { novel_id: novelId, html_content })
	);
}
