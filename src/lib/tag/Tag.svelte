<script lang="ts">
	import { run } from "svelte/legacy";

	import type { Tag } from "@/types";

	interface Props {
		tag: Tag;
		onClick?: () => void | undefined;
	}

	let { tag, onClick = undefined! }: Props = $props();

	let tagBtn: HTMLButtonElement = $state();

	run(() => {
		if (tagBtn) {
			if (tag.color) {
				tagBtn.style.color = tag.color;
			}
			if (tag.bgColor) {
				tagBtn.style.background = tag.bgColor;
			}
		}
	});
</script>

<button
	bind:this={tagBtn}
	class={`plain`}
	onclick={() => {
		if (typeof onClick === "function") {
			onClick();
		}
	}}
>
	{tag.name}
</button>

<style lang="scss">
	button {
		text-transform: capitalize;
		position: relative;
		width: max-content;
		border-radius: 8px;
		padding: 5px 8px;
		text-wrap: wrap;
		word-break: break-word;
		transition: opacity 150ms ease-in-out;

		&:hover {
			opacity: 0.8;
		}
	}
</style>
