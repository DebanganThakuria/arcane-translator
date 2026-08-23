<script lang="ts">
	import ShelfPage from '$components/ShelfPage.svelte';
	import { getNovels } from '$api/novels';
</script>

<svelte:head><title>Chinese novels · Arcane Translator</title></svelte:head>

<ShelfPage
	title="Chinese novels"
	fetcher={(page, limit) => getNovels(page, limit, { language: 'chinese' })}
	emptyTitle="No Chinese novels yet"
	emptyBody="Novels you add from Chinese sources will appear on this shelf."
/>
