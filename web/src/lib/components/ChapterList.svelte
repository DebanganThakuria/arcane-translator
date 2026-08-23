<script lang="ts">
	import { getChapterPage } from '$api/novels';
	import { progress } from '$stores/progress.svelte';
	import { fromUnix } from '$utils/time';
	import { pageWindow } from '$utils/pagination.svelte';
	import {
		ArrowDown,
		ArrowUp,
		CaretLeft,
		CaretRight,
		MagnifyingGlass
	} from 'phosphor-svelte';
	import type { ChapterPage } from '$lib/types';

	interface Props {
		novelId: string;
		/** Total chapters, so the header is correct before the first page lands. */
		totalCount: number;
	}

	let { novelId, totalCount }: Props = $props();

	const PER_PAGE = 50;

	let ascending = $state(true);
	let search = $state('');
	let searchTerm = $state('');
	let pageNumber = $state(1);

	let data = $state<ChapterPage | null>(null);
	let loading = $state(true);
	let failed = $state(false);

	const current = $derived(progress.get(novelId));
	const pages = $derived(pageWindow(pageNumber, data?.total_pages ?? 1));

	/** Debounced so typing does not fire a request per keystroke. */
	let searchTimer: ReturnType<typeof setTimeout>;
	function onSearchInput(value: string) {
		search = value;
		clearTimeout(searchTimer);
		searchTimer = setTimeout(() => {
			searchTerm = value.trim();
			pageNumber = 1;
		}, 250);
	}

	// Chapters are paged in the database: a novel here can run past 900 of them,
	// and loading every row to filter in the browser was the slow path.
	$effect(() => {
		const request = { page: pageNumber, term: searchTerm, asc: ascending };
		let cancelled = false;

		loading = true;
		failed = false;

		getChapterPage(novelId, request.page, PER_PAGE, {
			q: request.term || undefined,
			dir: request.asc ? 'asc' : 'desc'
		})
			.then((result) => {
				if (!cancelled) data = result;
			})
			.catch(() => {
				if (!cancelled) {
					failed = true;
					data = null;
				}
			})
			.finally(() => {
				if (!cancelled) loading = false;
			});

		return () => {
			cancelled = true;
		};
	});

	function goTo(next: number) {
		if (data && next >= 1 && next <= data.total_pages) pageNumber = next;
	}
</script>

