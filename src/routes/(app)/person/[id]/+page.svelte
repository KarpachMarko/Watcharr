<script lang="ts">
	import Error from "@/lib/Error.svelte";
	import PageError from "@/lib/PageError.svelte";
	import Poster from "@/lib/poster/Poster.svelte";
	import PosterList from "@/lib/poster/PosterList.svelte";
	import Spinner from "@/lib/Spinner.svelte";
	import DropDown from "@/lib/DropDown.svelte";
	import type {
		TMDBPersonCombinedCredits,
		TMDBPersonCombinedCreditsCast,
		TMDBPersonDetails,
	} from "@/types";
	import axios from "axios";
	import Checkbox from "@/lib/Checkbox.svelte";
	import Icon from "@/lib/Icon.svelte";

	let { data } = $props();

	let person: TMDBPersonDetails | undefined = $state();
	let pageError: Error | undefined = $state();
	let sortOption = $state("Vote count");
	let credits: TMDBPersonCombinedCredits | undefined = $state();
	let onMyListFilter = $state(false);

	$effect(() => {
		if (data.personId) {
			fetchPersonData();
		}
	});

	$effect(() => {
		if (sortOption && credits) {
			sortCredits(sortOption);
		}
	});

	async function fetchPersonData() {
		try {
			person = undefined;
			pageError = undefined;
			if (!data.personId) {
				return;
			}
			person = await getPerson(data.personId);
			await updatePersonCredits();
			sortCredits(sortOption);
		} catch (err: any) {
			person = undefined;
			pageError = err;
		}
	}

	async function getPerson(id: number) {
		return (await axios.get(`/content/person/${id}`)).data as TMDBPersonDetails;
	}

	async function updatePersonCredits() {
		credits = (await axios.get(`/content/person/${data.personId}/credits`))
			.data as TMDBPersonCombinedCredits;
		credits.cast = credits.cast.filter(
			(v, i, a) => a.findIndex((t) => t.id === v.id) === i,
		); // remove duplicate entries. If an actor has multiple roles in a single movie, it would otherwise show up multiple times
	}

	function newestOldestSort(
		a: TMDBPersonCombinedCreditsCast,
		b: TMDBPersonCombinedCreditsCast,
		/**
		 * 0 = Newest,
		 * 1 = Oldest
		 */
		n: 0 | 1,
	) {
		const dateA = new Date(a.release_date || a.first_air_date).valueOf();
		const dateB = new Date(b.release_date || b.first_air_date).valueOf();

		// Assume missing release date means future release (TBD)
		if (
			!a.release_date &&
			!a.first_air_date &&
			!b.release_date &&
			!b.first_air_date
		) {
			// Both releases have no date, return as equals
			// here to avoid an infinite loop.
			return 0;
		}
		if (!a.release_date && !a.first_air_date) return n === 0 ? -1 : 1;
		if (!b.release_date && !b.first_air_date) return n === 0 ? 1 : -1;

		if (n === 0) {
			return dateB - dateA;
		} else {
			return dateA - dateB;
		}
	}

	function sortCredits(sortOption: string) {
		if (!credits || !credits.cast) return;
		switch (sortOption) {
			case "Vote count":
				credits.cast.sort((a, b) => b.vote_count - a.vote_count);
				break;
			case "Newest":
				credits.cast.sort((a, b) => newestOldestSort(a, b, 0));
				break;
			case "Oldest":
				credits.cast.sort((a, b) => newestOldestSort(a, b, 1));
				break;
		}
		credits.cast = credits.cast;
	}
</script>

<svelte:head>
	<title>{person?.name ? `${person.name} - ` : ""}Person</title>
</svelte:head>

