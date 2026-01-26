<script lang="ts">
	import Error from "@/lib/Error.svelte";
	import Icon from "@/lib/Icon.svelte";
	import Spinner from "@/lib/Spinner.svelte";
	import CreateTagModal from "@/lib/tag/CreateTagModal.svelte";
	import Tag from "@/lib/tag/Tag.svelte";
	import WatchedList from "@/lib/WatchedList.svelte";
	import { store } from "@/store.svelte.js";
	import axios from "axios";
	import { onDestroy } from "svelte";

	let { data } = $props();

	let tagEditModalShown = $state(false);

	let tag = $derived(store.tags.find((t) => t.id === data.tagId));

	let reqController: AbortController | undefined;

	async function getTag(id: number) {
		try {
			reqController = new AbortController();
			return (
				await axios.get<any>(`/tag/${id}/watched`, {
					signal: reqController.signal,
					params: {
						p: 0 + 1,
						...store.sortAndFiltersForQueryParams,
					},
				})
			).data;
		} catch (err) {
			console.error(`getTag: Failed!`, id, err);
			throw err;
		}
	}

	onDestroy(() => {
		reqController?.abort("page destroyed");
	});
</script>

<svelte:head>
	<title>{tag ? `${tag.name} - Tag` : "Tag"}</title>
</svelte:head>

{#if tag}
	<div class="content">
		<div class="inner">
			<div class="basic-ctr">
				<Icon i="tag" wh={20} />
				<Tag
					{tag}
					onClick={() => {
						tagEditModalShown = !tagEditModalShown;
					}}
				/>
			</div>
		</div>
	</div>

	{#await getTag(tag.id)}
		<Spinner />
	{:then dbtag}
		{#if dbtag.watched?.length > 0}
			<WatchedList list={dbtag.watched} isLoading={false} />
		{:else}
			<div class="content empty-tag">
				<div class="inner">
					<Icon i="ticket" wh={80} />
					<h2 class="norm">This tag is empty!</h2>
					<h4 class="norm">Add entries to this tag via its page.</h4>
				</div>
			</div>
		{/if}
	{:catch err}
		<Error pretty="Failed to load tag entries!" error={err} />
	{/await}
{:else}
	<div class="content">
		<div class="inner">
			<strong>Tag does not exist!</strong>
		</div>
	</div>
{/if}

{#if tagEditModalShown}
	<CreateTagModal
		existingTag={tag}
		onClose={() => (tagEditModalShown = false)}
	/>
{/if}

<style lang="scss">
	.content {
		display: flex;
		width: 100%;
		justify-content: center;

		.inner {
			display: flex;
			flex-flow: column;
			gap: 5px;
			justify-content: center;
			align-items: center;
			width: 100%;
			max-width: 1200px;
			margin: 20px 30px;
			margin-top: 0;
		}
	}

	.empty-tag {
		h2 {
			margin-top: 10px;
		}

		h4 {
			font-weight: normal;
		}
	}

	.basic-ctr {
		display: flex;
		align-items: center;
		justify-content: center;
		flex-wrap: wrap;
		gap: 10px;
		max-width: 300px;
		width: 100%;
		fill: $text-color;
	}
</style>
