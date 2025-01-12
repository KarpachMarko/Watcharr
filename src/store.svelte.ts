import type {
	Filters,
	Follow,
	ImportedList,
	PrivateUser,
	ServerFeatures,
	Tag,
	Theme,
	UserSettings,
	WLDetailedViewOption,
	Watched,
} from "./types";
import type { Notification } from "./lib/util/notify";
import { browser } from "$app/environment";
import { toggleTheme } from "./lib/util/helpers";

export const defaultSort = ["DATEADDED", "DOWN"];

interface Store {
	userInfo: PrivateUser | undefined;
	userSettings: UserSettings | undefined;
	watchedList: Watched[];
	notifications: Notification[];
	activeSort: string[];
	activeFilters: Filters;
	appTheme: Theme;
	importedList:
		| {
				data: string;
				type:
					| "text-list"
					| "tmdb"
					| "movary"
					| "watcharr"
					| "myanimelist"
					| "ryot"
					| "todomovies"
					| "imdb";
		  }
		| undefined;
	parsedImportedList: ImportedList[] | undefined;
	searchQuery: string;
	serverFeatures: ServerFeatures | undefined;
	follows: Follow[];
	wlDetailedView: WLDetailedViewOption[];
	tags: Tag[];
}

export const store: Store = $state({
	watchedList: [],
	notifications: [],
	activeSort: defaultSort,
	activeFilters: { type: [], status: [] },
	appTheme: "light",
	importedList: undefined,
	parsedImportedList: undefined,
	searchQuery: "",
	userInfo: undefined,
	userSettings: undefined,
	serverFeatures: undefined,
	follows: [],
	wlDetailedView: [],
	tags: [],
});

/**
 * Reset everything in `store` back to default values.
 */
export const clearAllStores = () => {
	store.watchedList = [];
	store.notifications = [];
	store.activeSort = defaultSort;
	store.appTheme = "light";
	store.importedList = undefined;
	store.parsedImportedList = undefined;
	store.searchQuery = "";
	store.userInfo = undefined;
	store.userSettings = undefined;
	store.serverFeatures = undefined;
	store.follows = [];
	store.wlDetailedView = [];
	store.tags = [];
	clearActiveFilters();
};

export const clearActiveFilters = () => {
	store.activeFilters = { type: [], status: [] };
};

if (browser) {
	rehydrateStore();
}

/**
 * Restore state from localStorage and apply values into
 * our `store`.
 */
function rehydrateStore() {
	console.info("rehydrateStore: Running..");
	const raf = localStorage.getItem("activeFilter");
	if (raf) {
		store.activeSort = JSON.parse(raf);
		console.debug(
			"rehydrateStore: Restored activeSort:",
			$state.snapshot(store.activeSort),
		);
	}

	const filters = localStorage.getItem("activeFilterReal");
	if (filters) {
		store.activeFilters = JSON.parse(filters);
		console.debug(
			"rehydrateStore: Restored activeFilters:",
			$state.snapshot(store.activeFilters),
		);
	}

	const theme = localStorage.getItem("theme") as Theme;
	if (theme) {
		store.appTheme = theme;
		toggleTheme(theme);
		console.debug(
			"rehydrateStore: Restored appTheme:",
			$state.snapshot(store.appTheme),
		);
	} else {
		let defTheme: Theme = "light";
		if (window.matchMedia("(prefers-color-scheme: dark)").matches) {
			defTheme = "dark";
		}
		console.log(
			"Theme not set, setting default theme from system theme:",
			defTheme,
		);
		store.appTheme = defTheme;
		toggleTheme(defTheme);
		console.debug(
			"rehydrateStore: appTheme hydrated from system default:",
			$state.snapshot(store.appTheme),
		);
	}

	const wlDetailedViewR = localStorage.getItem("wlDetailedView");
	if (wlDetailedViewR) {
		store.wlDetailedView = JSON.parse(wlDetailedViewR);
		console.debug(
			"rehydrateStore: Restored wlDetailedView:",
			$state.snapshot(store.wlDetailedView),
		);
	}
}

/**
 * Start tracking changes for state we want to persist
 * in localStorage.
 */
export function startStoreSaver() {
	console.info("startStoreSaver: Creating savers.");

	$effect(() => {
		if (store.activeSort) {
			localStorage.setItem("activeFilter", JSON.stringify(store.activeSort));
			console.debug(
				"StoreSaver: Saved activeSort:",
				localStorage.getItem("activeFilter"),
			);
		}
	});

	$effect(() => {
		if (store.activeFilters) {
			localStorage.setItem(
				"activeFilterReal",
				JSON.stringify(store.activeFilters),
			);
			console.debug(
				"StoreSaver: Saved activeFilterReal:",
				localStorage.getItem("activeFilterReal"),
			);
		}
	});

	$effect(() => {
		if (store.appTheme) {
			localStorage.setItem("theme", store.appTheme);
			console.debug(
				"StoreSaver: Saved appTheme:",
				localStorage.getItem("theme"),
			);
		}
	});

	$effect(() => {
		if (store.wlDetailedView) {
			localStorage.setItem(
				"wlDetailedView",
				JSON.stringify(store.wlDetailedView),
			);
			console.debug(
				"StoreSaver: Saved wlDetailedView:",
				localStorage.getItem("wlDetailedView"),
			);
		} else {
			localStorage.removeItem("wlDetailedView");
			console.debug("StoreSaver: Removed wlDetailedView");
		}
	});
}
