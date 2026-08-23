<script lang="ts">
	import { page } from '$app/state';
	import { searchNovels } from '$api/novels';
	import NovelGrid from '$components/NovelGrid.svelte';
	import GridSkeleton from '$components/GridSkeleton.svelte';
	import EmptyState from '$components/EmptyState.svelte';
	import type { Novel } from '$lib/types';

	const query = $derived(page.url.searchParams.get('q')?.trim() ?? '');

	let results = $state<Novel[]>([]);
	let loading = $state(false);
	let failed = $state(false);

	$effect(() => {
		const term = query;
		if (!term) {
			results = [];
			loading = false;
			failed = false;
			return;
		}

		let cancelled = false;
		loading = true;
		failed = false;

		searchNovels(term)
			.then((found) => {
				if (!cancelled) results = found;
			})
			.catch(() => {
				if (!cancelled) {
					failed = true;
					results = [];
				}
			})
			.finally(() => {
				if (!cancelled) loading = false;
			});

		return () => {
			cancelled = true;
		};
	});
</script>

<svelte:head><title>{query ? `${query} · Search` : 'Search'} · Arcane Translator</title></svelte:head>

<div class="page">
	<header>
		<h1 class="display">Search</h1>
		{#if query}
			<p class="meta">
				{#if loading}
					Searching for &ldquo;{query}&rdquo;…
				{:else if failed}
					Search failed
				{:else}
					<span class="numeral">{results.length}</span>
					{results.length === 1 ? 'result' : 'results'} for &ldquo;{query}&rdquo;
				{/if}
			</p>
		{/if}
	</header>

	{#if !query}
		<EmptyState
			title="Search your library"
			body="Press Command K anywhere to open search, or type a title, author or genre."
		/>
	{:else if loading}
		<GridSkeleton count={10} />
	{:else if failed}
		<EmptyState
			title="Search failed"
			body="The translation server did not respond. Check that the backend is running on port 8088."
		/>
	{:else if results.length === 0}
		<EmptyState
			title="No matches"
			body="Nothing in your library matches &ldquo;{query}&rdquo;. Try a shorter term."
		/>
	{:else}
		<NovelGrid novels={results} />
	{/if}
</div>

<style>
	header {
		padding-bottom: var(--sp-4);
		margin-bottom: var(--sp-6);
		border-bottom: 1px solid var(--line);
	}

	h1 {
		font-size: var(--step-3);
	}

	.meta {
		margin-top: var(--sp-1);
	}
</style>
