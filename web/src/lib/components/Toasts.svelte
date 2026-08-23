<script lang="ts">
	import { toast } from '$stores/toast.svelte';
	import { CheckCircle, Info, Warning, X } from 'phosphor-svelte';

	const ICONS = { info: Info, success: CheckCircle, error: Warning };
</script>

<div class="dock" role="status" aria-live="polite">
	{#each toast.items as item (item.id)}
		{@const Icon = ICONS[item.tone]}
		<div class="toast" data-tone={item.tone}>
			<Icon size={16} weight="bold" />
			<div class="text">
				<p class="title">{item.title}</p>
				{#if item.body}<p class="body">{item.body}</p>{/if}
			</div>
			<button class="close" onclick={() => toast.dismiss(item.id)} aria-label="Dismiss">
				<X size={13} weight="bold" />
			</button>
		</div>
	{/each}
</div>

<style>
	.dock {
		position: fixed;
		right: var(--sp-4);
		bottom: var(--sp-4);
		z-index: var(--z-toast);
		display: grid;
		gap: var(--sp-2);
		width: min(23rem, calc(100vw - 2rem));
		pointer-events: none;
	}

	.toast {
		pointer-events: auto;
		display: grid;
		grid-template-columns: auto 1fr auto;
		gap: var(--sp-3);
		align-items: start;
		padding: var(--sp-3) var(--sp-4);
		background: var(--surface-2);
		border: 1px solid var(--line-strong);
		border-left: 2px solid var(--ink-3);
		border-radius: var(--r-1);
		box-shadow: var(--shadow-overlay);
	}

	.toast[data-tone='success'] {
		border-left-color: var(--positive);
		color: var(--positive);
	}

	.toast[data-tone='error'] {
		border-left-color: var(--danger);
		color: var(--danger);
	}

	.toast[data-tone='info'] {
		color: var(--ink-2);
	}

	.text {
		min-width: 0;
	}

	.title {
		font-size: var(--step-0);
		font-weight: 600;
		color: var(--ink);
	}

	.body {
		margin-top: 0.15rem;
		font-size: var(--step--1);
		color: var(--ink-2);
		overflow-wrap: anywhere;
	}

	.close {
		background: none;
		border: 0;
		padding: 0.15rem;
		color: var(--ink-3);
		cursor: pointer;
		border-radius: var(--r-1);
	}

	.close:hover {
		color: var(--ink);
	}

	@media (prefers-reduced-motion: no-preference) {
		.toast {
			animation: slide-in var(--dur-2) var(--ease);
		}
	}

	@keyframes slide-in {
		from {
			opacity: 0;
			transform: translateY(8px);
		}
	}
</style>
