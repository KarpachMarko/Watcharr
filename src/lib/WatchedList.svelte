<script lang="ts">
	import { goto } from "$app/navigation";
	import Icon from "@/lib/Icon.svelte";
	import Poster from "@/lib/poster/Poster.svelte";
	import PosterList from "@/lib/poster/PosterList.svelte";
	import { store, clearActiveFilters } from "@/store.svelte";
	import { type Watched } from "@/types";
	import Spinner from "./Spinner.svelte";

	interface Props {
		list: Watched[];
		isLoading: boolean;
		isPublicList?: boolean;
	}

	let { list, isPublicList = false, isLoading = false }: Props = $props();
</script>

<PosterList>
	{#if list?.length > 0}
		{#each list as w, i (w.id)}
			{#if w.media}
				<Poster
					bind:watched={list[i]}
					media={w.media}
					disableInteraction={isPublicList}
					fluidSize={true}
					pinned={w.pinned}
				/>
			{/if}
		{/each}
	{:else if !isLoading}
		{@const hasFiltersActive = store.hasActiveFilters}
		<div class="empty-list">
			<Icon i={hasFiltersActive ? "filter-circle" : "reel"} wh={80} />
			{#if isPublicList}
				<h2 class="norm">This list is empty!</h2>
				<h4 class="norm">
					Come back later to see if they have added anything.
				</h4>
			{:else}
				<h2 class="norm">Your list looks empty!</h2>
				<h4 class="norm">
					Try {`${hasFiltersActive ? "removing your active filters or" : ""}`} searching
					for something you would like to add.
				</h4>
				{#if !hasFiltersActive}
					<button onclick={() => goto("/import")}>Import</button>
				{/if}
			{/if}
			{#if hasFiltersActive}
				<button onclick={() => clearActiveFilters()}>Clear Filters</button>
			{/if}
		</div>
	{/if}
</PosterList>

{#if isLoading}
	<div style="margin-bottom: 60px;">
		<Spinner />
	</div>
{/if}

<style lang="scss">
	.empty-list {
		display: flex;
		flex-flow: column;
		gap: 5px;
		align-items: center;
		max-width: 400px;

		h2 {
			margin-top: 10px;
		}

		h4 {
			font-weight: normal;
			text-align: center;
		}

		button {
			width: max-content;
			padding-left: 20px;
			padding-right: 20px;
			margin-top: 15px;
		}
	}
</style>
