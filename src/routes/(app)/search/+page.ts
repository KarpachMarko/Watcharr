export async function load({ url }) {
	return {
		query: url.searchParams.get("query"),
	};
}
