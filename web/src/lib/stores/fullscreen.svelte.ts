import { browser } from '$app/environment';

/**
 * Fullscreen for the reader.
 *
 * Supported in every desktop browser and Android Chrome. iPhone Safari does not
 * implement the Fullscreen API for anything other than <video>, so `supported`
 * is false there and callers hide the control rather than offering a button
 * that does nothing. iPadOS Safari does support it.
 *
 * Safari only exposed the unprefixed API in 16.4, hence the webkit fallbacks.
 */

interface WebkitDocument extends Document {
	webkitFullscreenEnabled?: boolean;
	webkitFullscreenElement?: Element | null;
	webkitExitFullscreen?: () => Promise<void>;
}

interface WebkitElement extends HTMLElement {
	webkitRequestFullscreen?: () => Promise<void>;
}

function doc(): WebkitDocument | null {
	return browser ? (document as WebkitDocument) : null;
}

let active = $state(false);

export const fullscreen = {
	/** Whether this browser can put an element fullscreen at all. */
	get supported(): boolean {
		const d = doc();
		return Boolean(d && (d.fullscreenEnabled || d.webkitFullscreenEnabled));
	},

	get active() {
		return active;
	},

	/** Must be called from a user gesture; browsers reject it otherwise. */
	async toggle(target?: HTMLElement) {
		const d = doc();
		if (!d) return;

		const current = d.fullscreenElement ?? d.webkitFullscreenElement;

		try {
			if (current) {
				await (d.exitFullscreen?.() ?? d.webkitExitFullscreen?.());
				return;
			}

			const element = (target ?? d.documentElement) as WebkitElement;
			await (element.requestFullscreen?.() ?? element.webkitRequestFullscreen?.());
		} catch {
			// A rejected request is not worth interrupting reading over; the
			// state listener keeps the button honest either way.
		}
	},

	/** Keeps `active` in sync, including when the reader presses Escape. */
	watch(): () => void {
		const d = doc();
		if (!d) return () => {};

		const sync = () => {
			active = Boolean(d.fullscreenElement ?? d.webkitFullscreenElement);
		};

		sync();
		d.addEventListener('fullscreenchange', sync);
		d.addEventListener('webkitfullscreenchange', sync);

		return () => {
			d.removeEventListener('fullscreenchange', sync);
			d.removeEventListener('webkitfullscreenchange', sync);
		};
	}
};
