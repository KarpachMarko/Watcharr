<script lang="ts">
	import PageError from "@/lib/PageError.svelte";
	import Spinner from "@/lib/Spinner.svelte";
	import HorizontalList from "@/lib/HorizontalList.svelte";
	import {
		GameWebsiteCategory,
		type GameDetailsResponseWithWatched,
		type WatchedStatus,
	} from "@/types";
	import axios from "axios";
	import Activity from "@/lib/Activity.svelte";
	import Title from "@/lib/content/Title.svelte";
	import VideoEmbedModal from "@/lib/content/VideoEmbedModal.svelte";
	import Error from "@/lib/Error.svelte";
	import FollowedThoughts from "@/lib/content/FollowedThoughts.svelte";
	import { removeWatched, updatePlayed } from "@/lib/util/api.js";
	import GamePoster from "@/lib/poster/GamePoster.svelte";
	import tooltip from "@/lib/actions/tooltip.js";
	import Icon from "@/lib/Icon.svelte";
	import AddToTagButton from "@/lib/tag/AddToTagButton.svelte";
	import PageBackdrop from "@/lib/generic/PageBackdrop.svelte";
	import Poster from "@/lib/content/Poster.svelte";
	import MyReview from "@/lib/content/MyReview.svelte";

	let { data } = $props();

	let trailer: string | undefined = $state();
	let trailerShown = $state(false);
	let game: GameDetailsResponseWithWatched | undefined = $state();
	let pageError: Error | undefined = $state();
	let backdropSrc = $derived.by(() => {
		const base = "https://images.igdb.com/igdb/image/upload/t_1080p/";

		if (game?.artworks && game?.artworks?.length > 0) {
			return (
				base +
				game.artworks[Math.floor(Math.random() * game.artworks.length)]
					.image_id +
				".jpg"
			);
		} else if (game?.cover?.image_id) {
			return base + game.cover.image_id + ".jpg";
		}
	});

	$effect(() => {
		(async () => {
			try {
				game = undefined;
				pageError = undefined;
				if (!data.gameId) {
					return;
				}
				const resp = (await axios.get(`/game/${data.gameId}`))
					.data as GameDetailsResponseWithWatched;
				if (resp.videos?.length > 0) {
					const t = resp.videos.find(
						(v) => v.name?.toLowerCase() === "trailer",
					);
					// Doc says the video_id is "usually youtube", so we are gonna go with that assumption too ( 0 _ 0 )
					if (t?.video_id) {
						trailer = `https://www.youtube.com/embed/${t?.video_id}`;
					}
				}
				game = resp;
			} catch (err: any) {
				game = undefined;
				pageError = err;
			}
		})();
	});

	async function contentChanged(
		newStatus?: WatchedStatus,
		newRating?: number,
		newThoughts?: string,
		pinned?: boolean,
	) {
		if (!data.gameId) {
			console.error("contentChanged: no tvId");
			return;
		}
		if (!game) {
			console.error("contentChanged: no show");
			return;
		}
		game.watched = await updatePlayed(game.watched, {
			igdbId: data.gameId,
			status: newStatus,
			rating: newRating,
			thoughts: newThoughts,
			pinned: pinned,
		});
	}

	async function deleteWatched() {
		if (game?.watched) {
			if (await removeWatched(game.watched.id)) {
				game.watched = undefined;
			}
			return;
		}
		console.error("deleteWatched: no wlistItem.. can't delete");
	}
</script>

<svelte:head>
	<title>{game?.name ? `${game.name} - ` : ""}Game</title>
</svelte:head>

