<script lang="ts">
	interface Props {
		pretty: string;
		error: any;
		onRetry?: () => void | undefined;
	}

	let { pretty, error, onRetry = undefined! }: Props = $props();
</script>

<div>
	<div>
		<strong>{pretty}</strong>
		{#if error?.message}
			<p>{error.message}</p>
			{#if error.response?.data?.error}
				<p>{error.response.data.error}</p>
			{/if}
		{:else}
			<p>{JSON.stringify(error)}</p>
		{/if}
		{#if onRetry}
			<button onclick={onRetry}>Try Again</button>
		{/if}
	</div>
</div>

<style lang="scss">
	div {
		display: flex;
		justify-content: center;
		width: 100%;

		div {
			display: flex;
			justify-content: center;
			flex-flow: column;
			max-width: 500px;
			gap: 5px;
			text-transform: capitalize;
			font-size: 16px;

			p {
				font-size: 12px;
			}
		}
	}
</style>