<div>
	{#if pageError}
		<PageError pretty="Failed to load person!" error={pageError} />
	{:else if !person}
		<Spinner />
	{:else if Object.keys(person).length > 0}
		{#if Object.keys(person).length > 0}
			<div class="content">
				<img
					class="backdrop"
					src={"https://www.themoviedb.org/t/p/w1920_and_h800_multi_faces" +
						person.profile_path}
					alt=""
				/>
				<div class="vignette"></div>

				<div class="details-container">
					<img
						class="poster"
						src={"https://image.tmdb.org/t/p/w500" + person.profile_path}
						alt=""
					/>

					<div class="details">
						<span class="title-container">
							<a href={person.homepage} target="_blank">{person.name}</a>
							<span></span>
						</span>

						{#if person.biography}
							<span style="font-weight: bold; font-size: 14px;">Biography</span>
							<!-- Show just the first paragraph -->
							<p>{person.biography?.split("\n")[0]}</p>
						{/if}

						<div class="detail-info">
							{#if person.known_for_department}
								<div>
									<span>Department</span>
									<span>{person.known_for_department}</span>
								</div>
							{/if}
							{#if person.place_of_birth}
								<div>
									<span>Born In</span>
									<span>{person.place_of_birth}</span>
								</div>
							{/if}
							{#if person.birthday}
								<div>
									<span>Born On</span>
									<span
										>{new Date(
											Date.parse(person.birthday),
										).toLocaleDateString()}</span
									>
								</div>
							{/if}
							{#if person.deathday}
								<div>
									<span>Died On</span>
									<span
										>{new Date(
											Date.parse(person.deathday),
										).toLocaleDateString()}</span
									>
								</div>
							{/if}
						</div>
					</div>
				</div>
			</div>
			{#if credits}
				{#if credits?.cast?.length > 0}
					<div class="filters">
						<div class="listFilter">
							<span>On my list</span>
							<Checkbox name="On my list" bind:value={onMyListFilter} />
						</div>
						<DropDown
							bind:active={sortOption}
							placeholder="Vote count"
							options={["Vote count", "Newest", "Oldest"]}
							isDropDownItem={false}
							showActiveElementsInOptions={true}
						/>
					</div>
					<div class="page">
						<PosterList>
							{#each credits.cast as c, i (c.id)}
								<Poster
									media={c}
									bind:watched={credits.cast[i].watched}
									fluidSize
									hideIfNotOnList={onMyListFilter}
								/>
							{/each}
						</PosterList>
					</div>
				{:else}
					<div class="no-credits-message">
						<Icon i="star" wh={80} />
						<h2 class="norm">We found no credits!</h2>
						<h4 class="norm">It seems that this person has no credits.</h4>
					</div>
				{/if}
			{:else}
				<Spinner />
			{/if}
		{:else}
			person not found
		{/if}
	{:else}
		<Error error="Person not found" pretty="Person not found" />
	{/if}
</div>

<style lang="scss">
	.filters {
		align-items: center;
		display: flex;
		justify-content: flex-end;
		gap: 30px;
		margin: 0 auto;
		padding-left: 20px;
		padding-right: 20px;
		width: 100%;
		/* Same as in PosterList */
		max-width: 1200px;

		.listFilter {
			display: flex;
			align-items: center;
			gap: 8px;
		}
	}

	.content {
		position: relative;
		color: white;
		margin-bottom: 15px;

		img.backdrop {
			position: absolute;
			left: 0;
			top: 0;
			z-index: -2;
			width: 100%;
			height: 100%;
			object-fit: cover;
			filter: blur(4px) grayscale(80%);
			/* mix-blend-mode: multiply; */
		}

		.vignette {
			position: absolute;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			background-color: rgba($color: #000000, $alpha: 0.7);
			z-index: -1;
		}

		.details-container {
			display: flex;
			flex-flow: row;
			gap: 35px;
			max-width: 1100px;
			padding: 40px 80px;
			margin-left: auto;
			margin-right: auto;

			img.poster {
				width: 235px;
				height: 100%;
				box-shadow: 0px 0px 14px -4px #9c8080;
				border-radius: 12px;
			}

			.details {
				display: flex;
				flex-flow: column;
				gap: 5px;

				.title-container {
					a {
						color: white;
						text-decoration: none;
						font-size: 30px;
						font-weight: bold;
						padding-right: 3px;
					}

					span {
						font-size: 20px;
						color: rgba($color: #fff, $alpha: 0.7);
					}
				}

				p {
					font-size: 14px;
					max-height: 200px;
					overflow: hidden;
					text-overflow: ellipsis;
					white-space: pre-line;
				}

				.detail-info {
					display: flex;
					flex-wrap: wrap;
					gap: 35px;
					margin-top: 10px;
					font-size: 14px;

					div {
						display: flex;
						flex-flow: column;

						span:first-child {
							font-weight: bold;
						}
					}
				}
			}

			@media screen and (max-width: 700px) {
				padding: 40px;
			}

			@media screen and (max-width: 570px) {
				flex-flow: column;
				align-items: center;
			}
		}
	}

	.page {
		display: flex;
		flex-flow: column;
		align-items: center;
		gap: 30px;
		padding: 10px 0px;
	}

	.no-credits-message {
		display: flex;
		flex-flow: column;
		gap: 5px;
		align-items: center;
		margin-top: 20px;

		h2 {
			margin-top: 10px;
		}

		h4 {
			font-weight: normal;
		}
	}
</style>
