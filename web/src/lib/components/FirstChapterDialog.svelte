<script lang="ts">
	import Dialog from './Dialog.svelte';
	import { setFirstChapter } from '$api/translation';
	import { ManualExtractCancelled } from '$stores/manualExtract.svelte';

	interface Props {
		open: boolean;
		novelId: string;
		novelUrl: string;
		onclose: () => void;
		onsuccess: (chapterNumber: number) => void;
	}

	let { open, novelId, novelUrl, onclose, onsuccess }: Props = $props();

	let url = $state('');
	let submitting = $state(false);
	let error = $state('');

	const host = $derived.by(() => {
		try {
			return new URL(novelUrl).host;
		} catch {
			return '';
		}
	});

	async function submit(event: SubmitEvent) {
		event.preventDefault();
		error = '';

		const candidate = url.trim();
		if (!candidate) {
			error = 'Paste the URL of the first chapter.';
			return;
		}
		if (host && !candidate.includes(host)) {
			error = `That link is not on ${host}, which is where this novel lives.`;
			return;
		}

		submitting = true;
		try {
			const chapter = await setFirstChapter(novelId, candidate);
			url = '';
			onsuccess(chapter.number ?? 1);
		} catch (cause) {
			error =
				cause instanceof ManualExtractCancelled
					? 'Cancelled. The first chapter was not set.'
					: cause instanceof Error
						? cause.message
						: 'Could not translate that chapter.';
		} finally {
			submitting = false;
		}
	}
</script>

<Dialog
	{open}
	title="Set the first chapter"
	description="Later chapters are discovered from this one, so it only needs doing once."
	onclose={submitting ? () => {} : onclose}
>
	<form id="first-chapter" onsubmit={submit}>
		<div class="field">
			<label for="first-chapter-url">First chapter URL</label>
			<input
				id="first-chapter-url"
				class="input"
				type="url"
				inputmode="url"
				autocomplete="off"
				spellcheck="false"
				placeholder={host ? `https://${host}/...` : 'https://example.com/novel/1/chapter/1'}
				bind:value={url}
				disabled={submitting}
			/>
			{#if error}
				<p class="error" role="alert">{error}</p>
			{:else}
				<p class="hint">This chapter is translated straight away so you can start reading.</p>
			{/if}
		</div>
	</form>

	{#snippet footer()}
		<button class="btn" onclick={onclose} disabled={submitting}>Cancel</button>
		<button class="btn btn--primary" type="submit" form="first-chapter" disabled={submitting}>
			{submitting ? 'Translating…' : 'Translate it'}
		</button>
	{/snippet}
</Dialog>
