export async function load({ url }) {
	return {
		query: url.searchParams.get("query"),
		type: url.searchParams.get("type"),
	};
}
