<script lang="ts">
	import { goto } from "$app/navigation";
	import Icon from "@/lib/Icon.svelte";
	import Poster from "@/lib/poster/Poster.svelte";
	import PosterList from "@/lib/poster/PosterList.svelte";
	import { store, clearActiveFilters } from "@/store.svelte";
	import { MediaTypeE, type Watched } from "@/types";
	import GamePoster from "./poster/GamePoster.svelte";
	import Spinner from "./Spinner.svelte";

	interface Props {
		list: Watched[];
		isLoading: boolean;
		isPublicList?: boolean;
	}

	let { list, isPublicList = false, isLoading = false }: Props = $props();

	/**
	 * TODO Figure out if we need this still, don't think so, we are probably going
	 * to not modify sort when items are changed to avoid having to reload the whole
	 * list and stops items jumping around too which might be better for UX idk.
	 * Callback for when a watched list item is updated through poster,
	 * this allows us to run the filt() func again so the sorting is
	 * updated.
	 */
	function itemUpdated() {
		console.debug("itemUpdated");
		// filt();
	}
</script>

<PosterList>
	{#if list?.length > 0}
		{#each list as w, i (w.id)}
			{#if w.media.type == MediaTypeE.igdbGame}
				<GamePoster
					bind:watched={list[i]}
					media={{
						id: w.media.ids.igdb,
						coverId: w.media.extPosterPath,
						name: w.media.name || "",
						summary: w.media.summary,
						firstReleaseDate: w.media.releaseDate,
					}}
					disableInteraction={isPublicList}
					fluidSize={true}
					pinned={w.pinned}
					onUpdated={itemUpdated}
				/>
			{:else if w.media.type == MediaTypeE.tmdbMovie || w.media.type == MediaTypeE.tmdbShow}
				<Poster
					bind:watched={list[i]}
					media={w.media}
					disableInteraction={isPublicList}
					fluidSize={true}
					pinned={w.pinned}
					onUpdated={itemUpdated}
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
