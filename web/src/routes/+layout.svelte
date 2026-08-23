<script lang="ts">
	import '@fontsource-variable/geist';
	import '@fontsource-variable/cormorant-garamond';
	import '../styles/base.css';

	import { page } from '$app/state';
	import Masthead from '$components/Masthead.svelte';
	import Toasts from '$components/Toasts.svelte';
	import ManualExtractDialog from '$components/ManualExtractDialog.svelte';
	import { sources } from '$stores/sources.svelte';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();

	// The reader owns the whole viewport and carries its own theme and chrome.
	const isReader = $derived(/^\/novel\/[^/]+\/chapter\//.test(page.url.pathname));

	// Fired once on mount rather than from an effect, which would re-run when
	// the store's own loaded flag flips.
	sources.load();
</script>

{#if isReader}
	{@render children()}
{:else}
	<div class="shell">
		<Masthead />
		<main id="main" tabindex="-1">{@render children()}</main>

		<footer>
			<div class="page inner">
				<p>Chinese, Korean and Japanese webnovels, translated as you read.</p>
				<nav aria-label="Footer">
					<a href="/library">Library</a>
					<a href="/add">Add novel</a>
				</nav>
			</div>
		</footer>
	</div>
{/if}

<Toasts />
<ManualExtractDialog />

<style>
	.shell {
		display: flex;
		flex-direction: column;
		min-height: 100dvh;
	}

	main {
		flex: 1;
		padding-block: var(--sp-6) var(--sp-8);
	}

	/* The skip link targets main, but it should never draw a ring on click. */
	main:focus {
		outline: none;
	}

	footer {
		border-top: 1px solid var(--line);
		padding-block: var(--sp-5);
	}

	.inner {
		display: flex;
		flex-wrap: wrap;
		gap: var(--sp-4);
		align-items: center;
		justify-content: space-between;
		font-size: var(--step--1);
		color: var(--ink-3);
	}

	footer nav {
		display: flex;
		gap: var(--sp-4);
	}

	footer a:hover {
		color: var(--accent);
	}
</style>
