<script lang="ts">
	import { page } from '$app/state';
</script>

<svelte:head><title>Not found · Arcane Translator</title></svelte:head>

<div class="page wrap">
	<p class="code numeral">{page.status}</p>
	<h1 class="display">
		{page.status === 404 ? 'There is no page here' : 'Something went wrong'}
	</h1>
	<p class="detail">{page.error?.message ?? 'The link may be out of date.'}</p>
	<a class="btn btn--primary" href="/">Back to the shelf</a>
</div>

<style>
	.wrap {
		display: grid;
		justify-items: start;
		gap: var(--sp-3);
		max-width: 34rem;
		padding-block: var(--sp-8);
	}

	.code {
		font-size: var(--step--1);
		color: var(--accent);
	}

	h1 {
		font-size: clamp(2rem, 5vw, 3rem);
	}

	.detail {
		color: var(--ink-2);
		margin-bottom: var(--sp-2);
	}
</style>
