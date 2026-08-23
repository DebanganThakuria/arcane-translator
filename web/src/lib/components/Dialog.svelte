<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		open: boolean;
		title: string;
		description?: string;
		/** Caps the panel width; dialogs with a paste area want more room. */
		wide?: boolean;
		onclose: () => void;
		children: Snippet;
		footer?: Snippet;
	}

	let { open, title, description, wide = false, onclose, children, footer }: Props = $props();

	let el = $state<HTMLDialogElement>();

	// The native dialog gives us focus trapping, inert background and Esc for
	// free, so the only job here is keeping it in sync with `open`.
	$effect(() => {
		if (!el) return;
		if (open && !el.open) el.showModal();
		if (!open && el.open) el.close();
	});
</script>

<dialog
	bind:this={el}
	class:wide
	oncancel={(event) => {
		event.preventDefault();
		onclose();
	}}
	onclick={(event) => {
		// Clicks land on the dialog itself only when they hit the backdrop.
		if (event.target === el) onclose();
	}}
>
	<div class="panel">
		<header>
			<h2 class="display">{title}</h2>
			{#if description}<p class="meta">{description}</p>{/if}
		</header>

		<div class="body">{@render children()}</div>

		{#if footer}
			<footer>{@render footer()}</footer>
		{/if}
	</div>
</dialog>

<style>
	dialog {
		/* Explicit, because the global reset zeroes the UA's auto margins. */
		margin: auto;
		padding: 0;
		border: 0;
		background: transparent;
		max-width: min(30rem, calc(100vw - 2rem));
		width: 100%;
		color: var(--ink);
	}

	dialog.wide {
		max-width: min(44rem, calc(100vw - 2rem));
	}

	dialog::backdrop {
		background: #0b0d0e9e;
		backdrop-filter: blur(2px);
	}

	.panel {
		display: grid;
		gap: var(--sp-4);
		max-height: calc(100dvh - 4rem);
		overflow-y: auto;
		overscroll-behavior: contain;
		padding: var(--sp-5);
		background: var(--surface);
		border: 1px solid var(--line-strong);
		border-radius: var(--r-2);
		box-shadow: var(--shadow-overlay);
	}

	header {
		display: grid;
		gap: var(--sp-1);
	}

	h2 {
		font-size: var(--step-2);
	}

	footer {
		display: flex;
		justify-content: flex-end;
		gap: var(--sp-2);
		padding-top: var(--sp-1);
	}

	@media (prefers-reduced-motion: no-preference) {
		dialog[open] {
			animation: rise var(--dur-2) var(--ease);
		}
	}

	@keyframes rise {
		from {
			opacity: 0;
			transform: translateY(6px);
		}
	}
</style>
