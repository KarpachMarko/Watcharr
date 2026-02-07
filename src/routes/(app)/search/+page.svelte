<script lang="ts">
	import { page } from "$app/state";
	import { afterNavigate, goto } from "$app/navigation";
	import axios, { type GenericAbortSignal } from "axios";
	import Poster from "@/lib/poster/Poster.svelte";
	import PosterList from "@/lib/poster/PosterList.svelte";
	import { store } from "@/store.svelte.js";
	import PageError from "@/lib/PageError.svelte";
	import Spinner from "@/lib/Spinner.svelte";
	import PersonPoster from "@/lib/poster/PersonPoster.svelte";
	import {
		MediaTypeE,
		SearchType,
		type Media,
		type PublicUser,
		type SearchRequest,
	} from "@/types";
	import UsersList from "@/lib/UsersList.svelte";
	import { onDestroy, onMount } from "svelte";
	import Error from "@/lib/Error.svelte";
	import GamePoster from "@/lib/poster/GamePoster.svelte";
	import Icon from "@/lib/Icon.svelte";
	import infScroll from "@/lib/util/infScroll.js";
	import paginatedLoader, {
		PaginatedLoaderRunFnAction,
	} from "@/lib/util/paginatedLoader.svelte.js";

	let { data } = $props();

	const scroll = infScroll({ callback: onScrollToBottom });
	const dataLoader = paginatedLoader<Media>(load);

	let searchType: SearchType | undefined = $state(
		data.type ? (data.type as SearchType) : SearchType.multi,
	);
	let nextLoadParams: SearchRequest = $derived({
		page: dataLoader.state.page + 1,
		query: store.searchQuery,
		type: searchType,
	});

	async function load(signal: GenericAbortSignal) {
		console.debug("load: loadParams:", nextLoadParams, "page params:", data);
		if (nextLoadParams.page === dataLoader.state.page) {
			console.warn("load: Already on this page, not loading it again!");
			return;
		}
		const r = await axios.get(`/search`, {
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

	function setActiveSearchFilterQueryParam(to: SearchType | undefined) {
		console.debug("setActiveSearchFilterQueryParam: to:", to);
		const curLocation = new URL(page.url);
		if (!to || searchType === to) {
			curLocation.searchParams.delete("type");
			searchType = undefined;
		} else {
			curLocation.searchParams.set("type", to);
			searchType = to;
		}
		goto(`?${curLocation.searchParams.toString()}`);
	}

	function setActiveSearchFilter(to: SearchType | undefined) {
		console.debug("setActiveSearchFilter: to:", to);
		setActiveSearchFilterQueryParam(to);
		dataLoader.runFn(PaginatedLoaderRunFnAction.Reset);
	}

	async function searchUsers(query: string) {
		return (await axios.get(`/user/search`, { params: { q: query } }))
			.data as PublicUser[];
	}

	onMount(() => {
		if (!store.searchQuery && data?.query) {
			store.searchQuery = data?.query;
		}
		dataLoader.runFn(PaginatedLoaderRunFnAction.Reset);
	});

	afterNavigate((e) => {
		if (!e.from?.route?.id?.toLowerCase()?.includes("/search")) {
			// AfterNavigate will also be called when this page is mounted,
			// but that won't work for us since the OnMount hook also runs
			// a clean search, which can cause errors when both ran at same
			// time. We can't remove the OnMount hook since it's the only
			// hook to be ran if watcharr is first loaded at a search url.
			// `e.type` is always `goto` (that's how we search) so we can't
			// use that. The only alternative to only run this hook after a
			// navigation on the search page (query change), seems to be
			// checking the `from` property an making sure it's from the
			// `/search` route already.
			return;
		}
		// After navigate ensure searchType is kept up to date (eg back
		// btn pressed in browser).
		searchType = data.type ? (data.type as SearchType) : SearchType.multi;
		console.log(
			"Query changed (or just loaded first query), performing search",
		);
		dataLoader.abortReq("navigated away");
		dataLoader.runFn(PaginatedLoaderRunFnAction.Reset);
	});

	onDestroy(() => {
		console.debug("SEARCH PAGE DESTROYED");
		store.searchQuery = "";
		scroll.destroy();
		dataLoader.abortReq("page destroyed");
	});
</script>

<svelte:head>
	<title
		>Search Results{store.searchQuery
			? ` for '${store.searchQuery}'`
			: ""}</title
	>
</svelte:head>

<!-- <span style="position: sticky;top: 70px;">{curPage} / {maxContentPage}</span> -->
<div class="content">
	<div class="inner">
		{#if data?.query}
			<!-- Uses data?.query instead of store.searchQuery,
			 	so that the debounce of search is respected. -->
			{#await searchUsers(data?.query) then results}
				{#if results?.length > 0}
					<UsersList users={results} />
				{/if}
			{:catch err}
				<PageError pretty="Failed to load users!" error={err} />
			{/await}

			<div
				class={`results-filters-header${dataLoader.state.reqLoading ? " search-running" : ""}`}
			>
				<h2>Results</h2>
				<div>
					<button
						class="plain"
						data-active={searchType === SearchType.movie}
						onclick={() => setActiveSearchFilter(SearchType.movie)}
					>
						<Icon i="film" wh={20} /> Movies
					</button>
					<button
						class="plain"
						data-active={searchType === SearchType.show}
						onclick={() => setActiveSearchFilter(SearchType.show)}
					>
						<Icon i="tv" wh={20} /> TV Shows
					</button>
					{#if store.serverFeatures?.games}
						<button
							class="plain"
							data-active={searchType === SearchType.game}
							onclick={() => setActiveSearchFilter(SearchType.game)}
						>
							<Icon i="gamepad" wh={20} /> Games
						</button>
					{/if}
					<button
						class="plain"
						data-active={searchType === SearchType.person}
						onclick={() => setActiveSearchFilter(SearchType.person)}
					>
						<Icon i="people-nocircle" wh={20} /> People
					</button>
				</div>
			</div>
			<PosterList>
				{#if dataLoader.state.data?.length > 0}
					{#each dataLoader.state.data as w, i (`${i}-${w.type}`)}
						{#if w.type === MediaTypeE.tmdbPerson}
							<PersonPoster
								id={w.ids.tmdb}
								name={w.name}
								path={w.extPosterPath}
							/>
						{:else if w.type === MediaTypeE.igdbGame && w.ids.igdb}
							<GamePoster
								media={{
									id: w.ids.igdb,
									coverId: w.extPosterPath,
									name: w.name || "",
									summary: w.summary,
									firstReleaseDate: w.releaseDate,
								}}
								bind:watched={dataLoader.state.data[i].watched}
								fluidSize
							/>
						{:else if dataLoader.state.data[i].type === MediaTypeE.tmdbMovie || dataLoader.state.data[i].type === MediaTypeE.tmdbShow}
							<Poster
								media={w}
								bind:watched={dataLoader.state.data[i].watched}
								fluidSize
							/>
						{/if}
					{/each}
				{:else if !dataLoader.state.reqLoading && !dataLoader.state.reqLoadError}
					<!-- If search is running or we have an error, no point in showing 'no results' message. -->
					<h2 class="norm" title="Hovering over me doesn't change the facts ;(">
						No Results!
					</h2>
				{/if}
			</PosterList>

			{#if dataLoader.state.reqLoading}
				<div style="margin-bottom: 60px;">
					<Spinner />
				</div>
			{/if}

			{#if dataLoader.state.reqLoadError}
				<div style="margin-bottom: 60px;">
					<Error
						pretty="Failed to load results!"
						error={dataLoader.state.reqLoadError}
						onRetry={() => {
							dataLoader.state.reqLoadError = undefined;
							dataLoader.runFn(
								PaginatedLoaderRunFnAction.ResetIfOnFirstOrNoPage,
							);
						}}
					/>
				</div>
			{/if}
		{:else}
			<h2>No Search Query!</h2>
		{/if}
	</div>
</div>

<style lang="scss">
	.results-filters-header {
		display: flex;
		flex-flow: row;
		flex-wrap: wrap;
		align-items: center;
		gap: 10px;

		div {
			display: flex;
			flex-flow: row;
			flex-wrap: wrap;
			gap: 10px;
			margin: 0 15px;

			button {
				display: flex;
				flex-flow: row;
				flex-wrap: wrap;
				gap: 8px;
				align-items: center;
				height: fit-content;
				padding: 8px 12px;
				border-radius: 8px;
				font-size: 14px;
				color: $text-color;
				fill: $text-color;
				transition:
					background-color 150ms ease,
					color 150ms ease,
					outline 150ms ease;

				&:hover,
				&[data-active="true"] {
					color: $bg-color;
					fill: $bg-color;
					background-color: $accent-color-hover;
				}

				&[data-active="true"] {
					outline: 3px solid $accent-color;
				}

				@media screen and (max-width: 500px) {
					flex-flow: column;
				}
			}

			@media screen and (max-width: 500px) {
				width: 100%;
				justify-content: center;
			}
		}

		&.search-running {
			button {
				opacity: 0.8;
				pointer-events: none;
			}
		}
	}

	.content {
		display: flex;
		width: 100%;
		justify-content: center;

		.inner {
			width: 100%;
			max-width: 1200px;

			h2 {
				margin-left: 15px;
			}
		}
	}
</style>
