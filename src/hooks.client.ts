import type { ClientInit } from "@sveltejs/kit";

export const init: ClientInit = async () => {
	console.info(
		`%cWATCHARR v${__WATCHARR_VERSION__}`,
		"background: white;color: black;font-size: large;padding: 3px 5px;",
	);
};
