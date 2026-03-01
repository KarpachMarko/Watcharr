<script lang="ts">
	interface Props {
		src: string;
	}

	let { src }: Props = $props();

	let backdropLoaded = $state(false);
</script>

<div class="backdrop">
	<img
		{src}
		alt=""
		class:loaded={backdropLoaded}
		onload={() => (backdropLoaded = true)}
	/>
	<div class="mask"></div>
</div>

<style lang="scss">
	div.backdrop {
		$w: 1920px;
		$h: 800px;

		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		max-width: $w;
		height: $h;
		left: 50%;
		transform: translateX(-50%);
		overflow: hidden;
		pointer-events: none;
		z-index: -1;

		img {
			position: absolute;
			left: 50%;
			transform: translateX(-50%);
			opacity: 0;
			transition: opacity 150ms ease-in;
			/* We force the width and height so if we have a very small image
			to display, it stretches. */
			width: $w;
			height: $h;
			object-fit: cover;

			&.loaded {
				opacity: 1;
			}
		}

		.mask {
			display: block;
			position: absolute;
			width: 100%;
			height: 100%;
			top: 0;
			left: 0;
			background-repeat: no-repeat;
			/* I stole this gradient from letterboxd!!! IM SO SORRY IT JUST LOOKED TOO AMAZING WAAAAAGHGHHH */
			background-image: $blended-mask;
		}
	}
</style>
