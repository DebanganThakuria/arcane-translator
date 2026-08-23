<script lang="ts">
	import Dialog from './Dialog.svelte';
	import { manualExtract, MIN_HTML_LENGTH } from '$stores/manualExtract.svelte';

	let html = $state('');

	const pending = $derived(manualExtract.pending);
	const length = $derived(html.trim().length);
	const ready = $derived(length >= MIN_HTML_LENGTH);

	// Each new request starts from a clean textarea.
	$effect(() => {
		if (pending) html = '';
	});

	function submit() {
		if (ready) manualExtract.submit(html.trim());
	}
</script>

<Dialog
	open={Boolean(pending)}
	wide
	title="Paste the page source"
	description="This source blocked the server, so the page was opened in a new tab. Copy its source and paste it below."
	onclose={() => manualExtract.cancel()}
>
	<ol class="steps">
		<li>Switch to the tab that just opened.</li>
		<li>View the page source, then select all and copy.</li>
		<li>Paste it here and continue.</li>
	</ol>

	{#if pending}
		<p class="url" title={pending.url}>{pending.url}</p>
	{/if}

	<div class="field">
		<label for="page-source">Page source</label>
		<!-- svelte-ignore a11y_autofocus -->
		<textarea
			id="page-source"
			class="textarea"
			autofocus
			bind:value={html}
			placeholder="Paste the full HTML here"
			onkeydown={(event) => {
				if (event.key === 'Enter' && (event.metaKey || event.ctrlKey)) submit();
			}}
		></textarea>
		<p class="hint">
			<span class="numeral">{length.toLocaleString('en-US')}</span> characters pasted, {MIN_HTML_LENGTH}
			minimum.
		</p>
	</div>

	{#snippet footer()}
		<button class="btn" onclick={() => manualExtract.cancel()}>Cancel</button>
		<button class="btn btn--primary" disabled={!ready} onclick={submit}>Continue</button>
	{/snippet}
</Dialog>

<style>
	.steps {
		display: grid;
		gap: var(--sp-2);
		counter-reset: step;
		font-size: var(--step--1);
		color: var(--ink-2);
	}

	.steps li {
		display: grid;
		grid-template-columns: 1.25rem 1fr;
		gap: var(--sp-2);
		counter-increment: step;
	}

	.steps li::before {
		content: counter(step);
		font-family: var(--font-mono);
		font-size: 0.7rem;
		color: var(--accent);
		padding-top: 0.15em;
	}

	.url {
		font-family: var(--font-mono);
		font-size: var(--step--1);
		color: var(--ink-2);
		background: var(--surface-sunk);
		border-radius: var(--r-1);
		padding: var(--sp-2) var(--sp-3);
		overflow-wrap: anywhere;
	}
</style>
