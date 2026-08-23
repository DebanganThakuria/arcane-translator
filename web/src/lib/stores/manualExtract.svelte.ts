/**
 * Some source sites block server-side scraping. When that happens the backend
 * accepts an `html_content` field instead, and we ask the reader to fetch the
 * page themselves and paste the source in.
 *
 * This store holds the single in-flight request; `ManualExtractDialog` renders
 * it. Callers just await a promise and never touch the DOM.
 */

interface PendingRequest {
	url: string;
	resolve: (html: string) => void;
	reject: (reason: Error) => void;
}

export class ManualExtractCancelled extends Error {
	constructor() {
		super('Manual extraction was cancelled.');
		this.name = 'ManualExtractCancelled';
	}
}

let pending = $state<PendingRequest | null>(null);

export const manualExtract = {
	get pending() {
		return pending;
	},

	/** Opens the paste dialog and resolves with the HTML the reader supplied. */
	request(url: string): Promise<string> {
		pending?.reject(new ManualExtractCancelled());
		return new Promise<string>((resolve, reject) => {
			pending = { url, resolve, reject };
		});
	},

	submit(html: string) {
		pending?.resolve(html);
		pending = null;
	},

	cancel() {
		pending?.reject(new ManualExtractCancelled());
		pending = null;
	}
};

export const MIN_HTML_LENGTH = 50;
