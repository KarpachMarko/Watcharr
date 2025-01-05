import { store } from "@/store.svelte";

export interface Notification {
	/**
	 * Notification ID.
	 * Used to reference an exiting notification.
	 */
	id?: number;

	/**
	 * Text shown in popup;
	 */
	text: string;

	/**
	 * Type of notification, controls the style.
	 * Loading notifs never hide, so after request completion
	 * or failure, the notif must be updated or `unNotify`ed.
	 */
	type?: "error" | "success" | "loading";

	/**
	 * How long in milliseconds the popup will stay shown for.
	 */
	time?: number;
}

export function notify(n: Notification) {
	if (n.id) {
		const notif = store.notifications.find((not) => not.id === n.id);
		if (notif) {
			notif.type = n.type;
			notif.text = n.text;
		} else {
			console.error("Can't update notif that doesnt exist", n);
		}
	} else {
		n.id = Math.random();
		store.notifications.push({ ...n });
	}
	if (n.type !== "loading" && n.time !== Infinity) {
		setTimeout(() => unNotify(n.id!), n.time ?? 2500);
	}
	return n.id;
}

export function unNotify(id?: number) {
	if (!id) {
		console.warn("unNotify: Tried removing a notification without an id.");
		return;
	}
	store.notifications = store.notifications.filter((e) => e.id !== id);
}
