<script lang="ts">
	import type { Snippet } from "svelte";

	interface Props {
		title: string;
		row?: boolean;
		children: Snippet;
	}

	let { title, row = true, children }: Props = $props();
</script>

<div class="results-filters-header" class:row>
	<h2>{title}</h2>
	{@render children()}
</div>

<style lang="scss">
	.results-filters-header {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: 10px;
		margin: 0 15px;

		&.row {
			flex-flow: row;
		}
	}
</style>
