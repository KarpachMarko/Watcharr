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
			background-image: linear-gradient(
					to right,
					white,
					rgba(255, 255, 255, 0.984) 0.97%,
					rgba(255, 255, 255, 0.945) 2.07833333%,
					rgba(255, 255, 255, 0.882) 3.29666667%,
					rgba(255, 255, 255, 0.804) 4.60166667%,
					rgba(255, 255, 255, 0.71) 5.96666667%,
					rgba(255, 255, 255, 0.608),
					rgba(255, 255, 255, 0.5),
					rgba(255, 255, 255, 0.396) 10.16%,
					rgba(255, 255, 255, 0.294) 11.505%,
					rgba(255, 255, 255, 0.204) 12.78%,
					rgba(255, 255, 255, 0.12) 13.95833333%,
					rgba(255, 255, 255, 0.06) 15.01666667%,
					rgba(255, 255, 255, 0.016) 15.92833333%,
					rgba(255, 255, 255, 0) 16.66666667%,
					rgba(255, 255, 255, 0) 83.33333333%,
					rgba(255, 255, 255, 0.016) 84.07166667%,
					rgba(255, 255, 255, 0.06) 84.98333333%,
					rgba(255, 255, 255, 0.12) 86.04166667%,
					rgba(255, 255, 255, 0.204) 87.22%,
					rgba(255, 255, 255, 0.294) 88.495%,
					rgba(255, 255, 255, 0.396) 89.84%,
					rgba(255, 255, 255, 0.5),
					rgba(255, 255, 255, 0.608),
					rgba(255, 255, 255, 0.71) 94.03333333%,
					rgba(255, 255, 255, 0.804) 95.39833333%,
					rgba(255, 255, 255, 0.882) 96.70333333%,
					rgba(255, 255, 255, 0.945) 97.92166667%,
					rgba(255, 255, 255, 0.984) 99.03%,
					white
				),
				linear-gradient(
					to top,
					white,
					white 21.48148148%,
					rgba(255, 255, 255, 0.984) 23.63703704%,
					rgba(255, 255, 255, 0.945) 26.1%,
					rgba(255, 255, 255, 0.882) 28.80740741%,
					rgba(255, 255, 255, 0.804) 31.70740741%,
					rgba(255, 255, 255, 0.71) 34.74074074%,
					rgba(255, 255, 255, 0.608),
					rgba(255, 255, 255, 0.5) 40.97407407%,
					rgba(255, 255, 255, 0.396) 44.05925926%,
					rgba(255, 255, 255, 0.294) 47.04814815%,
					rgba(255, 255, 255, 0.204) 49.88148148%,
					rgba(255, 255, 255, 0.12) 52.5%,
					rgba(255, 255, 255, 0.06) 54.85185185%,
					rgba(255, 255, 255, 0.016) 56.87777778%,
					rgba(255, 255, 255, 0) 58.51851852%
				);
		}
	}
</style>
