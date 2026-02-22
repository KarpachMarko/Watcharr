<script lang="ts">
	interface Props {
		src: string;
	}

	let { src }: Props = $props();
</script>

<img {src} alt="" />

<style lang="scss">
	img {
		width: 245px;
		height: 100%;
		box-shadow: 0px 0px 14px -4px #000;
		border-radius: 12px;
	}
</style>
