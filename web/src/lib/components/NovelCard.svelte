<script lang="ts">
	import Cover from './Cover.svelte';
	import { sources } from '$stores/sources.svelte';
	import { progress } from '$stores/progress.svelte';
	import { fromUnix } from '$utils/time';
	import type { Novel } from '$lib/types';

	interface Props {
		novel: Novel;
		eager?: boolean;
	}

	let { novel, eager = false }: Props = $props();

	const saved = $derived(progress.get(novel.id));
	const chapterRead = $derived(saved?.lastChapter ?? novel.last_read_chapter_number ?? 0);

	const percent = $derived(
		novel.chapters_count > 0 && chapterRead > 0
			? Math.min(100, Math.round((chapterRead / novel.chapters_count) * 100))
			: 0
	);
</script>

<a class="card" href="/novel/{novel.id}">
	<div class="art">
		<Cover {novel} {eager} />
		<span class="mark" title="{sources.language(novel.source)} source">{sources.mark(novel.source)}</span>
		{#if percent > 0}
			<div class="track" aria-hidden="true">
				<div class="fill" style:width="{percent}%"></div>
			</div>
		{/if}
	</div>

	<h3 class="title">{novel.title}</h3>

	<p class="meta">
		<span class="numeral">{novel.chapters_count}</span> ch
		{#if percent > 0}
			<span class="sep" aria-hidden="true"></span>
			<span class="read">{percent}% read</span>
		{:else if novel.last_updated}
			<span class="sep" aria-hidden="true"></span>
			{fromUnix(novel.last_updated)}
		{/if}
	</p>
</a>

<style>
	.card {
		display: grid;
		gap: var(--sp-2);
		align-content: start;
	}

	.art {
		position: relative;
		transition: transform var(--dur-2) var(--ease);
	}

	.card:hover .art,
	.card:focus-visible .art {
		transform: translateY(-3px);
	}

	.mark {
		position: absolute;
		top: var(--sp-2);
		left: var(--sp-2);
		padding: 0.1rem 0.3rem;
		font-family: var(--font-mono);
		font-size: 0.6rem;
		font-weight: 600;
		letter-spacing: 0.06em;
		color: #ece9e1;
		border-radius: var(--r-1);
		background: #0b0d0ec7;
		backdrop-filter: blur(6px);
	}

	.track {
		position: absolute;
		inset: auto 0 0 0;
		height: 3px;
		background: #0b0d0e8f;
	}

	.fill {
		height: 100%;
		background: var(--accent);
	}

	.title {
		font-size: var(--step-0);
		font-weight: 500;
		line-height: 1.35;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
		transition: color var(--dur-1) var(--ease);
	}

	.card:hover .title,
	.card:focus-visible .title {
		color: var(--accent);
	}

	.meta {
		display: flex;
		align-items: center;
		gap: var(--sp-2);
		font-size: var(--step--1);
		color: var(--ink-3);
	}

	.sep {
		width: 2px;
		height: 2px;
		border-radius: 50%;
		background: currentColor;
		opacity: 0.6;
	}

	.read {
		color: var(--accent);
	}
</style>
