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
	const raf = localStorage.getItem("activeFilter");
	if (raf) {
		store.activeSort = JSON.parse(raf);
	}

	const filters = localStorage.getItem("activeFilterReal");
	if (filters) {
		store.activeFilters = JSON.parse(filters);
	}

	const theme = localStorage.getItem("theme") as Theme;
	if (theme) {
		store.appTheme = theme;
		toggleTheme(theme);
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
	}

	const wlDetailedViewR = localStorage.getItem("wlDetailedView");
	if (wlDetailedViewR) {
		store.wlDetailedView = JSON.parse(wlDetailedViewR);
	}
}

/**
 * Start tracking changes for state we want to persist
 * in localStorage.
 */
export function startStoreSaver() {
	// Save changes
	$effect(() => {
		if (store.activeSort)
			localStorage.setItem("activeFilter", JSON.stringify(store.activeSort));
	});

	$effect(() => {
		if (store.activeFilters)
			localStorage.setItem(
				"activeFilterReal",
				JSON.stringify(store.activeFilters),
			);
	});

	$effect(() => {
		if (store.appTheme) localStorage.setItem("theme", store.appTheme);
	});

	$effect(() => {
		if (store.wlDetailedView) {
			localStorage.setItem(
				"wlDetailedView",
				JSON.stringify(store.wlDetailedView),
			);
		} else {
			localStorage.removeItem("wlDetailedView");
		}
	});
}
