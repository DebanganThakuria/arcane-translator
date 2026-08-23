<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { getChapterByNumber, getNovel, deleteChapter } from '$api/novels';
	import { translateChapter } from '$api/translation';
	import ReaderSettings from '$components/ReaderSettings.svelte';
	import BookSpread from '$components/BookSpread.svelte';
	import ConfirmDialog from '$components/ConfirmDialog.svelte';
	import { reader, FONT_STACKS } from '$stores/reader.svelte';
	import { progress } from '$stores/progress.svelte';
	import { fullscreen } from '$stores/fullscreen.svelte';
	import { toast } from '$stores/toast.svelte';
	import { ManualExtractCancelled } from '$stores/manualExtract.svelte';
	import { readingMinutes } from '$utils/time';
	import {
		ArrowLeft,
		ArrowRight,
		ArrowsIn,
		ArrowsOut,
		CaretLeft,
		CaretRight,
		ListBullets,
		Sliders,
		Sparkle,
		Trash
	} from 'phosphor-svelte';
	import type { Chapter, Novel } from '$lib/types';

	const novelId = $derived(page.params.id ?? '');
	const chapterNumber = $derived(Number(page.params.number ?? 1));

	let novel = $state<Novel | null>(null);
	let chapter = $state<Chapter | null>(null);
	let hasNext = $state(false);
	let loading = $state(true);
	let loadError = $state('');

	let translating = $state(false);
	let deleting = $state(false);
	let confirmOpen = $state(false);
	let scrollPercent = $state(0);

	/** Index of the visible two-page spread, and how many the chapter fills. */
	let spread = $state(0);
	let spreadCount = $state(1);
	let readerEl = $state<HTMLDivElement>();

	// Keeps the fullscreen button in step, including when Escape exits.
	$effect(() => fullscreen.watch());

	const prefs = $derived(reader.prefs);
	const paged = $derived(prefs.mode === 'paged');
	const minutes = $derived(readingMinutes(chapter?.word_count));

	/** In paged mode position is which spread you are on, not how far you scrolled. */
	const percent = $derived(
		paged
			? spreadCount > 1
				? Math.round((spread / (spreadCount - 1)) * 100)
				: 100
			: scrollPercent
	);

	async function load() {
		loading = true;
		loadError = '';
		scrollPercent = 0;
		spread = 0;
		spreadCount = 1;
		try {
			const [fetchedNovel, fetchedChapter] = await Promise.all([
				getNovel(novelId),
				getChapterByNumber(novelId, chapterNumber)
			]);
			novel = fetchedNovel;
			chapter = fetchedChapter;
			progress.save(novelId, chapterNumber, 0, fetchedChapter.title);
			scroller().scrollTo(0, 0);
			checkNext();
		} catch (cause) {
			loadError = cause instanceof Error ? cause.message : 'Could not load this chapter.';
		} finally {
			loading = false;
		}
	}

	async function checkNext() {
		hasNext = await getChapterByNumber(novelId, chapterNumber + 1)
			.then(() => true)
			.catch(() => false);
	}

	$effect(() => {
		novelId;
		chapterNumber;
		load();
	});

	/**
	 * Going fullscreen makes the reader element the scroll container, because a
	 * fullscreen element is sized to the screen and the document stops
	 * scrolling. Everything that reads or writes scroll position has to follow.
	 */
	function scroller(): HTMLElement {
		return fullscreen.active && readerEl ? readerEl : document.documentElement;
	}

	/**
	 * Outside fullscreen the visible progress bar is driven by CSS
	 * scroll-timeline, so no JavaScript runs per frame for it. This listener
	 * exists to persist the position, and it bails out unless the whole-number
	 * percentage changed, which caps it at about a hundred updates per chapter.
	 */
	$effect(() => {
		if (!chapter || paged) return;

		// Re-attaches when fullscreen toggles, since the scroller changes.
		const element = fullscreen.active && readerEl ? readerEl : null;
		const target: HTMLElement | Window = element ?? window;

		let frame = 0;
		let saveTimer: ReturnType<typeof setTimeout>;

		function measure() {
			frame = 0;
			const { scrollTop, scrollHeight, clientHeight } = element ?? document.documentElement;
			const scrollable = scrollHeight - clientHeight;
			const next = scrollable > 0 ? Math.round((scrollTop / scrollable) * 100) : 0;
			if (next === scrollPercent) return;

			scrollPercent = Math.min(100, Math.max(0, next));
			clearTimeout(saveTimer);
			saveTimer = setTimeout(() => {
				progress.save(novelId, chapterNumber, scrollPercent, chapter?.title);
			}, 600);
		}

		function onScroll() {
			frame ||= requestAnimationFrame(measure);
		}

		target.addEventListener('scroll', onScroll, { passive: true });
		return () => {
			target.removeEventListener('scroll', onScroll);
			cancelAnimationFrame(frame);
			clearTimeout(saveTimer);
			// Persist wherever the reader stopped before the chapter unmounts.
			if (scrollPercent > 0) progress.save(novelId, chapterNumber, scrollPercent, chapter?.title);
		};
	});

	// Paged mode has no scroll to sample, so persist on each page turn instead.
	$effect(() => {
		if (!chapter || !paged) return;
		const at = percent;
		const timer = setTimeout(() => progress.save(novelId, chapterNumber, at, chapter?.title), 400);
		return () => clearTimeout(timer);
	});

	/** Forward past the last spread continues into the next chapter. */
	function turnForward() {
		if (spread < spreadCount - 1) spread += 1;
		else if (hasNext) navigate(1);
	}

	/** Back from the first spread steps into the previous chapter. */
	function turnBack() {
		if (spread > 0) spread -= 1;
		else if (chapterNumber > 1) navigate(-1);
	}

	function navigate(direction: -1 | 1) {
		const target = chapterNumber + direction;
		if (target < 1) return;
		if (direction === 1 && !hasNext) return;
		goto(`/novel/${novelId}/chapter/${target}`);
	}

	/**
	 * Translates the following chapter. Called from the header it runs in the
	 * background so reading is not interrupted, and the chapter is waiting by
	 * the time the reader reaches the end. Called from the footer it advances.
	 */
	async function translateNext(options: { advance?: boolean } = {}) {
		if (!chapter || translating) return;

		translating = true;
		try {
			const next = await translateChapter(novelId, chapter.next_chapter_url);
			hasNext = true;

			if (options.advance) {
				goto(`/novel/${novelId}/chapter/${next.number ?? chapterNumber + 1}`);
			} else {
				toast.success('Next chapter ready', next.title);
			}
		} catch (cause) {
			if (!(cause instanceof ManualExtractCancelled)) {
				toast.fromError('Could not translate the next chapter', cause);
			}
		} finally {
			translating = false;
		}
	}

	async function removeChapter() {
		if (!chapter) return;
		deleting = true;
		try {
			await deleteChapter(novelId, chapter.id);
			progress.clear(novelId);
			toast.success('Chapter deleted');
			await goto(`/novel/${novelId}`);
		} catch (cause) {
			toast.fromError('Delete failed', cause);
			deleting = false;
			confirmOpen = false;
		}
	}

	function onKeydown(event: KeyboardEvent) {
		const target = event.target as HTMLElement | null;
		if (target && /^(INPUT|TEXTAREA|SELECT)$/.test(target.tagName)) return;
		if (event.metaKey || event.ctrlKey || event.altKey) return;

		if (event.key === 'f' || event.key === 'F') {
			event.preventDefault();
			fullscreen.toggle(readerEl);
			return;
		}

		// Paged: arrows and space turn leaves, and run on into the next chapter
		// at the end. Scroll: arrows move between chapters, space is the
		// browser's own page-down and is left alone.
		if (paged) {
			if (event.key === 'ArrowLeft' || event.key === 'PageUp') {
				event.preventDefault();
				turnBack();
			}
			if (event.key === 'ArrowRight' || event.key === 'PageDown' || event.key === ' ') {
				event.preventDefault();
				turnForward();
			}
			return;
		}

		if (event.key === 'ArrowLeft') navigate(-1);
		if (event.key === 'ArrowRight') navigate(1);
	}
