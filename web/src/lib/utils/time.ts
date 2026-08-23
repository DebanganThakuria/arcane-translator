/**
 * Relative time, replacing the date-fns dependency. The backend serialises
 * timestamps as Unix *seconds*; `fromUnix` is the only entry point that should
 * be given raw API values.
 */

const MINUTE = 60_000;
const HOUR = 60 * MINUTE;
const DAY = 24 * HOUR;
const MONTH = 30 * DAY;
const YEAR = 365 * DAY;

function plural(count: number, unit: string): string {
	return `${count} ${unit}${count === 1 ? '' : 's'}`;
}

/** "3 hours ago", "just now", "2 months ago". */
export function relativeTime(msSinceEpoch: number): string {
	const delta = Date.now() - msSinceEpoch;
	if (!Number.isFinite(delta)) return 'unknown';
	if (delta < MINUTE) return 'just now';
	if (delta < HOUR) return `${plural(Math.floor(delta / MINUTE), 'minute')} ago`;
	if (delta < DAY) return `${plural(Math.floor(delta / HOUR), 'hour')} ago`;
	if (delta < MONTH) return `${plural(Math.floor(delta / DAY), 'day')} ago`;
	if (delta < YEAR) return `${plural(Math.floor(delta / MONTH), 'month')} ago`;
	return `${plural(Math.floor(delta / YEAR), 'year')} ago`;
}

/** Same, for the Unix-seconds timestamps the Go API returns. */
export function fromUnix(seconds: number | undefined | null): string {
	if (!seconds) return 'unknown';
	return relativeTime(Number(seconds) * 1000);
}

/** Reading time at a relaxed 200 words per minute. */
export function readingMinutes(wordCount: number | undefined): number {
	return Math.max(1, Math.ceil((wordCount || 900) / 200));
}

export function formatCount(value: number): string {
	return value.toLocaleString('en-US');
}
