<script lang="ts">
	import { goto } from '$app/navigation';
	import { browser } from '$app/environment';
	import { searchNovels } from '$api/novels';
	import Cover from './Cover.svelte';
	import { ArrowRight, ClockCounterClockwise, MagnifyingGlass, Plus } from 'phosphor-svelte';
	import type { Novel } from '$lib/types';

	interface Props {
		open: boolean;
		onclose: () => void;
	}

	let { open, onclose }: Props = $props();

	const RECENT_KEY = 'arcane:recent-searches';
	const MAX_RECENT = 6;

	let query = $state('');
	let results = $state<Novel[]>([]);
	let searching = $state(false);
	let recent = $state<string[]>(loadRecent());
	let input = $state<HTMLInputElement>();
	let dialog = $state<HTMLDialogElement>();

	function loadRecent(): string[] {
		if (!browser) return [];
		try {
			const parsed = JSON.parse(localStorage.getItem(RECENT_KEY) ?? '[]');
			return Array.isArray(parsed) ? (parsed as string[]) : [];
		} catch {
			return [];
		}
	}

	function remember(term: string) {
		recent = [term, ...recent.filter((item) => item !== term)].slice(0, MAX_RECENT);
		if (browser) localStorage.setItem(RECENT_KEY, JSON.stringify(recent));
	}

	function clearRecent() {
		recent = [];
		if (browser) localStorage.removeItem(RECENT_KEY);
	}

	$effect(() => {
		if (!dialog) return;
		if (open && !dialog.open) {
			dialog.showModal();
			query = '';
			results = [];
			queueMicrotask(() => input?.focus());
		}
		if (!open && dialog.open) dialog.close();
	});

	// Live preview of matches while typing, debounced so each keystroke does not
	// hit the API.
	$effect(() => {
		const term = query.trim();
		if (!open || term.length < 2) {
			results = [];
			searching = false;
			return;
		}

		searching = true;
		let cancelled = false;
		const timer = setTimeout(async () => {
			try {
				const found = await searchNovels(term);
				if (!cancelled) results = found.slice(0, 6);
			} catch {
				if (!cancelled) results = [];
			} finally {
				if (!cancelled) searching = false;
			}
		}, 220);

		return () => {
			cancelled = true;
			clearTimeout(timer);
		};
	});

	function submit(term = query) {
		const trimmed = term.trim();
		if (!trimmed) return;
		remember(trimmed);
		onclose();
		goto(`/search?q=${encodeURIComponent(trimmed)}`);
	}

	function jump(href: string) {
		onclose();
		goto(href);
	}
</script>

<dialog
	bind:this={dialog}
	oncancel={(event) => {
		event.preventDefault();
		onclose();
	}}
	onclick={(event) => {
		if (event.target === dialog) onclose();
	}}
