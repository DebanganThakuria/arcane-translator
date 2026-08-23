import { browser } from '$app/environment';

/**
 * The Go backend listens on :8088 on the same host that serves this app, so the
 * API host is derived from the page origin. That keeps LAN access working when
 * the app is opened from a phone on the same network.
 */
export const API_BASE_URL = browser ? `http://${window.location.hostname}:8088` : 'http://localhost:8088';

export class ApiError extends Error {
	readonly status: number;

	constructor(message: string, status: number) {
		super(message);
		this.name = 'ApiError';
		this.status = status;
	}
}

/** Maps the transport failure onto something worth showing a reader. */
function describe(status: number, detail: string | undefined): string {
	if (detail) return detail;
	switch (status) {
		case 404:
			return 'Not found. Check that the URL is correct.';
		case 403:
			return 'The source site refused the request. It may be blocking scrapers.';
		case 429:
			return 'The source site is rate limiting us. Wait a moment and retry.';
		case 500:
			return 'The server failed while handling that request.';
		case 503:
			return 'The server is unavailable. Is the backend running on port 8088?';
		default:
			return `Request failed with status ${status}.`;
	}
}

async function extractDetail(response: Response): Promise<string | undefined> {
	const raw = await response.text().catch(() => '');
	if (!raw) return undefined;
	try {
		const parsed = JSON.parse(raw) as { detail?: string; error?: string; message?: string };
		return parsed.detail ?? parsed.error ?? parsed.message;
	} catch {
		// Go's http.Error writes plain text.
		return raw.trim().slice(0, 300) || undefined;
	}
}

export async function request<T>(path: string, init?: RequestInit): Promise<T> {
	let response: Response;
	try {
		response = await fetch(`${API_BASE_URL}${path}`, {
			...init,
			headers:
				init?.body !== undefined
					? { 'Content-Type': 'application/json', ...init?.headers }
					: init?.headers
		});
	} catch {
		throw new ApiError(
			'Cannot reach the translation server. Check that the backend is running on port 8088.',
			0
		);
	}

	if (!response.ok) {
		throw new ApiError(describe(response.status, await extractDetail(response)), response.status);
	}

	if (response.status === 204) return undefined as T;
	return (await response.json()) as T;
}

export function post<T>(path: string, body: unknown): Promise<T> {
	return request<T>(path, { method: 'POST', body: JSON.stringify(body) });
}

export function del(path: string): Promise<void> {
	return request<void>(path, { method: 'DELETE' });
}
