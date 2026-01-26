<script lang="ts">
	import Error from "@/lib/Error.svelte";
	import infScroll from "@/lib/util/infScroll";
	import paginatedLoader from "@/lib/util/paginatedLoader.svelte";
	import WatchedList from "@/lib/WatchedList.svelte";
	import { store } from "@/store.svelte";
	import type { Watched } from "@/types";
	import axios, { type GenericAbortSignal } from "axios";
	import { onDestroy, untrack } from "svelte";

	const scroll = infScroll({ callback: onScrollToBottom });
	const dataLoader = paginatedLoader<Watched>(load);

	let nextLoadParams: {
		p: number;
		[x: string]: any;
	} = $derived({
		p: dataLoader.state.page + 1,
		...store.sortAndFiltersForQueryParams,
	});

	async function load(signal: GenericAbortSignal) {
		console.debug("load: loadParams:", nextLoadParams);
		if (nextLoadParams.p === dataLoader.state.page) {
			console.warn("load: Already on this page, not loading it again!");
			return;
		}
		const r = await axios.get(`/watched`, {
			params: nextLoadParams,
			signal,
		});
		scroll.dataLoaded();
		return r;
	}

	async function onScrollToBottom() {
		// If an error is being shown, no more infinite scroll.
		if (dataLoader.state.reqLoadError) {
			return;
		}
		dataLoader.runFn();
	}

	$effect(() => {
		// When our sort/filter query params change,
		// load our list again.
		// Since it exists at load, this performs our
		// initial load of data too.
		if (store.sortAndFiltersForQueryParams) {
			untrack(() => {
				// We don't want to trigger another re-run of this
				// effect when state inside these funcs changes.
				dataLoader.reset();
				dataLoader.runFn();
			});
		}
	});

	onDestroy(() => {
		console.log("MAIN PAGE DESTROYED");
		scroll.destroy();
		dataLoader.abortReq("page destroyed");
	});
</script>

<svelte:head>
	<title>Watched List</title>
</svelte:head>

<span style="position: fixed; top: 80px; background-color: white; z-index: 60;"
	><b>listPage</b>: {dataLoader.state.page} listPageMax: {dataLoader.state
		.pageMax} listLoading:
	{dataLoader.state.reqLoading}
	<b>sort:</b>
	{JSON.stringify(store.activeSort)} <b>filter:</b>
	{JSON.stringify(store.activeFilters)} <b>queryp:</b>
	{JSON.stringify(store.sortAndFiltersForQueryParams)}</span
>

{#if dataLoader.state.data.length >= 0 && !dataLoader.state.reqLoadError}
	<!-- Hide WatchedList if there is a load error and we have no
	 	 items in our list array. WatchedList stays rendered if we
		 do have items because we could get a load error loading
		 the second page for example. -->
	<WatchedList
		list={dataLoader.state.data}
		isLoading={dataLoader.state.reqLoading}
	/>
{/if}

{#if dataLoader.state.reqLoadError}
	<div style="margin-bottom: 60px;">
		<Error
			pretty="Failed to load results!"
			error={dataLoader.state.reqLoadError}
			onRetry={() => {
				dataLoader.state.reqLoadError = undefined;
				dataLoader.runFn();
			}}
		/>
	</div>
{/if}

<!-- TODO: A 'Thats It! All {watcheds_count} of your records.' message
 or similar message to indicate that we are now at the bottom of the
 users list (only show when listPage === listPageMax & we have any items.. etc
 probs best adding this to WatchedList component.. but message probs needs to be
 a bit different depending on viewing own list or someone elses). -->
