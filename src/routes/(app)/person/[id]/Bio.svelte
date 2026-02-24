<script lang="ts">
	interface Props {
		bio?: string;
	}

	const { bio }: Props = $props();

	const splitBio:
		| {
				preview: string;
				rest: string;
		  }
		| undefined = $derived.by(() => {
		if (!bio) {
			return;
		}
		let index = bio.indexOf("\n");
		return {
			// The preview we show.
			preview: bio.slice(0, index),
			rest: bio.slice(index + 2),
		};
	});
	let bioExtended = $state(false);
</script>

{#if splitBio}
	<span>Biography</span>
	<p>
		{splitBio.preview}
		{#if bioExtended}
			<br /><br />{splitBio.rest}
		{/if}
		{#if splitBio.rest}
			<button class="plain" onclick={() => (bioExtended = !bioExtended)}>
				{bioExtended ? "Less" : "More"}
			</button>
		{/if}
	</p>
{/if}

<style lang="scss">
	span {
		font-weight: bold;
		font-size: 14px;
	}

	p {
		font-size: 14px;
		white-space: pre-line;
	}

	button {
		font-weight: bold;
		color: white;
	}
</style>
