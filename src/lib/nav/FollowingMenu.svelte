<script lang="ts">
	import { store } from "@/store.svelte";

	interface Props {
		close: () => {};
	}

	let { close }: Props = $props();
</script>

<div class="menu">
	<div>
		{#if store.follows?.length > 0}
			<h4 class="norm sm-caps">following</h4>
			<div class="list">
				{#each store.follows as f}
					<a
						href="/lists/{f.followedUser.id}/{f.followedUser.username}"
						onclick={() => close()}
					>
						{f.followedUser.username}
					</a>
				{/each}
			</div>
		{:else}
			<span style="margin-top: 0;">You are not following anyone.</span>
		{/if}
	</div>
</div>

<style lang="scss">
	div {
		width: 180px;

		&:before {
			right: 53px;
		}

		h4 {
			position: sticky;
			top: -10px;
			background-color: $bg-color;
		}

		.list {
			list-style: none;
			display: flex;
			flex-flow: column;
			width: 100%;
			height: 100%;

			a {
				overflow: hidden;
				white-space: nowrap;
				text-overflow: ellipsis;
			}
		}
	}
</style>
