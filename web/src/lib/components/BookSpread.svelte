<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		/** Translated chapter HTML from our own backend. */
		html: string;
		/** Rendered at the top of the first leaf so it flows with the text. */
		heading: Snippet;
		justify: boolean;
		/** Re-measure whenever any of these change the text metrics. */
		metrics: unknown;
		spread: number;
		spreadCount: number;
		onchange: (spread: number, spreadCount: number) => void;
	}

	let { html, heading, justify, metrics, spread = 0, spreadCount = 1, onchange }: Props = $props();

	let frame = $state<HTMLDivElement>();
	let leaves = $state<HTMLDivElement>();

	/**
	 * The text is laid out with CSS multi-column at the frame's exact height, so
	 * columns three and beyond become overflow columns to the right. The browser
	 * does all the line breaking; turning a page only moves the scroll offset.
	 *
	 * Scroll rather than transform: Chrome paints multicol overflow columns from
	 * the scroll origin, so translating the container moves the box without
	 * revealing the later columns.
	 */
	function step(): number {
		if (!leaves) return 0;
		const gutter = parseFloat(getComputedStyle(leaves).columnGap) || 0;
		return leaves.clientWidth + gutter;
	}

	function measure() {
		if (!leaves) return;

		const width = step();
		if (width <= 0) return;

		const gutter = width - leaves.clientWidth;
		// scrollWidth omits the trailing gutter, so add it back before dividing.
		const measured = Math.max(1, Math.round((leaves.scrollWidth + gutter) / width));
		const clamped = Math.min(spread, measured - 1);

		if (measured !== spreadCount || clamped !== spread) onchange(clamped, measured);
	}

	// Keep the scroll offset in step with the requested spread.
	$effect(() => {
		if (!leaves) return;
		const width = step();
		if (width > 0) leaves.scrollLeft = spread * width;
	});

	$effect(() => {
		// Content and every typographic control that changes where lines break.
		html;
		justify;
		metrics;

		if (!frame) return;

		// A webfont landing after first paint re-flows the columns.
		let cancelled = false;
		document.fonts?.ready.then(() => {
			if (!cancelled) measure();
		});

		const observer = new ResizeObserver(measure);
		observer.observe(frame);

		// Two frames: one for the style change, one for the column re-layout.
		const raf = requestAnimationFrame(() => requestAnimationFrame(measure));

		return () => {
			cancelled = true;
			observer.disconnect();
			cancelAnimationFrame(raf);
		};
	});
</script>

<div class="frame" bind:this={frame}>
	<div class="leaves" class:justify bind:this={leaves}>
		{@render heading()}
		<!-- eslint-disable-next-line svelte/no-at-html-tags -->
		<div class="prose">{@html html}</div>
	</div>

	<div class="gutter" aria-hidden="true"></div>
</div>

<style>
	.frame {
		position: relative;
		height: 100%;
		overflow: hidden;
	}

	.leaves {
		height: 100%;
		columns: 2;
		column-gap: var(--sp-8);
		column-fill: auto;
		font-family: var(--read-font);
		font-size: var(--read-size);
		line-height: var(--read-leading);
		letter-spacing: -0.003em;
		/*
		 * hidden, not clip: this element must stay a scroll container so the
		 * page turn can move scrollLeft. The scrollbar is never shown.
		 */
		overflow: hidden;
		scroll-behavior: smooth;
	}

	@media (prefers-reduced-motion: reduce) {
		.leaves {
			scroll-behavior: auto;
		}
	}

	.leaves.justify {
		text-align: justify;
		hyphens: auto;
	}

	/* The centre rule reads as the fold between two leaves. */
	.gutter {
		position: absolute;
		top: 0;
		bottom: 0;
		left: 50%;
		width: 1px;
		background: var(--read-line);
		opacity: 0.7;
		pointer-events: none;
	}

	.leaves :global(p) {
		margin-block: 0 0.9em;
		/* Never strand a single line of a paragraph across the fold. */
		orphans: 2;
		widows: 2;
	}

	.leaves :global(h1),
	.leaves :global(h2),
	.leaves :global(h3) {
		font-family: var(--font-display);
		font-weight: 600;
		margin-block: 1.2em 0.5em;
		break-after: avoid;
	}

	.leaves :global(img) {
		max-width: 100%;
		height: auto;
		border-radius: var(--r-1);
	}

	.leaves :global(a) {
		color: var(--read-accent);
	}

	/* One column on narrow screens: a two-page spread needs the width. */
	@media (max-width: 820px) {
		.leaves {
			columns: 1;
		}

		.gutter {
			display: none;
		}
	}
</style>
