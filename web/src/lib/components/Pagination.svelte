<script lang="ts">
	import { CaretLeft, CaretRight } from 'phosphor-svelte';
	import { PER_PAGE_OPTIONS, pageWindow, type Pagination } from '$utils/pagination.svelte';

	interface Props {
		pagination: Pagination;
	}

	let { pagination }: Props = $props();

	const data = $derived(pagination.data);
	const pages = $derived(pageWindow(pagination.current, data.total_pages));
</script>

{#if data.total_pages > 1}
	<nav class="pager" aria-label="Pagination">
		<button
			class="btn btn--icon"
			disabled={pagination.current <= 1 || pagination.loading}
			onclick={() => pagination.goTo(pagination.current - 1)}
			aria-label="Previous page"
		>
			<CaretLeft size={15} weight="bold" />
		</button>

		<ol class="pages">
			{#each pages as page, i (`${page}-${i}`)}
				{#if page === -1}
					<li class="gap" aria-hidden="true">&hellip;</li>
				{:else}
					<li>
						<button
							class="page numeral"
							class:current={page === pagination.current}
							aria-current={page === pagination.current ? 'page' : undefined}
							disabled={pagination.loading}
							onclick={() => pagination.goTo(page)}
						>
							{page}
						</button>
					</li>
				{/if}
			{/each}
		</ol>

		<button
			class="btn btn--icon"
			disabled={pagination.current >= data.total_pages || pagination.loading}
			onclick={() => pagination.goTo(pagination.current + 1)}
			aria-label="Next page"
		>
			<CaretRight size={15} weight="bold" />
		</button>

		<div class="per-page">
			<label for="per-page">Per page</label>
			<select
				id="per-page"
				class="select"
				value={pagination.limit}
				onchange={(event) => pagination.setLimit(Number(event.currentTarget.value))}
			>
				{#each PER_PAGE_OPTIONS as option (option)}
					<option value={option}>{option}</option>
				{/each}
			</select>
		</div>
	</nav>
{/if}

<style>
	.pager {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: var(--sp-2);
		padding-top: var(--sp-5);
		margin-top: var(--sp-5);
		border-top: 1px solid var(--line);
	}

	.pages {
		display: flex;
		gap: var(--sp-1);
	}

	.page {
		min-width: var(--control-h);
		height: var(--control-h);
		padding-inline: var(--sp-2);
		background: transparent;
		border: 1px solid transparent;
		border-radius: var(--r-1);
		color: var(--ink-2);
		cursor: pointer;
		transition:
			background var(--dur-1) var(--ease),
			color var(--dur-1) var(--ease);
	}

	.page:hover:not(:disabled) {
		background: var(--accent-wash);
		color: var(--ink);
	}

	.page.current {
		background: var(--accent);
		border-color: var(--accent);
		color: var(--accent-ink);
		font-weight: 600;
	}

	.gap {
		display: grid;
		place-items: center;
		min-width: 1.5rem;
		color: var(--ink-3);
	}

	.per-page {
		display: flex;
		align-items: center;
		gap: var(--sp-2);
		margin-left: auto;
		font-size: var(--step--1);
		color: var(--ink-2);
	}

	.per-page .select {
		width: auto;
		height: var(--control-h-sm);
		font-size: var(--step--1);
	}

	@media (max-width: 560px) {
		.per-page {
			margin-left: 0;
			width: 100%;
		}
	}
</style>
