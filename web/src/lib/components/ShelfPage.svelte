<script lang="ts">
	import NovelGrid from './NovelGrid.svelte';
	import GridSkeleton from './GridSkeleton.svelte';
	import EmptyState from './EmptyState.svelte';
	import Pagination from './Pagination.svelte';
	import { createPagination } from '$utils/pagination.svelte';
	import { ArrowClockwise, Plus } from 'phosphor-svelte';
	import type { PaginatedNovels } from '$lib/types';

	interface Props {
		title: string;
		/** Re-runs whenever the fetcher identity changes, e.g. a new genre param. */
		fetcher: (page: number, limit: number) => Promise<PaginatedNovels>;
		emptyTitle: string;
		emptyBody: string;
	}

	let { title, fetcher, emptyTitle, emptyBody }: Props = $props();

	const pagination = createPagination((page, limit) => fetcher(page, limit));

	$effect(() => {
		// Reloads when the URL page changes and when the shelf itself changes,
		// e.g. navigating from one genre to another.
		fetcher;
		pagination.current;
		pagination.limit;
		pagination.load();
	});

	const data = $derived(pagination.data);
</script>

<div class="page">
	<header>
		<div>
			<h1 class="display">{title}</h1>
			{#if !pagination.loading && !pagination.failed}
				<p class="meta">
					<span class="numeral">{data.total_count}</span>
					{data.total_count === 1 ? 'novel' : 'novels'}
				</p>
			{/if}
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

	{#if pagination.loading}
		<GridSkeleton count={12} />
	{:else if pagination.failed}
		<EmptyState
			title="Could not load this shelf"
			body="The translation server did not respond. Check that the backend is running on port 8088."
		>
			{#snippet action()}
				<button class="btn btn--primary" onclick={() => pagination.load()}>Try again</button>
			{/snippet}
		</EmptyState>
	{:else if data.novels.length === 0}
		<EmptyState title={emptyTitle} body={emptyBody}>
			{#snippet action()}
				<a class="btn btn--primary" href="/add">
					<Plus size={15} weight="bold" />
					Add a novel
				</a>
			{/snippet}
		</EmptyState>
	{:else}
		<NovelGrid novels={data.novels} />
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
		margin-bottom: var(--sp-6);
		border-bottom: 1px solid var(--line);
	}

	h1 {
		font-size: var(--step-3);
	}

	:global(.spin) {
		animation: spin 900ms linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
