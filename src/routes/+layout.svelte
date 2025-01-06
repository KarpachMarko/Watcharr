<script lang="ts">
	import Icon from "@/lib/Icon.svelte";
	import SpinnerTiny from "@/lib/SpinnerTiny.svelte";
	import { unNotify } from "@/lib/util/notify";
	import { startStoreSaver, store } from "@/store.svelte";
	import { onMount } from "svelte";
	import { pwaInfo } from "virtual:pwa-info";

	interface Props {
		children?: import("svelte").Snippet;
	}

	let { children }: Props = $props();

	console.log(
		`%cWATCHARR v${__WATCHARR_VERSION__}`,
		"background: white;color: black;font-size: large;padding: 3px 5px;",
	);

	function resetTooltipPos() {
		const t = document.getElementById("tooltip");
		if (t) {
			t.style.top = "0";
			t.style.left = "0";
		}
	}

	onMount(() => {
		startStoreSaver();
		window.addEventListener("resize", resetTooltipPos);

		return () => {
			window.removeEventListener("resize", resetTooltipPos);
		};
	});
</script>

<svelte:head>
	{#if pwaInfo?.webManifest?.linkTag}
		<!-- eslint-disable-next-line -->
		{@html pwaInfo.webManifest.linkTag}
	{/if}
</svelte:head>

<div id="tooltip"></div>
<div id="notifications">
	{#each store.notifications as n}
		<div class={`${n.type} notif`}>
			{#if n.type === "loading"}
				<SpinnerTiny />
			{/if}
			<!-- only comes from our strings (which may have html) -->
			<!-- eslint-disable-next-line -->
			<span>{@html n.text}</span>
			<button
				class="plain"
				onclick={() => {
					unNotify(n.id);
				}}
			>
				<Icon i="close" />
			</button>
		</div>
	{/each}
</div>

{@render children?.()}
