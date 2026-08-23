<script lang="ts">
	import ShelfPage from '$components/ShelfPage.svelte';
	import { getNovels } from '$api/novels';
</script>

<svelte:head><title>Korean novels · Arcane Translator</title></svelte:head>

<ShelfPage
	title="Korean novels"
	fetcher={(page, limit) => getNovels(page, limit, { language: 'korean' })}
	emptyTitle="No Korean novels yet"
	emptyBody="Novels you add from Korean sources will appear on this shelf."
/>
