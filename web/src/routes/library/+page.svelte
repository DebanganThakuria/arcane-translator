<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { getNovels } from '$api/novels';
	import NovelGrid from '$components/NovelGrid.svelte';
	import GridSkeleton from '$components/GridSkeleton.svelte';
	import EmptyState from '$components/EmptyState.svelte';
	import Pagination from '$components/Pagination.svelte';
	import { createPagination } from '$utils/pagination.svelte';
	import type { NovelQuery } from '$lib/types';
	import { ArrowClockwise, ArrowsDownUp, MagnifyingGlass, Plus } from 'phosphor-svelte';

	type SortKey = 'last_updated' | 'date_added' | 'title' | 'chapters' | 'last_read';

	const SORTS: { key: SortKey; label: string }[] = [
		{ key: 'last_read', label: 'Recently read' },
		{ key: 'last_updated', label: 'Last updated' },
		{ key: 'date_added', label: 'Date added' },
		{ key: 'title', label: 'Title' },
		{ key: 'chapters', label: 'Chapter count' }
	];

	/*
	 * Filters live in the query string, so a filtered view can be bookmarked and
	 * survives a reload. They replace history rather than pushing, so the back
	 * button leaves the library instead of unwinding filter changes one by one.
	 */
	function param<T extends string>(key: string, fallback: T): T {
		return (page.url.searchParams.get(key) as T) || fallback;
	}

	const search = $derived(param('q', ''));
	const status = $derived(param<'all' | 'ongoing' | 'completed'>('status', 'all'));
	const language = $derived(param<'all' | 'chinese' | 'korean' | 'japanese' | 'english'>('lang', 'all'));
	const sortBy = $derived(param<SortKey>('sort', 'last_read'));
	const descending = $derived(param('dir', 'desc') === 'desc');

	// The server does the filtering, so the query is part of the fetch.
	const query = $derived<NovelQuery>({
		q: search.trim() || undefined,
		status: status === 'all' ? undefined : status,
		language: language === 'all' ? undefined : language,
		sort: sortBy,
		dir: descending ? 'desc' : 'asc'
	});

	const pagination = createPagination((pageNumber, limit) =>
		getNovels(pageNumber, limit, query)
	);

	$effect(() => {
		// Any change to paging or to a filter re-runs the query.
		pagination.current;
		pagination.limit;
		query;
		pagination.load();
	});

	const activeFilters = $derived(
		Number(search.trim() !== '') + Number(status !== 'all') + Number(language !== 'all')
	);

	function setParam(key: string, value: string, fallback: string) {
		const params = new URLSearchParams(page.url.searchParams);
		if (value === fallback) params.delete(key);
		else params.set(key, value);
		// A narrowed result set may have fewer pages than the current one.
		params.delete('page');

		const next = params.toString();
		goto(`${page.url.pathname}${next ? `?${next}` : ''}`, {
			replaceState: true,
			noScroll: true,
			keepFocus: true
		});
	}

	/** Debounced so typing does not fire a request per keystroke. */
	let searchTimer: ReturnType<typeof setTimeout>;
	function onSearchInput(value: string) {
		clearTimeout(searchTimer);
		searchTimer = setTimeout(() => setParam('q', value, ''), 250);
	}

	function clearFilters() {
		const params = new URLSearchParams(page.url.searchParams);
		for (const key of ['q', 'status', 'lang', 'page']) params.delete(key);
		const next = params.toString();
		goto(`${page.url.pathname}${next ? `?${next}` : ''}`, {
			replaceState: true,
			noScroll: true
		});
	}
</script>

<svelte:head><title>Library · Arcane Translator</title></svelte:head>

