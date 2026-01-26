import type { PaginationResponse } from "@/types";
import type { AxiosResponse, GenericAbortSignal } from "axios";

export default function paginatedLoader<T>(
	fn: (
		sig: GenericAbortSignal,
	) => Promise<AxiosResponse<PaginationResponse<T>, any> | undefined>,
) {
	let reqController = new AbortController();

	const state: {
		data: T[];
		page: number;
		pageMax: number;
		reqLoading: boolean;
		reqLoadError: Error | undefined;
	} = $state({
		data: [],
		page: 0,
		pageMax: 1,
		reqLoading: false,
		reqLoadError: undefined,
	});

	/**
	 * Resets state.
	 */
	const reset = () => {
		state.data = [];
		state.page = 0;
		state.pageMax = 1;
		state.reqLoading = false;
		state.reqLoadError = undefined;
	};

	/**
	 * Abort the request.
	 */
	const abortReq = (reason: string) => {
		reqController.abort(reason);
	};

	/**
	 * Runs the paginated request passed in as `fn` with our
	 * paginated logic wrapped around.
	 */
	const runFn = async () => {
		const logStyle = "font-weight: bold; font-size: 18px;";
		if (state.reqLoading) {
			console.warn("%cpaginatedLoader->runFn: already running", logStyle);
			return;
		}
		if (state.page >= state.pageMax) {
			console.warn("%cpaginatedLoader->runFn: max page reached", logStyle);
			return;
		}

		state.reqLoading = true;
		reqController = new AbortController();
		try {
			const resp = await fn(reqController.signal);
			if (!resp) {
				state.reqLoading = false;
				console.error(
					"paginatedLoader->runFn: fn returned nothing! May not have wanted to for some reason..",
				);
				return;
			}
			state.page = resp.data.page;
			state.pageMax = resp.data.totalPages;
			console.debug(
				`%cpaginatedLoader->runFn: Loaded Page=${state.page} Max=${state.pageMax}`,
				logStyle,
			);
			if (resp.data.results.length <= 0) {
				state.reqLoading = false;
				console.warn("loadWatchedList: No results.");
				return;
			}
			state.data.push(...resp.data.results);
			state.data = state.data;
		} catch (err: any) {
			if (err?.code === "ERR_CANCELED") {
				console.warn("loadWatchedList: Cancelled, not showing error.");
			} else {
				console.error("loadWatchedList: failed!", err);
				state.reqLoadError = err;
			}
		}
		state.reqLoading = false;
	};

	return {
		runFn,
		reset,
		abortReq,
		state,
	};
}
