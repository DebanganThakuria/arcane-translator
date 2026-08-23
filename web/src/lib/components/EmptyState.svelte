<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		title: string;
		body?: string;
		action?: Snippet;
	}

	let { title, body, action }: Props = $props();
</script>

<div class="empty">
	<h3 class="display">{title}</h3>
	{#if body}<p>{body}</p>{/if}
	{#if action}<div class="action">{@render action()}</div>{/if}
</div>

<style>
	.empty {
		display: grid;
		justify-items: center;
		gap: var(--sp-3);
		padding: var(--sp-8) var(--sp-5);
		text-align: center;
		border: 1px dashed var(--line-strong);
		border-radius: var(--r-2);
		background: var(--surface);
	}

	h3 {
		font-size: var(--step-2);
		color: var(--ink);
	}

	p {
		max-width: 42ch;
		color: var(--ink-2);
		font-size: var(--step-0);
	}

	.action {
		margin-top: var(--sp-2);
	}
</style>
