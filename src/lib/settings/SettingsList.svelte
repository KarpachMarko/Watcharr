<script lang="ts">
	interface Props {
		children?: import("svelte").Snippet;
	}

	let { children }: Props = $props();
</script>

<div class="settings">
	{@render children?.()}
</div>

<style lang="scss">
	.settings {
		display: flex;
		flex-flow: column;
		gap: 15px;
		width: 100%;
	}
</style>
