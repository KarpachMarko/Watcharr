<script lang="ts">
	interface Props {
		title: string;
		desc: string;
		type: "warn";
	}

	let { title, desc, type }: Props = $props();
</script>

<div class={type}>
	<h4 class="norm">{title}</h4>
	<p>{desc}</p>
</div>

<style lang="scss">
	div {
		padding: 10px 12px;
		border-radius: 10px;

		&.warn {
			background-color: $warn;
		}
	}
</style>