</script>

<svelte:head>
	<title>{chapter ? `${chapter.title} · ${novel?.title ?? ''}` : 'Reading'}</title>
</svelte:head>

<svelte:window onkeydown={onKeydown} />

<div
	class="reader"
	bind:this={readerEl}
	data-reader={prefs.theme}
	class:paged
	class:fullscreen={fullscreen.active}
>
	<div class="progress-bar" aria-hidden="true" style:--turned="{percent}%"></div>

	<header>
		<div class="bar">
			<a class="back" href="/novel/{novelId}">
				<CaretLeft size={14} weight="bold" />
				<span class="back-label">{novel?.title ?? 'Back'}</span>
			</a>

			{#if chapter?.word_count}
				<p class="stats"><span class="numeral">{minutes}</span> min read</p>
			{/if}
		</div>
	</header>

	<!--
		The controls that must stay reachable live on the edge rather than in a
		sticky bar, so nothing hovers over the text while reading.
	-->
	<div class="rail">
		{#if chapter && !hasNext}
			<button
				class="tool prefetch"
				onclick={() => translateNext()}
				disabled={translating}
				aria-label="Translate the next chapter in the background"
				title="Translate the next chapter now, while you read this one"
			>
				<Sparkle size={16} weight="bold" class={translating ? 'spin' : ''} />
			</button>
		{/if}

		{#if fullscreen.supported}
			<button
				class="tool"
				onclick={() => fullscreen.toggle(readerEl)}
				aria-label={fullscreen.active ? 'Leave fullscreen' : 'Enter fullscreen'}
				title={fullscreen.active ? 'Leave fullscreen (F)' : 'Fullscreen (F)'}
			>
				{#if fullscreen.active}
					<ArrowsIn size={16} weight="bold" />
				{:else}
					<ArrowsOut size={16} weight="bold" />
				{/if}
			</button>
		{/if}

		<button
			class="tool"
			popovertarget="reader-settings"
			aria-label="Reading settings"
			title="Reading settings"
		>
			<Sliders size={16} weight="bold" />
		</button>

		<button
			class="tool danger"
			onclick={() => (confirmOpen = true)}
			aria-label="Delete this chapter"
			title="Delete this chapter"
		>
			<Trash size={16} weight="bold" />
		</button>
	</div>

	<ReaderSettings />

	<main
		style:--read-font={FONT_STACKS[prefs.font]}
		style:--read-size="{prefs.fontSize}px"
		style:--read-leading={prefs.lineHeight}
		style:--read-measure="{prefs.measure}ch"
	>
		{#if loading}
			<div class="loading">
				<div class="skel" style:width="55%" style:height="2rem"></div>
				{#each { length: 12 } as _, i (i)}
					<div class="skel" style:width="{88 - (i % 4) * 9}%"></div>
				{/each}
			</div>
		{:else if loadError || !chapter}
			<div class="failure">
				<h1>Chapter {chapterNumber} is not available</h1>
				<p>{loadError || 'It has not been translated yet.'}</p>
				<a class="action" href="/novel/{novelId}">Back to the novel</a>
			</div>
		{:else if paged}
			{@const leaf = chapter}
			<BookSpread
				html={leaf.content}
				justify={prefs.justify}
				metrics={[prefs.fontSize, prefs.lineHeight, prefs.measure, prefs.font]}
				{spread}
				{spreadCount}
				onchange={(next, count) => {
					spread = next;
					spreadCount = count;
				}}
			>
				{#snippet heading()}
					<h1>{leaf.title}</h1>
					{#if leaf.original_tile}
						<p class="original">{leaf.original_tile}</p>
					{/if}
				{/snippet}
			</BookSpread>

			<nav class="leaf-nav">
				<button
					class="action"
					disabled={spread === 0 && chapterNumber <= 1}
					onclick={turnBack}
					aria-label={spread === 0 ? 'Previous chapter' : 'Previous page'}
				>
					<CaretLeft size={15} weight="bold" />
					<span>{spread === 0 ? 'Previous chapter' : 'Back'}</span>
				</button>

				<p class="leaf-count meta">
					<span class="numeral">{spread + 1}</span> of <span class="numeral">{spreadCount}</span>
				</p>

				{#if spread < spreadCount - 1}
					<button class="action primary" onclick={turnForward} aria-label="Next page">
						<span>Next</span>
						<CaretRight size={15} weight="bold" />
					</button>
				{:else if hasNext}
					<button class="action primary" onclick={() => navigate(1)}>
						<span>Next chapter</span>
						<ArrowRight size={15} weight="bold" />
					</button>
				{:else}
					<button
						class="action primary"
						disabled={translating}
						onclick={() => translateNext({ advance: true })}
					>
						<Sparkle size={15} weight="bold" class={translating ? 'spin' : ''} />
						<span>{translating ? 'Translating…' : 'Translate next'}</span>
					</button>
				{/if}
			</nav>
		{:else}
			<nav class="top-nav">
				<button class="action" disabled={chapterNumber <= 1} onclick={() => navigate(-1)}>
					<ArrowLeft size={14} weight="bold" />
					<span>Previous</span>
				</button>

				<a class="action" href="/novel/{novelId}">
					<ListBullets size={14} weight="bold" />
					<span>Chapters</span>
				</a>

				{#if hasNext}
					<button class="action" onclick={() => navigate(1)}>
						<span>Next</span>
						<ArrowRight size={14} weight="bold" />
					</button>
				{/if}
			</nav>

			<article>
				<h1>{chapter.title}</h1>
				{#if chapter.original_tile}
					<p class="original">{chapter.original_tile}</p>
				{/if}
				<!-- Chapter HTML comes from our own translation backend. -->
				<div class="body prose" class:justified={prefs.justify}>{@html chapter.content}</div>
			</article>

			<nav class="chapter-nav">
				<button class="action" disabled={chapterNumber <= 1} onclick={() => navigate(-1)}>
					<ArrowLeft size={15} weight="bold" />
					Previous
				</button>

				{#if hasNext}
					<button class="action primary" onclick={() => navigate(1)}>
						Next chapter
						<ArrowRight size={15} weight="bold" />
					</button>
				{:else}
					<button
						class="action primary"
						disabled={translating}
						onclick={() => translateNext({ advance: true })}
					>
						<Sparkle size={15} weight="bold" class={translating ? 'spin' : ''} />
						{translating ? 'Translating…' : 'Translate next chapter'}
					</button>
				{/if}
			</nav>
		{/if}
	</main>

	<ConfirmDialog
		open={confirmOpen}
		title="Delete this chapter?"
		description="{chapter?.title ?? 'This chapter'} will be removed. You can translate it again later."
		busy={deleting}
		onconfirm={removeChapter}
		oncancel={() => (confirmOpen = false)}
	/>
</div>

<style>
	.reader {
		min-height: 100dvh;
		background: var(--read-bg);
		color: var(--read-ink);
		padding-bottom: var(--sp-8);
	}

	/* Paged mode owns exactly one viewport: the spread flexes into whatever the
	 * header and the leaf nav leave behind, and nothing scrolls.
	 */
	.reader.paged {
		height: 100dvh;
		min-height: 0;
		display: flex;
		flex-direction: column;
		overflow: hidden;
		padding-bottom: 0;
	}

	.reader.paged main {
		flex: 1;
		min-height: 0;
		display: flex;
		flex-direction: column;
		max-width: calc(2 * var(--read-measure) + 12rem);
		padding-block: var(--sp-5) 0;
	}

	/* --- Progress ---------------------------------------------------------
	 * Driven by the scroll timeline, so the bar costs no JavaScript at all.
	 * Browsers without scroll-timeline simply show a full-width rule.
	 */
	.progress-bar {
		position: fixed;
		inset: 0 auto auto 0;
		z-index: 20;
		height: 2px;
		width: 100%;
		background: var(--read-accent);
		transform-origin: 0 50%;
	}

	@supports (animation-timeline: scroll()) {
		.reader:not(.paged):not(.fullscreen) .progress-bar {
			animation: fill-progress linear;
			animation-timeline: scroll(root block);
		}
	}

	/* Paged mode has no scrolling, and in fullscreen the root scroller is not
	   the one that moves, so both drive the bar from the measured position. */
	.reader.paged .progress-bar,
	.reader.fullscreen .progress-bar {
		width: var(--turned);
		transition: width 320ms var(--ease);
	}

	@keyframes fill-progress {
		from {
			transform: scaleX(0);
		}
		to {
			transform: scaleX(1);
		}
	}

	/* --- Chrome ------------------------------------------------------------
	 * The bar scrolls away with the text. Nothing overlays the column while
	 * reading; the controls that must persist live on the rail instead.
	 */
	header {
		border-bottom: 1px solid var(--read-line);
	}

	.reader.fullscreen header {
		display: none;
	}

	/*
	 * A fullscreen element is sized to the screen, and the document stops
	 * scrolling, so in scroll mode the reader has to become the scroll
	 * container itself or the text is simply clipped.
	 */
	.reader:fullscreen:not(.paged) {
		height: 100%;
		overflow-y: auto;
	}

	.reader:-webkit-full-screen:not(.paged) {
		height: 100%;
		overflow-y: auto;
	}

	.bar {
		display: flex;
		align-items: center;
		gap: var(--sp-4);
		width: 100%;
		max-width: 62rem;
		margin-inline: auto;
		padding: var(--sp-3) var(--sp-5);
		height: 3.5rem;
		font-family: var(--font-ui);
	}

	.back {
		display: flex;
		align-items: center;
		gap: var(--sp-2);
		min-width: 0;
		color: var(--read-ink-2);
		font-size: var(--step--1);
	}

	.back:hover {
		color: var(--read-accent);
	}

	.back-label {
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		max-width: 32ch;
	}

	.stats {
		display: flex;
		align-items: center;
		gap: var(--sp-3);
		margin-left: auto;
		font-size: var(--step--1);
		color: var(--read-ink-2);
	}

	/* --- Side rail ----------------------------------------------------------
	 * Fixed to the right edge and vertically centred: always in reach, never
	 * over the column. It rests at low contrast and comes up on hover or focus.
	 */
	.rail {
		position: fixed;
		right: var(--sp-3);
		top: 50%;
		transform: translateY(-50%);
		z-index: 12;
		display: flex;
		flex-direction: column;
		gap: var(--sp-1);
		padding: var(--sp-1);
		border-radius: var(--r-2);
		background: color-mix(in srgb, var(--read-bg) 80%, transparent);
		backdrop-filter: blur(8px);
		border: 1px solid transparent;
		opacity: 0.45;
		transition:
			opacity var(--dur-2) var(--ease),
			border-color var(--dur-2) var(--ease);
	}

	.rail:hover,
	.rail:focus-within {
		opacity: 1;
		border-color: var(--read-line);
	}

	.tool {
		display: grid;
		place-items: center;
		width: var(--control-h);
		height: var(--control-h);
		background: transparent;
		border: 1px solid transparent;
		border-radius: var(--r-1);
		color: var(--read-ink-2);
		cursor: pointer;
		transition:
			background var(--dur-1) var(--ease),
			color var(--dur-1) var(--ease);
	}

	.tool:hover:not(:disabled) {
		background: color-mix(in srgb, var(--read-accent) 14%, transparent);
		color: var(--read-ink);
	}

	.tool.danger:hover:not(:disabled) {
		color: var(--read-accent);
	}

	/* The prefetch action is the one worth noticing, so it keeps the accent. */
	.prefetch {
		color: var(--read-accent);
	}

	.prefetch:disabled {
		cursor: progress;
	}

	@media (max-width: 640px) {
		main {
			/* The rail shrinks here, so the reserve can too. */
			--rail-gutter: 2.25rem;
			max-width: min(calc(var(--read-measure) + 8rem), calc(100vw - 5rem));
		}

		.rail {
			right: var(--sp-2);
			opacity: 0.85;
		}

		.tool {
			width: var(--control-h-sm);
			height: var(--control-h-sm);
		}

		.back-label {
			max-width: 18ch;
		}
	}


	/* --- Text -------------------------------------------------------------- */
	main {
		/*
		 * The rail is fixed to the right edge, so the text column reserves a
		 * gutter on both sides wide enough that a line can never run underneath
		 * it. Equal padding keeps the column centred.
		 */
		--rail-gutter: 3.5rem;

		width: 100%;
		/*
		 * Capped against the viewport as well as the measure. The reader can set
		 * a line width up to 96ch, which on its own would push the column out
		 * past the rail no matter how much padding it carried.
		 */
		max-width: min(calc(var(--read-measure) + 8rem), calc(100vw - 9rem));
		margin-inline: auto;
		padding: var(--sp-7) var(--rail-gutter) 0;
	}

	article {
		max-width: var(--read-measure);
		margin-inline: auto;
	}

	h1 {
		font-family: var(--font-display);
		font-size: clamp(1.8rem, 4.5vw, 2.6rem);
		font-weight: 600;
		line-height: 1.15;
	}

	.original {
		font-size: var(--step-0);
		color: var(--read-ink-2);
		font-family: var(--font-ui);
	}

	.body {
		margin-top: var(--sp-5);
		font-family: var(--read-font);
		font-size: var(--read-size);
		line-height: var(--read-leading);
		/* Verdana at long measures benefits from a touch of negative tracking. */
		letter-spacing: -0.003em;
	}

	.body.justified {
		text-align: justify;
		hyphens: auto;
	}

	.body :global(p) {
		margin-block: 0 1.15em;
	}

	.body :global(img) {
		max-width: 100%;
		height: auto;
		border-radius: var(--r-1);
	}

	.body :global(a) {
		color: var(--read-accent);
	}

	.body :global(h1),
	.body :global(h2),
	.body :global(h3) {
		font-family: var(--font-display);
		font-weight: 600;
		margin-block: 1.4em 0.5em;
	}

	/* --- Top nav ------------------------------------------------------------ */
	.top-nav {
		display: flex;
		justify-content: center;
		gap: var(--sp-2);
		max-width: var(--read-measure);
		margin: 0 auto var(--sp-5);
		padding-bottom: var(--sp-4);
		border-bottom: 1px solid var(--read-line);
		font-family: var(--font-ui);
	}

	.top-nav .action {
		height: var(--control-h-sm);
		font-size: var(--step--1);
	}

	/* --- Footer nav --------------------------------------------------------- */
	.chapter-nav {
		display: flex;
		justify-content: space-between;
		gap: var(--sp-3);
		max-width: var(--read-measure);
		margin: var(--sp-8) auto 0;
		padding-top: var(--sp-5);
		border-top: 1px solid var(--read-line);
		font-family: var(--font-ui);
	}

	.action {
		display: inline-flex;
		align-items: center;
		gap: var(--sp-2);
		padding: 0 1rem;
		height: var(--control-h);
		line-height: 1;
		background: transparent;
		border: 1px solid var(--read-line);
		border-radius: var(--r-1);
		color: var(--read-ink);
		font-size: var(--step-0);
		font-weight: 500;
		cursor: pointer;
		transition:
			background var(--dur-1) var(--ease),
			transform var(--dur-1) var(--ease);
	}

	.action:hover:not(:disabled) {
		background: color-mix(in srgb, var(--read-accent) 10%, transparent);
	}

	.action:active:not(:disabled) {
		transform: translateY(1px);
	}

	.action:disabled {
		opacity: 0.4;
		cursor: not-allowed;
	}

	.action.primary {
		background: var(--read-accent);
		border-color: var(--read-accent);
		color: var(--read-bg);
	}

	.action.primary:hover:not(:disabled) {
		filter: brightness(1.1);
		background: var(--read-accent);
	}

	/* --- Leaf nav (paged mode) ---------------------------------------------- */
	.leaf-nav {
		flex-shrink: 0;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--sp-3);
		padding-block: var(--sp-4);
		font-family: var(--font-ui);
	}

	.leaf-count {
		color: var(--read-ink-2);
		white-space: nowrap;
	}

	.leaf-nav .action {
		min-width: 9rem;
	}

	@media (max-width: 640px) {
		.leaf-nav .action {
			min-width: 0;
		}

		.leaf-nav .action span {
			display: none;
		}
	}

	/* --- Loading and failure ------------------------------------------------ */
	.loading {
		display: grid;
		gap: var(--sp-4);
		max-width: var(--read-measure);
		margin-inline: auto;
	}

	.skel {
		height: 1rem;
		border-radius: var(--r-1);
		background: color-mix(in srgb, var(--read-ink) 9%, transparent);
		animation: pulse 1.6s ease-in-out infinite;
	}

	@keyframes pulse {
		50% {
			opacity: 0.45;
		}
	}

	.failure {
		display: grid;
		justify-items: center;
		gap: var(--sp-3);
		padding-block: var(--sp-8);
		text-align: center;
		font-family: var(--font-ui);
	}

	.failure h1 {
		font-size: var(--step-2);
	}

	.failure p {
		color: var(--read-ink-2);
	}

	:global(.spin) {
		animation: spin 900ms linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	:global(.reader :focus-visible) {
		outline-color: var(--read-accent);
	}

	@media (max-width: 640px) {
		main {
			padding-top: var(--sp-6);
		}

		.bar {
			padding-inline: var(--sp-4);
		}

		.back-label {
			max-width: 12ch;
		}

		.chapter-nav {
			flex-direction: column-reverse;
		}

		.action {
			justify-content: center;
		}
	}
</style>
