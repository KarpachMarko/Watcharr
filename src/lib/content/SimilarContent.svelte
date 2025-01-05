<script lang="ts">
	import type { ContentType, TMDBMovieSimilar, TMDBShowSimilar } from "@/types";
	import HorizontalList from "../HorizontalList.svelte";
	import { store } from "@/store.svelte";
	import { getWatchedDependedProps } from "@/lib/util/helpers";
	import Poster from "../poster/Poster.svelte";

	interface Props {
		type: ContentType;
		similar: TMDBShowSimilar | TMDBMovieSimilar;
	}

	let { type, similar }: Props = $props();
</script>

{#if similar?.results?.length > 0}
	<HorizontalList title="Similar">
		{#each similar.results as content}
			<Poster
				media={{ ...content, media_type: type }}
				{...getWatchedDependedProps(content.id, type, store.watchedList)}
				small={true}
			/>
		{/each}
	</HorizontalList>
{/if}
