<script lang="ts">
	import { fallbackCover } from '$utils/cover';
	import type { Novel } from '$lib/types';

	interface Props {
		novel: Pick<Novel, 'id' | 'title' | 'cover'>;
		/** Covers above the fold on the reader's entry point should load eagerly. */
		eager?: boolean;
	}

	let { novel, eager = false }: Props = $props();

	let broken = $state(false);
	let loaded = $state(false);

	const fallback = $derived(fallbackCover(novel));
	const showImage = $derived(Boolean(novel.cover) && !broken);

	// A new novel in a recycled component must re-arm the image state.
	$effect(() => {
		novel.id;
		broken = false;
		loaded = false;
	});
</script>

<div class="cover" style:background={showImage && loaded ? 'var(--surface-sunk)' : fallback.background}>
	{#if showImage}
		<img
			src={novel.cover}
			alt=""
			loading={eager ? 'eager' : 'lazy'}
			decoding="async"
			class:loaded
			onload={() => (loaded = true)}
			onerror={() => (broken = true)}
		/>
	{/if}

	{#if !showImage}
		<span class="glyph" style:color={fallback.ink}>{fallback.glyph}</span>
	{/if}
</div>

<style>
	.cover {
		position: relative;
		container-type: inline-size;
		aspect-ratio: 2 / 3;
		width: 100%;
		overflow: hidden;
		border-radius: var(--r-2);
		border: 1px solid var(--line);
		/* Inner top highlight reads as a page edge catching light. */
		box-shadow: inset 0 1px 0 #ffffff1f;
	}

	img {
		width: 100%;
		height: 100%;
		object-fit: cover;
		opacity: 0;
		transition: opacity var(--dur-2) var(--ease);
	}

	img.loaded {
		opacity: 1;
	}

	.glyph {
		position: absolute;
		inset: 0;
		display: grid;
		place-items: center;
		font-family: var(--font-display);
		font-size: clamp(1.75rem, 26cqw, 3.5rem);
		font-weight: 600;
		letter-spacing: 0.02em;
		opacity: 0.9;
		user-select: none;
	}
</style>
