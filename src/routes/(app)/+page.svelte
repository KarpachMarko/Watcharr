<script lang="ts">
	import WatchedList from "@/lib/WatchedList.svelte";
	import { store } from "@/store.svelte";
</script>

<svelte:head>
	<title>Watched List</title>
</svelte:head>

<WatchedList list={store.watchedList} />
