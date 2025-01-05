<script lang="ts">
	interface Props {
		children?: import("svelte").Snippet;
	}

	let { children }: Props = $props();
</script>

<div>
	{@render children?.()}
</div>

<style lang="scss">
	div {
		display: flex;
		flex-flow: row;
		gap: 12px;
		margin-top: 15px;
		flex-wrap: wrap;
	}
</style>
