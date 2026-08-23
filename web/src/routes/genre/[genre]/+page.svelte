<script lang="ts">
	import { page } from '$app/state';
	import ShelfPage from '$components/ShelfPage.svelte';
	import { getNovels } from '$api/novels';

	const genre = $derived(decodeURIComponent(page.params.genre ?? ''));

	// A new identity per genre so ShelfPage's effect reloads on navigation.
	const fetcher = $derived((pageNumber: number, limit: number) =>
		getNovels(pageNumber, limit, { genre })
	);
</script>

<svelte:head><title>{genre} · Arcane Translator</title></svelte:head>

<ShelfPage
	title={genre}
	{fetcher}
	emptyTitle="No novels tagged {genre}"
	emptyBody="Genres are filled in as chapters are translated, so this shelf grows as you read."
/>
