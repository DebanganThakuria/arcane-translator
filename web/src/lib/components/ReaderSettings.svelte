<script lang="ts">
	import {
		reader,
		FONT_LABELS,
		FONT_STACKS,
		LIMITS,
		type ReaderFont,
		type ReaderMode,
		type ReaderTheme
	} from '$stores/reader.svelte';

	const THEMES: { key: ReaderTheme; label: string }[] = [
		{ key: 'sepia', label: 'Sepia' },
		{ key: 'light', label: 'Light' },
		{ key: 'dark', label: 'Dark' }
	];

	const MODES: { key: ReaderMode; label: string }[] = [
		{ key: 'scroll', label: 'Scroll' },
		{ key: 'paged', label: 'Two pages' }
	];

	const FONTS = Object.keys(FONT_LABELS) as ReaderFont[];

	const prefs = $derived(reader.prefs);
</script>

<div class="sheet" popover="auto" id="reader-settings">
	<section>
		<h3>Theme</h3>
		<div class="swatches">
			{#each THEMES as option (option.key)}
				<button
					data-reader={option.key}
					class:selected={prefs.theme === option.key}
					onclick={() => reader.set('theme', option.key)}
					aria-pressed={prefs.theme === option.key}
				>
					<span class="sample">Aa</span>
					<span class="name">{option.label}</span>
				</button>
			{/each}
		</div>
	</section>

	<section>
		<h3>Layout</h3>
		<div class="modes" role="group" aria-label="Reading layout">
			{#each MODES as option (option.key)}
				<button
					class="mode"
					class:selected={prefs.mode === option.key}
					aria-pressed={prefs.mode === option.key}
					onclick={() => reader.set('mode', option.key)}
				>
					<span class="glyph" aria-hidden="true">
						{#if option.key === 'scroll'}
							<span class="col"></span>
						{:else}
							<span class="col"></span><span class="col"></span>
						{/if}
					</span>
					{option.label}
				</button>
			{/each}
		</div>
	</section>

	<section>
		<h3>Typeface</h3>
		<div class="fonts">
			{#each FONTS as font (font)}
				<button
					class="font"
					class:selected={prefs.font === font}
					style:font-family={FONT_STACKS[font]}
					onclick={() => reader.set('font', font)}
					aria-pressed={prefs.font === font}
				>
					{FONT_LABELS[font]}
				</button>
			{/each}
		</div>
	</section>

	<section>
		<label class="slider" for="reader-size">
			<span>Size</span>
			<output class="numeral">{prefs.fontSize}px</output>
		</label>
		<input
			id="reader-size"
			type="range"
			min={LIMITS.fontSize.min}
			max={LIMITS.fontSize.max}
			step={LIMITS.fontSize.step}
			value={prefs.fontSize}
			oninput={(event) => reader.set('fontSize', Number(event.currentTarget.value))}
		/>

		<label class="slider" for="reader-leading">
			<span>Line height</span>
			<output class="numeral">{prefs.lineHeight.toFixed(1)}</output>
		</label>
		<input
			id="reader-leading"
			type="range"
			min={LIMITS.lineHeight.min}
			max={LIMITS.lineHeight.max}
			step={LIMITS.lineHeight.step}
			value={prefs.lineHeight}
			oninput={(event) => reader.set('lineHeight', Number(event.currentTarget.value))}
		/>

		<label class="slider" for="reader-measure">
			<span>Line width</span>
			<output class="numeral">{prefs.measure} ch</output>
		</label>
		<input
			id="reader-measure"
			type="range"
			min={LIMITS.measure.min}
			max={LIMITS.measure.max}
			step={LIMITS.measure.step}
			value={prefs.measure}
			oninput={(event) => reader.set('measure', Number(event.currentTarget.value))}
		/>
	</section>

	<section class="row">
		<label class="toggle" for="reader-justify">
			<input
				id="reader-justify"
				type="checkbox"
				checked={prefs.justify}
				onchange={(event) => reader.set('justify', event.currentTarget.checked)}
			/>
			<span>Justify text</span>
		</label>
		<button class="reset" onclick={() => reader.reset()}>Reset</button>
	</section>
</div>

<style>
	.sheet {
		position: fixed;
		inset: auto var(--sp-4) auto auto;
		top: 4.25rem;
		margin: 0;
		width: min(19rem, calc(100vw - 2rem));
		display: grid;
		gap: var(--sp-5);
		padding: var(--sp-5);
		background: var(--read-bg-2);
		color: var(--read-ink);
		border: 1px solid var(--read-line);
		border-radius: var(--r-2);
		box-shadow: 0 1px 2px #0000001f, 0 16px 48px -8px #00000038;
		font-family: var(--font-ui);
		font-size: var(--step-0);
	}

	.sheet:not(:popover-open) {
		display: none;
	}

	h3 {
		font-size: var(--step--1);
		font-weight: 600;
		color: var(--read-ink-2);
		margin-bottom: var(--sp-3);
	}

	/* --- Theme swatches --------------------------------------------------- */
	.swatches {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: var(--sp-2);
	}

	.swatches button {
		display: grid;
		gap: var(--sp-2);
		padding: var(--sp-2);
		background: var(--read-bg);
		border: 1px solid var(--read-line);
		border-radius: var(--r-1);
		cursor: pointer;
		color: var(--read-ink-2);
	}

	.sample {
		display: grid;
		place-items: center;
		height: 2.5rem;
		border-radius: var(--r-1);
		background: var(--read-bg);
		color: var(--read-ink);
		font-family: var(--font-display);
		font-size: 1.1rem;
	}

	.name {
		font-size: 0.7rem;
		text-align: center;
	}

	.swatches button.selected {
		border-color: var(--read-accent);
		box-shadow: inset 0 0 0 1px var(--read-accent);
		color: var(--read-ink);
	}

	/* --- Layout mode ------------------------------------------------------ */
	.modes {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: var(--sp-2);
	}

	.mode {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: var(--sp-2);
		padding: var(--sp-2);
		background: var(--read-bg);
		border: 1px solid var(--read-line);
		border-radius: var(--r-1);
		color: var(--read-ink-2);
		font-size: var(--step--1);
		cursor: pointer;
	}

	.mode.selected {
		border-color: var(--read-accent);
		box-shadow: inset 0 0 0 1px var(--read-accent);
		color: var(--read-ink);
	}

	/* A miniature of the layout it selects: one column, or two. */
	.glyph {
		display: flex;
		gap: 2px;
		width: 1rem;
		height: 0.85rem;
	}

	.glyph .col {
		flex: 1;
		border: 1px solid currentColor;
		border-radius: 1px;
		opacity: 0.7;
	}

	/* --- Fonts ------------------------------------------------------------ */
	.fonts {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: var(--sp-2);
	}

	.font {
		padding: var(--sp-2);
		background: var(--read-bg);
		border: 1px solid var(--read-line);
		border-radius: var(--r-1);
		color: var(--read-ink-2);
		font-size: var(--step--1);
		cursor: pointer;
	}

	.font.selected {
		border-color: var(--read-accent);
		box-shadow: inset 0 0 0 1px var(--read-accent);
		color: var(--read-ink);
	}

	/* --- Sliders ---------------------------------------------------------- */
	.slider {
		display: flex;
		justify-content: space-between;
		align-items: baseline;
		font-size: var(--step--1);
		color: var(--read-ink-2);
		margin-bottom: var(--sp-1);
	}

	.slider + input {
		margin-bottom: var(--sp-4);
	}

	input[type='range'] {
		width: 100%;
		accent-color: var(--read-accent);
	}

	section input[type='range']:last-child {
		margin-bottom: 0;
	}

	/* --- Footer ----------------------------------------------------------- */
	.row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding-top: var(--sp-4);
		border-top: 1px solid var(--read-line);
	}

	.toggle {
		display: flex;
		align-items: center;
		gap: var(--sp-2);
		font-size: var(--step--1);
		color: var(--read-ink-2);
		cursor: pointer;
	}

	.toggle input {
		accent-color: var(--read-accent);
	}

	.reset {
		background: none;
		border: 0;
		padding: 0;
		font-size: var(--step--1);
		color: var(--read-accent);
		cursor: pointer;
	}

	button:focus-visible,
	input:focus-visible {
		outline: 2px solid var(--read-accent);
		outline-offset: 2px;
	}
</style>
