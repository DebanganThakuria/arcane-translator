<script lang="ts">
	import NovelCard from './NovelCard.svelte';
	import type { Novel } from '$lib/types';

	interface Props {
		novels: Novel[];
		/** Number of leading covers to load eagerly. */
		eagerCount?: number;
	}

	let { novels, eagerCount = 6 }: Props = $props();
</script>

<ul class="grid">
	{#each novels as novel, i (novel.id)}
		<li><NovelCard {novel} eager={i < eagerCount} /></li>
	{/each}
</ul>

<style>
	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(9.5rem, 1fr));
		gap: var(--sp-5) var(--sp-4);
	}

	@media (max-width: 480px) {
		.grid {
			grid-template-columns: repeat(2, 1fr);
			gap: var(--sp-4) var(--sp-3);
		}
	}
</style>
