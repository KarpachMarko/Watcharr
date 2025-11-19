<script lang="ts">
	import { store } from "@/store.svelte";
	import Menu from "../Menu.svelte";

	function sortClicked(type: string) {
		let mode = "UP";
		// If this sort is already the `activeSort`
		if (store.activeSort[0] == type) {
			if (store.activeSort[1] === "UP") {
				mode = "DOWN";
			} else if (store.activeSort[1] === "DOWN") {
				mode = "";
			}
		}
		store.activeSort = [type, mode];
	}
</script>

<Menu conf={{ width: "180px", right: "90px", arrowLeft: "21px" }}>
	<button
		class={`plain ${store.activeSort[0] == "DATEADDED" ? store.activeSort[1].toLowerCase() : ""}`}
		onclick={() => sortClicked("DATEADDED")}
	>
		Date Added
	</button>
	<button
		class={`plain ${store.activeSort[0] == "LASTCHANGED" ? store.activeSort[1].toLowerCase() : ""}`}
		onclick={() => sortClicked("LASTCHANGED")}
	>
		Last Changed
	</button>
	<button
		class={`plain ${store.activeSort[0] == "LASTFIN" ? store.activeSort[1].toLowerCase() : ""}`}
		onclick={() => sortClicked("LASTFIN")}
	>
		Last Finished
	</button>
	<button
		class={`plain ${store.activeSort[0] == "RATING" ? store.activeSort[1].toLowerCase() : ""}`}
		onclick={() => sortClicked("RATING")}
	>
		Rating
	</button>
	<button
		class={`plain ${store.activeSort[0] == "ALPHA" ? store.activeSort[1].toLowerCase() : ""}`}
		onclick={() => sortClicked("ALPHA")}
	>
		Alphabetical
	</button>
</Menu>

<style lang="scss">
	button {
		position: relative;

		&.down::before {
			content: "\2193";
		}

		&.up::before {
			content: "\2191";
		}

		&.on::before {
			content: "\2713";
		}

		&::before {
			position: absolute;
			top: 4px;
			left: 12px;
			font-family:
				system-ui,
				-apple-system,
				BlinkMacSystemFont;
			font-size: 18px;
		}
	}
</style>