<div class="head">
	<div class="search">
		<MagnifyingGlass size={14} />
		<input
			value={search}
			oninput={(event) => onSearchInput(event.currentTarget.value)}
			type="search"
			placeholder="Filter chapters…"
			aria-label="Filter chapters"
		/>
	</div>

	<button
		class="btn btn--sm"
		onclick={() => {
			ascending = !ascending;
			pageNumber = 1;
		}}
	>
		{#if ascending}
			<ArrowUp size={13} weight="bold" /> Oldest first
		{:else}
			<ArrowDown size={13} weight="bold" /> Newest first
		{/if}
	</button>
</div>

{#if failed}
	<p class="none">Could not load chapters.</p>
{:else if loading && !data}
	<ol class="list" aria-hidden="true">
		{#each { length: 8 } as _, i (i)}
			<li class="row-skeleton"><span class="skeleton"></span></li>
		{/each}
	</ol>
{:else if data && data.chapters.length === 0}
	<p class="none">
		{searchTerm ? 'No chapters match that filter.' : 'No chapters translated yet.'}
	</p>
{:else if data}
	{#if searchTerm}
		<p class="matches meta">
			<span class="numeral">{data.total_count}</span>
			{data.total_count === 1 ? 'match' : 'matches'}
		</p>
	{/if}

	<ol class="list" class:dimmed={loading}>
		{#each data.chapters as chapter (chapter.id)}
			{@const isCurrent = current?.lastChapter === chapter.number}
			<li>
				<a href="/novel/{novelId}/chapter/{chapter.number}" class:current={isCurrent}>
					<span class="num numeral">{chapter.number}</span>
					<span class="title">{chapter.title}</span>
					{#if isCurrent}
						<span class="badge">{Math.round(current.progress)}%</span>
					{/if}
					<span class="when">{fromUnix(chapter.date_translated)}</span>
				</a>
				{#if isCurrent && current.progress > 0}
					<div class="track" aria-hidden="true">
						<div class="fill" style:width="{current.progress}%"></div>
					</div>
				{/if}
			</li>
		{/each}
	</ol>

	{#if data.total_pages > 1}
		<nav class="pager" aria-label="Chapter pages">
			<button
				class="btn btn--sm btn--icon"
				disabled={pageNumber <= 1 || loading}
				onclick={() => goTo(pageNumber - 1)}
				aria-label="Previous page"
			>
				<CaretLeft size={13} weight="bold" />
			</button>

			<ol class="pages">
				{#each pages as candidate, i (`${candidate}-${i}`)}
					{#if candidate === -1}
						<li class="gap" aria-hidden="true">&hellip;</li>
					{:else}
						<li>
							<button
								class="page numeral"
								class:selected={candidate === pageNumber}
								aria-current={candidate === pageNumber ? 'page' : undefined}
								disabled={loading}
								onclick={() => goTo(candidate)}
							>
								{candidate}
							</button>
						</li>
					{/if}
				{/each}
			</ol>

			<button
				class="btn btn--sm btn--icon"
				disabled={pageNumber >= data.total_pages || loading}
				onclick={() => goTo(pageNumber + 1)}
				aria-label="Next page"
			>
				<CaretRight size={13} weight="bold" />
			</button>

			<p class="tally meta">
				<span class="numeral">{totalCount || data.total_count}</span> chapters
			</p>
		</nav>
	{/if}
{/if}

<style>
	.head {
		display: flex;
		align-items: center;
		gap: var(--sp-3);
		margin-bottom: var(--sp-3);
	}

	.search {
		display: flex;
		align-items: center;
		gap: var(--sp-2);
		flex: 1;
		min-width: 0;
		padding-inline: var(--sp-3);
		height: var(--control-h);
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
		font-size: var(--step-0);
	}

	.search input::-webkit-search-cancel-button {
		display: none;
	}

	.matches {
		margin-bottom: var(--sp-2);
	}

	.list {
		display: grid;
		border-top: 1px solid var(--line);
		transition: opacity var(--dur-1) var(--ease);
	}

	/* A page swap keeps the old rows visible but dimmed, so the list does not
	   collapse and shove the page around. */
	.list.dimmed {
		opacity: 0.5;
	}

	li {
		border-bottom: 1px solid var(--line);
	}

	.row-skeleton {
		padding: var(--sp-3) var(--sp-2);
	}

	.row-skeleton .skeleton {
		display: block;
		height: 1rem;
		width: 70%;
	}

	a {
		display: grid;
		grid-template-columns: 2.5rem 1fr auto auto;
		align-items: baseline;
		gap: var(--sp-3);
		padding: var(--sp-3) var(--sp-2);
		transition: background var(--dur-1) var(--ease);
	}

	a:hover {
		background: var(--accent-wash);
	}

	.num {
		color: var(--ink-3);
		text-align: right;
		font-size: 0.75rem;
	}

	.title {
		min-width: 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	a.current .title {
		color: var(--accent);
		font-weight: 600;
	}

	.badge {
		font-family: var(--font-mono);
		font-size: 0.7rem;
		padding: 0.1rem 0.35rem;
		border-radius: var(--r-1);
		background: var(--accent);
		color: var(--accent-ink);
	}

	.when {
		font-size: var(--step--1);
		color: var(--ink-3);
		white-space: nowrap;
	}

	.track {
		height: 2px;
		background: var(--surface-sunk);
	}

	.fill {
		height: 100%;
		background: var(--accent);
	}

	.none {
		padding: var(--sp-6);
		text-align: center;
		color: var(--ink-2);
	}

	/* --- Pager -------------------------------------------------------------- */
	.pager {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: var(--sp-2);
		padding-top: var(--sp-4);
	}

	.pages {
		display: flex;
		gap: var(--sp-1);
	}

	.page {
		min-width: var(--control-h-sm);
		height: var(--control-h-sm);
		padding-inline: var(--sp-2);
		background: transparent;
		border: 1px solid transparent;
		border-radius: var(--r-1);
		color: var(--ink-2);
		cursor: pointer;
	}

	.page:hover:not(:disabled) {
		background: var(--accent-wash);
		color: var(--ink);
	}

	.page.selected {
		background: var(--accent);
		border-color: var(--accent);
		color: var(--accent-ink);
		font-weight: 600;
	}

	.gap {
		display: grid;
		place-items: center;
		min-width: 1.25rem;
		color: var(--ink-3);
	}

	.tally {
		margin-left: auto;
	}

	@media (max-width: 640px) {
		a {
			grid-template-columns: 2rem 1fr auto;
		}

		.when {
			display: none;
		}
	}
</style>
