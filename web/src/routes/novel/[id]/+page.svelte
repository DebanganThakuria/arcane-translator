<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { getChapterPage, getNovel, deleteNovel } from '$api/novels';
	import { refreshNovel } from '$api/translation';
	import Cover from '$components/Cover.svelte';
	import ChapterList from '$components/ChapterList.svelte';
	import ConfirmDialog from '$components/ConfirmDialog.svelte';
	import FirstChapterDialog from '$components/FirstChapterDialog.svelte';
	import EmptyState from '$components/EmptyState.svelte';
	import { progress } from '$stores/progress.svelte';
	import { sources } from '$stores/sources.svelte';
	import { toast } from '$stores/toast.svelte';
	import { ManualExtractCancelled } from '$stores/manualExtract.svelte';
	import { fromUnix } from '$utils/time';
	import { ArrowClockwise, ArrowLeft, ArrowSquareOut, BookOpen, Trash } from 'phosphor-svelte';
	import type { ChapterPage, Novel } from '$lib/types';

	const novelId = $derived(page.params.id ?? '');

	let novel = $state<Novel | null>(null);
	// Only the first page of chapters, plus the number bounds for the whole
	// novel. Loading all of them just to find the first one was wasteful.
	let chapterInfo = $state<ChapterPage | null>(null);
	let loading = $state(true);
	let loadError = $state('');

	let refreshing = $state(false);
	let deleting = $state(false);
	let confirmOpen = $state(false);
	let firstChapterOpen = $state(false);

	async function load() {
		loading = true;
		loadError = '';
		try {
			const [fetchedNovel, fetchedChapters] = await Promise.all([
				getNovel(novelId),
				getChapterPage(novelId, 1, 1, { dir: 'asc' }).catch(() => null)
			]);
			novel = fetchedNovel;
			chapterInfo = fetchedChapters;
		} catch (cause) {
			loadError = cause instanceof Error ? cause.message : 'Could not load this novel.';
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		novelId;
		load();
	});

	const saved = $derived(progress.get(novelId));

	const translatedCount = $derived(chapterInfo?.total_count ?? 0);
	const firstChapter = $derived(chapterInfo?.first_number || null);
	const lastChapter = $derived(chapterInfo?.last_number || null);

	/**
	 * Where "continue" goes: local progress first, then the server's record,
	 * clamped to the range the novel actually has.
	 */
	const resumeChapter = $derived.by(() => {
		if (firstChapter === null || lastChapter === null) return null;
		const candidate = saved?.lastChapter || novel?.last_read_chapter_number || 0;
		if (candidate >= firstChapter && candidate <= lastChapter) return candidate;
		return firstChapter;
	});

	const started = $derived(Boolean(saved?.lastChapter || novel?.last_read_timestamp));

	async function refresh() {
		if (!novel) return;
		refreshing = true;
		try {
			await refreshNovel(novel.id, novel.url);
			await load();
			toast.success('Novel refreshed');
		} catch (cause) {
			if (!(cause instanceof ManualExtractCancelled)) {
				toast.fromError('Refresh failed', cause);
			}
		} finally {
			refreshing = false;
		}
	}

	async function remove() {
		deleting = true;
		try {
			await deleteNovel(novelId);
			toast.success('Novel deleted', novel?.title);
			await goto('/library');
		} catch (cause) {
			toast.fromError('Delete failed', cause);
			deleting = false;
			confirmOpen = false;
		}
	}
</script>

<svelte:head><title>{novel?.title ?? 'Novel'} · Arcane Translator</title></svelte:head>

