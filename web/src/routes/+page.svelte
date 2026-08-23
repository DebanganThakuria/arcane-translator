<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { getRecentlyRead, getNovelsByFilter, getStats } from '$api/novels';
	import Cover from '$components/Cover.svelte';
	import NovelGrid from '$components/NovelGrid.svelte';
	import GridSkeleton from '$components/GridSkeleton.svelte';
	import EmptyState from '$components/EmptyState.svelte';
	import { progress } from '$stores/progress.svelte';
	import { formatCount, fromUnix } from '$utils/time';
	import { ArrowRight, BookOpen, Plus } from 'phosphor-svelte';
	import type { Novel, NovelStats } from '$lib/types';

	const LANGUAGES = [
		{ key: 'chinese', label: 'Chinese', href: '/chinese' },
		{ key: 'korean', label: 'Korean', href: '/korean' },
		{ key: 'japanese', label: 'Japanese', href: '/japanese' }
	] as const;

	let stats = $state<NovelStats>({ novel_count: 0, chapter_count: 0 });
	let recent = $state<Novel[]>([]);
	let loading = $state(true);

	type LanguageKey = (typeof LANGUAGES)[number]['key'];

	// In the URL so the chosen shelf survives a reload and can be linked to.
	const active = $derived.by<LanguageKey>(() => {
		const value = page.url.searchParams.get('lang') as LanguageKey | null;
		return LANGUAGES.some((language) => language.key === value) ? value! : 'chinese';
	});

	function selectLanguage(key: LanguageKey) {
		const query = key === 'chinese' ? '' : `?lang=${key}`;
		goto(`/${query}`, { replaceState: true, noScroll: true, keepFocus: true });
	}

	let shelves = $state<Record<string, Novel[]>>({});
	let shelvesLoading = $state(true);

	$effect(() => {
		let cancelled = false;

		Promise.all([
			getStats().catch(() => ({ novel_count: 0, chapter_count: 0 })),
			getRecentlyRead(7).catch(() => ({ novels: [] }))
		]).then(([fetchedStats, fetchedRecent]) => {
			if (cancelled) return;
			stats = fetchedStats;
			recent = fetchedRecent.novels;
			loading = false;
		});

		Promise.all(
			LANGUAGES.map((language) =>
				getNovelsByFilter('language', language.key, 1, 12)
					.then((result) => [language.key, result.novels] as const)
					.catch(() => [language.key, [] as Novel[]] as const)
			)
		).then((entries) => {
			if (cancelled) return;
			shelves = Object.fromEntries(entries);
			shelvesLoading = false;
		});

		return () => {
			cancelled = true;
		};
	});

	/** The most recently read novel gets the resume panel; the rest form a shelf. */
	const resuming = $derived(recent[0]);
	const rest = $derived(recent.slice(1));

	const resumeChapter = $derived(
		resuming ? (progress.get(resuming.id)?.lastChapter ?? resuming.last_read_chapter_number ?? 1) : 1
	);

	const resumePercent = $derived(resuming ? Math.round(progress.get(resuming.id)?.progress ?? 0) : 0);
	const resumeTitle = $derived(resuming ? progress.get(resuming.id)?.chapterTitle : undefined);
</script>

