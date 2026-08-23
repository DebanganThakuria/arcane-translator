<script lang="ts">
	interface Props {
		count?: number;
	}

	let { count = 12 }: Props = $props();
</script>

<ul class="grid" aria-hidden="true">
	{#each { length: count } as _, i (i)}
		<li>
			<div class="skeleton art"></div>
			<div class="skeleton line"></div>
			<div class="skeleton line short"></div>
		</li>
	{/each}
</ul>

<style>
	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(9.5rem, 1fr));
		gap: var(--sp-5) var(--sp-4);
	}

	@media (max-width: 480px) {
		.grid {
			grid-template-columns: repeat(2, 1fr);
			gap: var(--sp-4) var(--sp-3);
		}
	}

	li {
		display: grid;
		gap: var(--sp-2);
	}

	.art {
		aspect-ratio: 2 / 3;
		border-radius: var(--r-2);
	}

	.line {
		height: 0.75rem;
	}

	.short {
		width: 55%;
	}
</style>