{#if pageError}
	<PageError pretty="Failed to load game!" error={pageError} />
{:else if !game}
	<Spinner />
{:else if Object.keys(game).length > 0}
	{#if backdropSrc}
		<PageBackdrop src={backdropSrc} />
	{/if}
	<div>
		<div class="content">
			<div class="details-wrap">
				<div class="details-container">
					<Poster
						src={"https://images.igdb.com/igdb/image/upload/t_cover_big/" +
							game.cover.image_id +
							".jpg"}
					/>

					<div class="details">
						<Title
							title={game.name}
							homepage={game.websites?.find(
								(w) => w.category == GameWebsiteCategory.Official,
							)?.url}
							releaseYear={new Date(game.first_release_date).getFullYear()}
							voteAverage={game.rating}
							voteCount={game.rating_count}
						/>

						<span class="quick-info">
							{#if game.genres?.length > 0}
								<div>
									{#each game.genres as g, i}
										<span
											>{g.name}{i !== game.genres.length - 1 ? ", " : ""}</span
										>
									{/each}
								</div>
							{:else}
								<span>Unknown Genres</span>
							{/if}
							<span></span>
							<div>
								{#if game.game_modes?.length > 0}
									{#each game.game_modes as g, i}
										<span
											>{g.name}{i !== game.game_modes.length - 1
												? ", "
												: ""}</span
										>
									{/each}
								{:else}
									<span>Unknown Game Modes</span>
								{/if}
							</div>
						</span>

						<p>{game.summary}</p>

						<div class="btns">
							{#if trailer}
								<button onclick={() => (trailerShown = !trailerShown)}
									>View Trailer</button
								>
								{#if trailerShown}
									<VideoEmbedModal
										embed={trailer}
										closed={() => (trailerShown = false)}
									/>
								{/if}
							{/if}
							{#if game.watched}
								<div class="other-side">
									<AddToTagButton watchedItem={game.watched} />
									<button
										onclick={() => {
											if (game?.watched?.pinned) {
												contentChanged(undefined, undefined, undefined, false);
											} else {
												contentChanged(undefined, undefined, undefined, true);
											}
										}}
										use:tooltip={{
											text: `${game.watched?.pinned ? "Unpin from" : "Pin to"} top of list`,
											pos: "bot",
										}}
									>
										<Icon i={game.watched?.pinned ? "unpin" : "pin"} wh={19} />
									</button>
									<button
										class="delete-btn"
										onclick={() => deleteWatched()}
										use:tooltip={{ text: "Delete", pos: "bot" }}
									>
										<Icon i="trash" wh={19} />
									</button>
								</div>
							{/if}
						</div>

						<!-- <ProvidersList providers={game["watch/providers"]} /> -->
					</div>
				</div>
			</div>

			<MyReview
				watched={game.watched}
				contentTitle={game.name}
				onRatingChanged={(n) => contentChanged(undefined, n)}
				onStatusChanged={(n) => contentChanged(n)}
				onThoughtsChanged={(newThoughts) => {
					return contentChanged(undefined, undefined, newThoughts);
				}}
			/>
		</div>

		<div class="page">
			{#if data.gameId}
				<FollowedThoughts mediaType="game" mediaId={data.gameId} />
			{/if}

			{#if game.similar_games?.length > 0}
				<HorizontalList title="Similar">
					{#each game.similar_games as g, i}
						<GamePoster
							media={{
								id: g.id,
								coverId: g.cover.image_id,
								name: g.name,
								summary: g.summary,
								firstReleaseDate: g.first_release_date,
							}}
							bind:watched={game.similar_games[i].watched}
							small={true}
						/>
					{/each}
				</HorizontalList>
			{/if}

			{#if game.watched}
				<Activity bind:activity={game.watched.activity} />
			{/if}
		</div>
	</div>
{:else}
	<Error error="Game not found" pretty="Game not found" />
{/if}

<style lang="scss">
	@use "../../../../lib/content/page.scss";

	.content {
		position: relative;
		color: white;

		.details-container .details {
			.quick-info {
				display: flex;
				gap: 10px;
				margin-bottom: 8px;
			}

			p {
				font-size: 16px;
				margin-bottom: 18px;
			}

			.btns {
				display: flex;
				flex-flow: row;
				flex-wrap: wrap;
				gap: 8px;
				margin-top: auto;

				a.btn,
				button {
					max-width: fit-content;
					overflow: hidden;
					animation: 50ms cubic-bezier(0.86, 0, 0.07, 1) forwards otherbtn;
					white-space: nowrap;
					gap: 6px;
					justify-content: flex-start;
					font-size: 14px;

					@keyframes otherbtn {
						from {
							width: 0px;
						}
						to {
							width: 100%;
						}
					}
				}

				.other-side {
					display: flex;
					flex-flow: row;
					gap: 8px;

					@media screen and (min-width: 900px) {
						margin-left: auto;
					}
				}

				.delete-btn {
					&:hover {
						color: $error;
					}
				}
			}
		}
	}

	.page {
		display: flex;
		flex-flow: column;
		align-items: center;
		margin-left: auto;
		margin-right: auto;
		gap: 30px;
		padding: 20px 50px;
		max-width: 1200px;

		@media screen and (max-width: 500px) {
			padding: 20px;
		}
	}

	.creators {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		gap: 35px;
		margin: 10px 60px;

		div {
			display: flex;
			flex-flow: column;
			min-width: 150px;

			span:first-child {
				font-weight: bold;
			}
		}
	}
</style>
