<script lang="ts">
	import { page } from '$app/state';
	import SearchPalette from './SearchPalette.svelte';
	import { theme } from '$stores/theme.svelte';
	import { MagnifyingGlass, Moon, Plus, Sun } from 'phosphor-svelte';

	// Home is reached through the wordmark, so it is not repeated here.
	const LINKS = [
		{ href: '/library', label: 'Library' },
		{ href: '/chinese', label: 'Chinese' },
		{ href: '/korean', label: 'Korean' },
		{ href: '/japanese', label: 'Japanese' }
	];

	let searchOpen = $state(false);

	function isCurrent(href: string): boolean {
		return page.url.pathname.startsWith(href);
	}
</script>

<svelte:window
	onkeydown={(event) => {
		if (event.key === 'k' && (event.metaKey || event.ctrlKey)) {
			event.preventDefault();
			searchOpen = !searchOpen;
		}
	}}
/>

<a class="skip" href="#main">Skip to content</a>

<header>
	<div class="page bar">
		<a class="mark" href="/" aria-label="Arcane Translator, home">
			<span class="sigil" aria-hidden="true">A</span>
			<span class="wordmark">Arcane Translator</span>
		</a>

		<nav aria-label="Sections">
			{#each LINKS as link (link.href)}
				{@const current = isCurrent(link.href)}
				<a href={link.href} class:current aria-current={current ? 'page' : undefined}>
					{link.label}
				</a>
			{/each}
		</nav>

		<div class="tools">
			<button class="search" onclick={() => (searchOpen = true)} aria-label="Search novels">
				<MagnifyingGlass size={15} />
				<span class="search-label" aria-hidden="true">Search</span>
				<kbd aria-hidden="true">&#8984;K</kbd>
			</button>

			<button
				class="btn btn--icon"
				onclick={() => theme.toggle()}
				aria-label="Switch to {theme.value === 'dark' ? 'light' : 'dark'} theme"
			>
				{#if theme.value === 'dark'}
					<Sun size={15} weight="bold" />
				{:else}
					<Moon size={15} weight="bold" />
				{/if}
			</button>

			<a class="btn btn--primary add" href="/add">
				<Plus size={14} weight="bold" />
				<span>Add novel</span>
			</a>
		</div>
	</div>
</header>

<SearchPalette open={searchOpen} onclose={() => (searchOpen = false)} />

<style>
	.skip {
		position: absolute;
		left: var(--sp-4);
		top: -3rem;
		z-index: var(--z-dialog);
		padding: var(--sp-2) var(--sp-4);
		background: var(--accent);
		color: var(--accent-ink);
		border-radius: var(--r-1);
		font-size: var(--step--1);
		font-weight: 600;
		transition: top var(--dur-1) var(--ease);
	}

	.skip:focus-visible {
		top: var(--sp-3);
	}

	header {
		position: sticky;
		top: 0;
		z-index: var(--z-sticky);
		background: color-mix(in srgb, var(--ground) 88%, transparent);
		backdrop-filter: blur(12px);
		border-bottom: 1px solid var(--line);
	}

	.bar {
		display: flex;
		align-items: center;
		gap: var(--sp-5);
		height: 68px;
	}

	.mark {
		display: flex;
		align-items: center;
		gap: var(--sp-3);
		flex-shrink: 0;
	}

	/* The seal is the wordmark's anchor: one letter, cinnabar, like a stamp. */
	.sigil {
		display: grid;
		place-items: center;
		width: 1.75rem;
		height: 1.75rem;
		border-radius: var(--r-1);
		background: var(--accent);
		color: var(--accent-ink);
		font-family: var(--font-display);
		font-size: 1.15rem;
		line-height: 1;
		font-weight: 600;
	}

	.wordmark {
		font-family: var(--font-display);
		font-size: 1.3rem;
		font-weight: 600;
		letter-spacing: 0.01em;
		white-space: nowrap;
	}

	nav {
		display: flex;
		gap: var(--sp-1);
		min-width: 0;
		overflow-x: auto;
		scrollbar-width: none;
	}

	nav::-webkit-scrollbar {
		display: none;
	}

	nav a {
		position: relative;
		padding: var(--sp-2) var(--sp-3);
		border-radius: var(--r-1);
		font-size: var(--step-0);
		color: var(--ink-2);
		white-space: nowrap;
		transition: color var(--dur-1) var(--ease);
	}

	nav a:hover {
		color: var(--ink);
	}

	nav a.current {
		color: var(--ink);
		font-weight: 600;
	}

	nav a.current::after {
		content: '';
		position: absolute;
		left: var(--sp-3);
		right: var(--sp-3);
		bottom: -0.35rem;
		height: 2px;
		background: var(--accent);
	}

	.tools {
		display: flex;
		align-items: center;
		gap: var(--sp-2);
		margin-left: auto;
	}

	.search {
		display: flex;
		align-items: center;
		gap: var(--sp-2);
		height: var(--control-h);
		padding-inline: var(--sp-3);
		background: var(--surface);
		border: 1px solid var(--line-strong);
		border-radius: var(--r-1);
		color: var(--ink-3);
		cursor: pointer;
		transition: border-color var(--dur-1) var(--ease);
	}

	.search:hover {
		border-color: var(--ink-3);
		color: var(--ink-2);
	}

	.search-label {
		font-size: var(--step--1);
	}

	kbd {
		font-family: var(--font-mono);
		font-size: 0.65rem;
		color: var(--ink-3);
	}

	/* --- Narrow viewports -------------------------------------------------
	 * The bar stays one line at every width. Labels drop before the nav does,
	 * and the nav itself becomes a horizontal scroller rather than wrapping.
	 */
	@media (max-width: 1080px) {
		.bar {
			gap: var(--sp-4);
		}

		.search-label,
		kbd {
			display: none;
		}

		.search {
			width: var(--control-h);
			justify-content: center;
			padding-inline: 0;
		}
	}

	@media (max-width: 860px) {
		.wordmark {
			display: none;
		}
	}

	@media (max-width: 720px) {
		.add span {
			display: none;
		}

		.add {
			width: var(--control-h);
			padding-inline: 0;
		}
	}
</style>