<div class="page">
	<a class="back" href="/library">
		<ArrowLeft size={14} weight="bold" />
		Library
	</a>

	{#if loading}
		<div class="layout">
			<div class="rail">
				<div class="skeleton art"></div>
			</div>
			<div class="main">
				<div class="skeleton line" style:width="60%" style:height="2rem"></div>
				<div class="skeleton line" style:width="35%"></div>
				<div class="skeleton line"></div>
				<div class="skeleton line" style:width="85%"></div>
			</div>
		</div>
	{:else if loadError || !novel}
		<EmptyState title="Novel not found" body={loadError || 'This novel is no longer in your library.'}>
			{#snippet action()}
				<a class="btn btn--primary" href="/library">Back to library</a>
			{/snippet}
		</EmptyState>
	{:else}
		<div class="layout">
			<aside class="rail">
				<Cover {novel} eager />

				<div class="actions">
					{#if resumeChapter !== null}
						<a class="btn btn--primary btn--block" href="/novel/{novel.id}/chapter/{resumeChapter}">
							<BookOpen size={15} weight="bold" />
							{started ? `Continue chapter ${resumeChapter}` : 'Start reading'}
						</a>
						{#if started && firstChapter !== null && firstChapter !== resumeChapter}
							<a class="btn btn--block" href="/novel/{novel.id}/chapter/{firstChapter}">
								Start from the beginning
							</a>
						{/if}
					{/if}
				</div>

				{#if started && saved}
					<p class="resume-note meta">
						Chapter <span class="numeral">{saved.lastChapter}</span>, {Math.round(saved.progress)}%
						read, {fromUnix(Math.floor(saved.lastReadAt / 1000))}
					</p>
				{/if}

				<dl class="facts">
					<dt>Status</dt>
					<dd>{novel.status ?? 'Unknown'}</dd>

					<dt>Chapters</dt>
					<dd>
						<span class="numeral">{translatedCount}</span> translated
						{#if novel.chapters_count > translatedCount}
							<span class="of">of {novel.chapters_count}</span>
						{/if}
					</dd>

					<dt>Source</dt>
					<dd>{sources.byId(novel.source)?.name ?? novel.source}</dd>

					<dt>Updated</dt>
					<dd>{fromUnix(novel.last_updated)}</dd>
				</dl>

				{#if novel.genres?.length}
					<div class="genres">
						{#each novel.genres as genre (genre)}
							<a class="chip" href="/genre/{encodeURIComponent(genre)}">{genre}</a>
						{/each}
					</div>
				{/if}

				<div class="admin">
					<button class="btn btn--sm" onclick={refresh} disabled={refreshing}>
						<ArrowClockwise size={13} weight="bold" class={refreshing ? 'spin' : ''} />
						{refreshing ? 'Checking…' : 'Check for updates'}
					</button>
					<button class="btn btn--sm btn--danger" onclick={() => (confirmOpen = true)}>
						<Trash size={13} weight="bold" />
						Delete
					</button>
				</div>
			</aside>

			<div class="main">
				<h1 class="display">{novel.title}</h1>

				<p class="origin">
					{#if novel.original_title}<span class="original">{novel.original_title}</span>{/if}
					{#if novel.author}<span class="author">by {novel.author}</span>{/if}
					<a class="source-link" href={novel.url} target="_blank" rel="noreferrer noopener">
						Source
						<ArrowSquareOut size={12} weight="bold" />
					</a>
				</p>

				{#if novel.summary}
					<!-- Summary HTML is produced by our own translation backend. -->
					<div class="prose summary">{@html novel.summary}</div>
				{/if}

				<section class="chapters">
					<h2 class="display">Chapters</h2>

					{#if translatedCount === 0}
						<EmptyState
							title="No chapters yet"
							body="Point the translator at the first chapter and it will follow the links from there."
						>
							{#snippet action()}
								<button class="btn btn--primary" onclick={() => (firstChapterOpen = true)}>
									Set the first chapter
								</button>
							{/snippet}
						</EmptyState>
					{:else}
						<ChapterList novelId={novel.id} totalCount={translatedCount} />
					{/if}
				</section>
			</div>
		</div>

		<ConfirmDialog
			open={confirmOpen}
			title="Delete this novel?"
			description="{novel.title} and its {translatedCount} translated {translatedCount === 1
				? 'chapter'
				: 'chapters'} will be removed permanently."
			busy={deleting}
			onconfirm={remove}
			oncancel={() => (confirmOpen = false)}
		/>

		<FirstChapterDialog
			open={firstChapterOpen}
			novelId={novel.id}
			novelUrl={novel.url}
			onclose={() => (firstChapterOpen = false)}
			onsuccess={(chapterNumber) => {
				firstChapterOpen = false;
				toast.success('First chapter ready');
				goto(`/novel/${novelId}/chapter/${chapterNumber}`);
			}}
		/>
	{/if}
</div>

<style>
	.back {
		display: inline-flex;
		align-items: center;
		gap: var(--sp-2);
		margin-bottom: var(--sp-5);
		font-size: var(--step--1);
		color: var(--ink-2);
	}

	.back:hover {
		color: var(--accent);
	}

	.layout {
		display: grid;
		grid-template-columns: 17rem minmax(0, 1fr);
		gap: var(--sp-7);
		align-items: start;
	}

	/* --- Rail ------------------------------------------------------------- */
	.rail {
		display: grid;
		gap: var(--sp-4);
		position: sticky;
		top: calc(68px + var(--sp-5));
	}

	.rail .art {
		aspect-ratio: 2 / 3;
		border-radius: var(--r-2);
	}

	.actions {
		display: grid;
		gap: var(--sp-2);
	}

	.resume-note {
		margin-top: calc(var(--sp-3) * -1);
		color: var(--ink-3);
	}

	.facts {
		display: grid;
		grid-template-columns: auto 1fr;
		gap: var(--sp-2) var(--sp-4);
		padding-top: var(--sp-4);
		border-top: 1px solid var(--line);
		font-size: var(--step--1);
	}

	.facts dt {
		color: var(--ink-3);
	}

	.facts dd {
		color: var(--ink);
		text-align: right;
	}

	.of {
		color: var(--ink-3);
	}

	.genres {
		display: flex;
		flex-wrap: wrap;
		gap: var(--sp-2);
	}

	.admin {
		display: flex;
		flex-wrap: wrap;
		gap: var(--sp-2);
		padding-top: var(--sp-4);
		border-top: 1px solid var(--line);
	}

	/* --- Main ------------------------------------------------------------- */
	.main {
		display: grid;
		gap: var(--sp-4);
	}

	h1 {
		font-size: clamp(2rem, 4.4vw, 3rem);
	}

	.origin {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: var(--sp-2) var(--sp-4);
		font-size: var(--step-0);
		color: var(--ink-2);
	}

	.original {
		font-family: var(--font-display);
		font-size: var(--step-1);
	}

	.author {
		color: var(--ink);
	}

	.source-link {
		display: inline-flex;
		align-items: center;
		gap: var(--sp-1);
		font-size: var(--step--1);
		color: var(--accent);
	}

	.source-link:hover {
		text-decoration: underline;
		text-underline-offset: 2px;
	}

	.summary {
		max-width: var(--measure);
		color: var(--ink-2);
		padding-block: var(--sp-2);
	}

	.chapters {
		margin-top: var(--sp-5);
	}

	.chapters h2 {
		font-size: var(--step-2);
		padding-bottom: var(--sp-3);
		margin-bottom: var(--sp-4);
		border-bottom: 1px solid var(--line);
	}

	.line {
		height: 1rem;
	}

	/* --- Narrow ----------------------------------------------------------- */
	@media (max-width: 880px) {
		.layout {
			grid-template-columns: minmax(0, 1fr);
			gap: var(--sp-6);
		}

		.rail {
			position: static;
			grid-template-columns: 9rem minmax(0, 1fr);
			align-items: start;
			column-gap: var(--sp-5);
		}

		.rail :global(> *:first-child) {
			grid-row: span 4;
		}

		.facts {
			border-top: 0;
			padding-top: 0;
		}

		.facts dd {
			text-align: left;
		}
	}

	@media (max-width: 560px) {
		.rail {
			grid-template-columns: minmax(0, 1fr);
		}

		.rail :global(> *:first-child) {
			grid-row: auto;
			max-width: 11rem;
		}
	}
</style>
