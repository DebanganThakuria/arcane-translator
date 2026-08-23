<script lang="ts">
	import Dialog from './Dialog.svelte';

	interface Props {
		open: boolean;
		title: string;
		description: string;
		confirmLabel?: string;
		busy?: boolean;
		onconfirm: () => void;
		oncancel: () => void;
	}

	let {
		open,
		title,
		description,
		confirmLabel = 'Delete',
		busy = false,
		onconfirm,
		oncancel
	}: Props = $props();
</script>

<Dialog {open} {title} onclose={oncancel}>
	<p>{description}</p>

	{#snippet footer()}
		<button class="btn" onclick={oncancel} disabled={busy}>Cancel</button>
		<button class="btn btn--confirm" onclick={onconfirm} disabled={busy}>
			{busy ? 'Working…' : confirmLabel}
		</button>
	{/snippet}
</Dialog>

<style>
	p {
		color: var(--ink-2);
		text-wrap: pretty;
	}

	/* Solid destructive fill so the label always clears AA against it. */
	.btn--confirm {
		background: var(--danger);
		border-color: var(--danger);
		color: var(--accent-ink);
	}

	.btn--confirm:hover:not(:disabled) {
		filter: brightness(1.1);
	}
</style>
