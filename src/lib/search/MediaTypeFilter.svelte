<script lang="ts">
	import Icon from "../Icon.svelte";
	import { store } from "@/store.svelte";

	type FilterType = "movie" | "show" | "game" | "person";

	interface Props {
		active?: string;
		disabled?: boolean;
		onChange: (nowActive: FilterType) => void;
	}

	let { active, disabled, onChange }: Props = $props();
</script>

<div class:disabled>
	<button
		class="plain"
		data-active={active === "movie"}
		onclick={() => onChange("movie")}
	>
		<Icon i="film" wh={20} /> Movies
	</button>
	<button
		class="plain"
		data-active={active === "show"}
		onclick={() => onChange("show")}
	>
		<Icon i="tv" wh={20} /> TV Shows
	</button>
	{#if store.serverFeatures?.games}
		<button
			class="plain"
			data-active={active === "game"}
			onclick={() => onChange("game")}
		>
			<Icon i="gamepad" wh={20} /> Games
		</button>
	{/if}
	<button
		class="plain"
		data-active={active === "person"}
		onclick={() => onChange("person")}
	>
		<Icon i="people-nocircle" wh={20} /> People
	</button>
</div>

<style lang="scss">
	div {
		display: flex;
		flex-flow: row;
		flex-wrap: wrap;
		gap: 10px;
		margin: 0 15px;

		button {
			display: flex;
			flex-flow: row;
			flex-wrap: wrap;
			gap: 8px;
			align-items: center;
			height: fit-content;
			padding: 8px 12px;
			border-radius: 8px;
			font-size: 14px;
			color: $text-color;
			fill: $text-color;
			transition:
				background-color 150ms ease,
				color 150ms ease,
				outline 150ms ease;

			&:hover,
			&[data-active="true"] {
				color: $bg-color;
				fill: $bg-color;
				background-color: $accent-color-hover;
			}

			&[data-active="true"] {
				outline: 3px solid $accent-color;
			}

			@media screen and (max-width: 500px) {
				flex-flow: column;
			}
		}

		&.disabled {
			button {
				opacity: 0.8;
				pointer-events: none;
			}
		}

		@media screen and (max-width: 500px) {
			width: 100%;
			justify-content: center;
		}
	}
</style>