<div class="page">
	<header>
		<div>
			<h1 class="display">Library</h1>
			<p class="meta">
				<span class="numeral">{pagination.data.total_count}</span>
				{pagination.data.total_count === 1 ? 'novel' : 'novels'}
			</p>
		</div>

		<button
			class="btn"
			disabled={pagination.loading || pagination.refreshing}
			onclick={() => pagination.refresh()}
		>
			<ArrowClockwise size={14} weight="bold" class={pagination.refreshing ? 'spin' : ''} />
			{pagination.refreshing ? 'Refreshing…' : 'Refresh'}
		</button>
	</header>

	<div class="toolbar">
		<div class="search">
			<MagnifyingGlass size={15} />
			<input
				value={search}
				oninput={(event) => onSearchInput(event.currentTarget.value)}
				type="search"
				placeholder="Title, author or genre…"
				aria-label="Filter novels"
			/>
		</div>

		<label class="control">
			<span>Status</span>
			<select
				class="select"
				value={status}
				onchange={(event) => setParam('status', event.currentTarget.value, 'all')}
			>
				<option value="all">Any</option>
				<option value="ongoing">Ongoing</option>
				<option value="completed">Completed</option>
			</select>
		</label>

		<label class="control">
			<span>Language</span>
			<select
				class="select"
				value={language}
				onchange={(event) => setParam('lang', event.currentTarget.value, 'all')}
			>
				<option value="all">Any</option>
				<option value="chinese">Chinese</option>
				<option value="korean">Korean</option>
				<option value="japanese">Japanese</option>
				<option value="english">English</option>
			</select>
		</label>

		<label class="control">
			<span>Sort</span>
			<select
				class="select"
				value={sortBy}
				onchange={(event) => setParam('sort', event.currentTarget.value, 'last_read')}
			>
				{#each SORTS as sort (sort.key)}
					<option value={sort.key}>{sort.label}</option>
				{/each}
			</select>
		</label>

		<button
			class="btn btn--icon"
			onclick={() => setParam('dir', descending ? 'asc' : 'desc', 'desc')}
			aria-label={descending ? 'Sort ascending' : 'Sort descending'}
			title={descending ? 'Sorted descending' : 'Sorted ascending'}
		>
			<ArrowsDownUp size={15} weight="bold" />
		</button>
	</div>

	{#if pagination.loading}
		<GridSkeleton count={12} />
	{:else if pagination.failed}
		<EmptyState
			title="Could not load your library"
			body="The translation server did not respond. Check that the backend is running on port 8088."
		>
			{#snippet action()}
				<button class="btn btn--primary" onclick={() => pagination.load()}>Try again</button>
			{/snippet}
		</EmptyState>
	{:else if pagination.data.novels.length === 0}
		{#if activeFilters > 0}
			<EmptyState
				title="Nothing matches those filters"
				body="No novel in your library matches every filter you have set."
			>
				{#snippet action()}
					<button class="btn btn--primary" onclick={clearFilters}>Clear filters</button>
				{/snippet}
			</EmptyState>
		{:else}
			<EmptyState
				title="Your library is empty"
				body="Add a novel by pasting its URL from a source site."
			>
				{#snippet action()}
					<a class="btn btn--primary" href="/add">
						<Plus size={15} weight="bold" />
						Add your first novel
					</a>
				{/snippet}
			</EmptyState>
		{/if}
	{:else}
		{#if activeFilters > 0}
			<p class="count meta">
				<span class="numeral">{pagination.data.total_count}</span>
				{pagination.data.total_count === 1 ? 'novel matches' : 'novels match'}
				<button class="clear" onclick={clearFilters}>Clear filters</button>
			</p>
		{/if}

		<NovelGrid novels={pagination.data.novels} />
		<Pagination {pagination} />
	{/if}
</div>

<style>
	header {
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
		gap: var(--sp-4);
		padding-bottom: var(--sp-4);
		border-bottom: 1px solid var(--line);
	}

	h1 {
		font-size: var(--step-3);
	}

	.toolbar {
		display: flex;
		flex-wrap: wrap;
		align-items: flex-end;
		gap: var(--sp-3);
		padding-block: var(--sp-4);
		margin-bottom: var(--sp-5);
	}

	.search {
		display: flex;
		align-items: center;
		gap: var(--sp-2);
		flex: 1 1 15rem;
		min-width: 0;
		height: var(--control-h);
		padding-inline: var(--sp-3);
		background: var(--surface-2);
		border: 1px solid var(--line-strong);
		border-radius: var(--r-1);
		color: var(--ink-3);
	}

	.search:focus-within {
		border-color: var(--accent);
		box-shadow: 0 0 0 3px var(--accent-wash);
	}

	.search input {
		flex: 1;
		min-width: 0;
		border: 0;
		background: none;
		outline: none;
	}

	.search input::-webkit-search-cancel-button {
		display: none;
	}

	.control {
		display: grid;
		gap: var(--sp-1);
	}

	.control span {
		font-size: var(--step--1);
		color: var(--ink-3);
	}

	.control .select {
		width: auto;
		min-width: 8rem;
	}

	.count {
		display: flex;
		align-items: center;
		gap: var(--sp-3);
		margin-bottom: var(--sp-4);
	}

	.clear {
		background: none;
		border: 0;
		padding: 0;
		font-size: inherit;
		color: var(--accent);
		cursor: pointer;
		text-decoration: underline;
		text-underline-offset: 2px;
	}

	@media (max-width: 640px) {
		.control,
		.search {
			flex: 1 1 100%;
		}

		.control .select {
			width: 100%;
		}
	}
</style>