<div class="page">
	{#if loading}
		<div class="resume skeleton-panel">
			<div class="skeleton art"></div>
			<div class="lines">
				<div class="skeleton line" style:width="45%"></div>
				<div class="skeleton line" style:width="70%"></div>
				<div class="skeleton line" style:width="30%"></div>
			</div>
		</div>
	{:else if resuming}
		<!-- Resume is the reason this page exists: get back into the book. -->
		<section class="resume">
			<a class="art" href="/novel/{resuming.id}">
				<Cover novel={resuming} eager />
			</a>

			<div class="detail">
				<p class="kicker">Continue reading</p>
				<h1 class="display">{resuming.title}</h1>
				<p class="chapter">
					Chapter <span class="numeral">{resumeChapter}</span>{#if resumeTitle}<span
							class="chapter-title">{resumeTitle}</span>{/if}
					<span class="when">
						read {fromUnix(resuming.last_read_timestamp || resuming.last_updated)}
					</span>
				</p>

				{#if resumePercent > 0}
					<div class="progress">
						<div class="track"><div class="fill" style:width="{resumePercent}%"></div></div>
						<span class="numeral">{resumePercent}%</span>
					</div>
				{/if}

				<div class="actions">
					<a class="btn btn--primary" href="/novel/{resuming.id}/chapter/{resumeChapter}">
						<BookOpen size={15} weight="bold" />
						Resume chapter {resumeChapter}
					</a>
					<a class="btn" href="/novel/{resuming.id}">All chapters</a>
				</div>
			</div>

		</section>

		<p class="library-summary meta">
			Your library holds <span class="numeral">{formatCount(stats.novel_count)}</span> novels and
			<span class="numeral">{formatCount(stats.chapter_count)}</span> translated chapters.
		</p>
	{:else}
		<section class="welcome">
			<h1 class="display">Read webnovels in your own language</h1>
			<p>
				Point Arcane Translator at a Chinese, Korean or Japanese novel and it translates chapters as
				you read them.
			</p>
			<a class="btn btn--primary" href="/add">
				<Plus size={15} weight="bold" />
				Add your first novel
			</a>
		</section>
	{/if}

	{#if rest.length > 0}
		<section class="shelf">
			<header>
				<h2 class="display">Recently read</h2>
				<a class="more" href="/library">
					Library
					<ArrowRight size={13} weight="bold" />
				</a>
			</header>
			<NovelGrid novels={rest} eagerCount={0} />
		</section>
	{/if}

	<section class="shelf">
		<header>
			<h2 class="display">Browse by language</h2>
			<a class="more" href={LANGUAGES.find((language) => language.key === active)?.href ?? '/library'}>
				View all
				<ArrowRight size={13} weight="bold" />
			</a>
		</header>

		<div class="tabs" role="tablist" aria-label="Language">
			{#each LANGUAGES as language (language.key)}
				<button
					role="tab"
					aria-selected={active === language.key}
					class:active={active === language.key}
					onclick={() => selectLanguage(language.key)}
				>
					{language.label}
				</button>
			{/each}
		</div>

		{#if shelvesLoading}
			<GridSkeleton count={6} />
		{:else if (shelves[active] ?? []).length > 0}
			<NovelGrid novels={shelves[active]} eagerCount={0} />
		{:else}
			<EmptyState title="Nothing here yet" body="No {active} novels in your library.">
				{#snippet action()}
					<a class="btn btn--primary" href="/add">
						<Plus size={15} weight="bold" />
						Add a novel
					</a>
				{/snippet}
			</EmptyState>
		{/if}
	</section>
</div>

<style>
	/* --- Resume ---------------------------------------------------------- */
	.resume {
		display: grid;
		grid-template-columns: 8.5rem minmax(0, 1fr);
		gap: var(--sp-6);
		align-items: center;
		padding: var(--sp-5);
		background: var(--surface);
		border: 1px solid var(--line);
		border-radius: var(--r-2);
	}

	.resume .art {
		width: 8.5rem;
	}

	.detail {
		display: grid;
		gap: var(--sp-2);
		justify-items: start;
	}

	.kicker {
		font-size: var(--step--1);
		color: var(--accent);
	}

	.detail h1 {
		font-size: clamp(1.75rem, 3.4vw, var(--step-4));
	}

	.chapter {
		color: var(--ink-2);
	}

	.chapter-title::before {
		content: '·';
		margin-inline: 0.5ch;
	}

	.progress {
		display: flex;
		align-items: center;
		gap: var(--sp-3);
		width: min(100%, 22rem);
		font-size: var(--step--1);
		color: var(--ink-3);
	}

	.track {
		flex: 1;
		height: 3px;
		background: var(--surface-sunk);
		border-radius: 999px;
		overflow: hidden;
	}

	.fill {
		height: 100%;
		background: var(--accent);
	}

	.actions {
		display: flex;
		flex-wrap: wrap;
		gap: var(--sp-2);
		margin-top: var(--sp-2);
	}

	.when {
		color: var(--ink-3);
	}

	.when::before {
		content: '·';
		margin-inline: 0.5ch;
	}

	/* Deliberately outside the panel: these are library totals, and inside the
	   card they read as facts about the novel. */
	.library-summary {
		margin-top: var(--sp-4);
		color: var(--ink-3);
	}

	/* --- Skeleton -------------------------------------------------------- */
	.skeleton-panel {
		grid-template-columns: 8.5rem 1fr;
	}

	.skeleton-panel .art {
		aspect-ratio: 2 / 3;
		border-radius: var(--r-2);
	}

	.lines {
		display: grid;
		gap: var(--sp-3);
		align-content: center;
	}

	.line {
		height: 1rem;
	}

	/* --- First run -------------------------------------------------------- */
	.welcome {
		display: grid;
		justify-items: start;
		gap: var(--sp-4);
		padding: var(--sp-7) 0;
		max-width: 34rem;
	}

	.welcome h1 {
		font-size: clamp(2rem, 5vw, 3.25rem);
	}

	.welcome p {
		color: var(--ink-2);
		font-size: var(--step-1);
		text-wrap: pretty;
	}

	/* --- Shelves ---------------------------------------------------------- */
	.shelf {
		margin-top: var(--sp-8);
	}

	.shelf header {
		display: flex;
		align-items: baseline;
		justify-content: space-between;
		gap: var(--sp-4);
		padding-bottom: var(--sp-3);
		margin-bottom: var(--sp-5);
		border-bottom: 1px solid var(--line);
	}

	.more {
		display: inline-flex;
		align-items: center;
		gap: var(--sp-2);
		font-size: var(--step--1);
		color: var(--ink-2);
	}

	.more:hover {
		color: var(--accent);
	}

	.tabs {
		display: flex;
		gap: var(--sp-1);
		margin-bottom: var(--sp-5);
	}

	.tabs button {
		padding: var(--sp-2) var(--sp-4);
		background: transparent;
		border: 1px solid var(--line);
		border-radius: var(--r-1);
		color: var(--ink-2);
		font-size: var(--step-0);
		cursor: pointer;
		transition:
			background var(--dur-1) var(--ease),
			color var(--dur-1) var(--ease);
	}

	.tabs button:hover {
		color: var(--ink);
	}

	.tabs button.active {
		background: var(--accent);
		border-color: var(--accent);
		color: var(--accent-ink);
		font-weight: 600;
	}

	/* --- Narrow ----------------------------------------------------------- */
	@media (max-width: 900px) {
		.resume {
			grid-template-columns: 7rem minmax(0, 1fr);
		}

		.resume .art {
			width: 7rem;
		}

	}

	@media (max-width: 560px) {
		.resume {
			grid-template-columns: 1fr;
			gap: var(--sp-4);
			justify-items: start;
		}

	}
</style>
