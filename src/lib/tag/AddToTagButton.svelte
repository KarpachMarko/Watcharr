<script lang="ts">
	import type { Watched } from "@/types";
	import tooltip from "../actions/tooltip";
	import Icon from "../Icon.svelte";
	import TagMenu from "./TagMenu.svelte";
	import { tagWatched, untagWatched } from "./api";
	import { onMount } from "svelte";

	interface Props {
		watchedItem: Watched;
	}

	let { watchedItem }: Props = $props();

	let menuOpen = $state(false);

	onMount(() => {
		const onScroll = () => {
			menuOpen = false;
		};

		window.addEventListener("scroll", onScroll);

		return () => {
			window.removeEventListener("scroll", onScroll);
		};
	});
</script>

<div>
	<button
		onclick={() => (menuOpen = !menuOpen)}
		use:tooltip={{
			text: `Add to a Tag`,
			pos: "bot",
		}}
	>
		<Icon i={"tag"} wh={19} />
	</button>

	{#if menuOpen}
		<TagMenu
			titleText="Add To Tag"
			classes="from-add-to-tag-btn"
			selectedTags={watchedItem.tags}
			onTagClick={(tag, remove) => {
				console.debug("Tag: Adding content to tag. Remove?:", remove);
				if (remove) {
					untagWatched(watchedItem.id, tag);
				} else {
					tagWatched(watchedItem.id, tag);
				}
			}}
		/>
	{/if}
</div>

<style lang="scss">
	div {
		position: relative;
	}
</style>