>
	<div class="palette">
		<form
			class="bar"
			onsubmit={(event) => {
				event.preventDefault();
				submit();
			}}
		>
			<MagnifyingGlass size={17} />
			<input
				bind:this={input}
				bind:value={query}
				type="search"
				placeholder="Search your library"
				autocomplete="off"
				spellcheck="false"
				aria-label="Search novels"
			/>
			<kbd>esc</kbd>
		</form>

		<div class="results">
			{#if query.trim().length >= 2}
				<p class="group-label">
					{searching ? 'Searching…' : `${results.length} match${results.length === 1 ? '' : 'es'}`}
				</p>
				{#each results as novel (novel.id)}
					<button class="row novel" onclick={() => jump(`/novel/${novel.id}`)}>
						<span class="thumb"><Cover {novel} eager /></span>
						<span class="labels">
							<span class="label">{novel.title}</span>
							{#if novel.author}<span class="sub">{novel.author}</span>{/if}
						</span>
					</button>
				{/each}
				<button class="row" onclick={() => submit()}>
					<MagnifyingGlass size={15} />
					<span class="label">See all results for &ldquo;{query.trim()}&rdquo;</span>
					<ArrowRight size={14} />
				</button>
			{:else}
				{#if recent.length}
					<div class="group-head">
						<p class="group-label">Recent</p>
						<button class="link" onclick={clearRecent}>Clear</button>
					</div>
					{#each recent as term (term)}
						<button class="row" onclick={() => submit(term)}>
							<ClockCounterClockwise size={15} />
							<span class="label">{term}</span>
						</button>
					{/each}
				{/if}

				<p class="group-label">Go to</p>
				<button class="row" onclick={() => jump('/library')}>
					<ArrowRight size={15} />
					<span class="label">Library</span>
				</button>
				<button class="row" onclick={() => jump('/add')}>
					<Plus size={15} />
					<span class="label">Add a novel</span>
				</button>
			{/if}
		</div>
	</div>
</dialog>

<style>
	dialog {
		/* Explicit, because the global reset zeroes the UA's auto margins. */
		margin: 12vh auto auto;
		padding: 0;
		border: 0;
		background: transparent;
		width: min(34rem, calc(100vw - 2rem));
		color: var(--ink);
	}

	dialog::backdrop {
		background: #0b0d0eab;
		backdrop-filter: blur(3px);
	}

	.palette {
		background: var(--surface);
		border: 1px solid var(--line-strong);
		border-radius: var(--r-2);
		box-shadow: var(--shadow-overlay);
		overflow: hidden;
	}

	.bar {
		display: flex;
		align-items: center;
		gap: var(--sp-3);
		padding: var(--sp-4);
		border-bottom: 1px solid var(--line);
		color: var(--ink-3);
	}

	/* The input clears its own outline, so the row carries the focus ring. */
	.bar:focus-within {
		box-shadow: inset 2px 0 0 var(--accent);
		color: var(--ink-2);
	}

	.bar input {
		flex: 1;
		min-width: 0;
		border: 0;
		background: none;
		outline: none;
		font-size: var(--step-1);
		color: var(--ink);
	}

	.bar input::placeholder {
		color: var(--ink-3);
	}

	.bar input::-webkit-search-cancel-button {
		display: none;
	}

	kbd {
		font-family: var(--font-mono);
		font-size: 0.65rem;
		text-transform: uppercase;
		padding: 0.15rem 0.35rem;
		border: 1px solid var(--line-strong);
		border-radius: var(--r-1);
		color: var(--ink-3);
	}

	.results {
		max-height: min(26rem, 60vh);
		overflow-y: auto;
		overscroll-behavior: contain;
		padding: var(--sp-2);
	}

	.group-head {
		display: flex;
		align-items: baseline;
		justify-content: space-between;
	}

	.group-label {
		padding: var(--sp-2) var(--sp-3) var(--sp-1);
		font-size: var(--step--1);
		color: var(--ink-3);
	}

	.link {
		background: none;
		border: 0;
		padding: var(--sp-2) var(--sp-3);
		font-size: var(--step--1);
		color: var(--accent);
		cursor: pointer;
	}

	.row {
		display: flex;
		align-items: center;
		gap: var(--sp-3);
		width: 100%;
		padding: var(--sp-2) var(--sp-3);
		background: none;
		border: 0;
		border-radius: var(--r-1);
		text-align: left;
		color: var(--ink-2);
		cursor: pointer;
	}

	.row:hover,
	.row:focus-visible {
		background: var(--accent-wash);
		color: var(--ink);
	}

	.row .label {
		flex: 1;
		min-width: 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		color: var(--ink);
	}

	.novel {
		padding-block: var(--sp-2);
	}

	.thumb {
		width: 1.75rem;
		flex-shrink: 0;
	}

	.labels {
		flex: 1;
		min-width: 0;
		display: grid;
	}

	.sub {
		font-size: var(--step--1);
		color: var(--ink-3);
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
</style>
