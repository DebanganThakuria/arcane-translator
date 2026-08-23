<script lang="ts">
	import { goto } from '$app/navigation';
	import { extractNovel } from '$api/translation';
	import { sources } from '$stores/sources.svelte';
	import { toast } from '$stores/toast.svelte';
	import { ManualExtractCancelled } from '$stores/manualExtract.svelte';
	import { Spinner } from 'phosphor-svelte';

	let url = $state('');
	let sourceId = $state('');
	let submitting = $state(false);
	let error = $state('');

	const grouped = $derived.by(() => {
		const groups = new Map<string, typeof sources.all>();
		for (const site of sources.all) {
			const bucket = groups.get(site.language) ?? [];
			bucket.push(site);
			groups.set(site.language, bucket);
		}
		return [...groups.entries()].sort(([a], [b]) => a.localeCompare(b));
	});

	async function submit(event: SubmitEvent) {
		event.preventDefault();
		error = '';

		if (!url.trim() || !sourceId) {
			error = 'Enter the novel URL and pick the source site it came from.';
			return;
		}

		submitting = true;
		try {
			const novel = await extractNovel(url.trim(), sourceId);
			toast.success('Novel added', novel.title);
			await goto(`/novel/${novel.id}`);
		} catch (cause) {
			if (cause instanceof ManualExtractCancelled) {
				error = 'Cancelled. Nothing was added.';
			} else {
				error = cause instanceof Error ? cause.message : 'Could not add that novel.';
			}
		} finally {
			submitting = false;
		}
	}
</script>

<svelte:head><title>Add a novel · Arcane Translator</title></svelte:head>

<div class="page wrap">
	<header>
		<h1 class="display">Add a novel</h1>
		<p>
			Paste the link to a novel's main page. The chapters are translated one at a time as you read
			them, so only the details are fetched now.
		</p>
	</header>

	<form onsubmit={submit}>
		<div class="field">
			<label for="novel-url">Novel URL</label>
			<input
				id="novel-url"
				class="input"
				type="url"
				inputmode="url"
				autocomplete="off"
				spellcheck="false"
				placeholder="https://www.69shuba.com/book/84522.htm"
				bind:value={url}
				disabled={submitting}
			/>
			<p class="hint">Use the novel's index page, not a chapter page.</p>
		</div>

		<div class="field">
			<label for="novel-source">Source site</label>
			<select id="novel-source" class="select" bind:value={sourceId} disabled={submitting || !sources.loaded}>
				<option value="" disabled>
					{sources.loaded ? 'Choose a source' : 'Loading sources…'}
				</option>
				{#each grouped as [language, sites] (language)}
					<optgroup label={language}>
						{#each sites as site (site.id)}
							<option value={site.id}>{site.name}</option>
						{/each}
					</optgroup>
				{/each}
			</select>
			{#if sources.loaded && sources.all.length === 0}
				<p class="error">
					No sources came back from the server. Check that the backend is running on port 8088.
				</p>
			{/if}
		</div>

		{#if error}
			<p class="form-error" role="alert">{error}</p>
		{/if}

		<button class="btn btn--primary btn--block" type="submit" disabled={submitting}>
			{#if submitting}
				<Spinner size={15} weight="bold" class="spin" />
				Fetching details…
			{:else}
				Add novel
			{/if}
		</button>
	</form>

	<p class="aside">
		If the source blocks the server, you will be asked to open the page yourself and paste its source.
	</p>
</div>

<style>
	.wrap {
		max-width: 34rem;
	}

	header {
		margin-bottom: var(--sp-6);
	}

	h1 {
		font-size: var(--step-3);
		margin-bottom: var(--sp-2);
	}

	header p {
		color: var(--ink-2);
		text-wrap: pretty;
	}

	form {
		display: grid;
		gap: var(--sp-5);
		padding: var(--sp-5);
		background: var(--surface);
		border: 1px solid var(--line);
		border-radius: var(--r-2);
	}

	.form-error {
		padding: var(--sp-3);
		border-left: 2px solid var(--danger);
		background: var(--surface-2);
		border-radius: var(--r-1);
		font-size: var(--step--1);
		color: var(--danger);
	}

	.aside {
		margin-top: var(--sp-4);
		font-size: var(--step--1);
		color: var(--ink-3);
		text-wrap: pretty;
	}

	:global(.spin) {
		animation: spin 900ms linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
