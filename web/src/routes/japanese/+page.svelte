<script lang="ts">
	import ShelfPage from '$components/ShelfPage.svelte';
	import { getNovels } from '$api/novels';
</script>

<svelte:head><title>Japanese novels · Arcane Translator</title></svelte:head>

<ShelfPage
	title="Japanese novels"
	fetcher={(page, limit) => getNovels(page, limit, { language: 'japanese' })}
	emptyTitle="No Japanese novels yet"
	emptyBody="Novels you add from Japanese sources will appear on this shelf."
/>
