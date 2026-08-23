export type ToastTone = 'info' | 'success' | 'error';

export interface Toast {
	id: number;
	tone: ToastTone;
	title: string;
	body?: string;
}

let items = $state<Toast[]>([]);
let nextId = 0;

const LIFETIME: Record<ToastTone, number> = {
	info: 3500,
	success: 3500,
	error: 7000
};

function push(tone: ToastTone, title: string, body?: string): number {
	const id = nextId++;
	items.push({ id, tone, title, body });
	setTimeout(() => toast.dismiss(id), LIFETIME[tone]);
	return id;
}

export const toast = {
	get items() {
		return items;
	},

	info: (title: string, body?: string) => push('info', title, body),
	success: (title: string, body?: string) => push('success', title, body),
	error: (title: string, body?: string) => push('error', title, body),

	/** Reports a caught value without leaking `[object Object]` into the UI. */
	fromError(title: string, cause: unknown) {
		return push('error', title, cause instanceof Error ? cause.message : String(cause));
	},

	dismiss(id: number) {
		items = items.filter((item) => item.id !== id);
	}
};
