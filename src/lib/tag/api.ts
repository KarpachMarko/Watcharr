import axios from "axios";
import { notify } from "../util/notify";
import type { Tag } from "@/types";

export async function tagWatched(
	watchedId: number,
	tag: Tag,
): Promise<boolean> {
	const nid = notify({ text: `Tagging`, type: "loading" });
	if (typeof watchedId !== "number") {
		notify({
			id: nid,
			text: "Failed To Tag! Watched entry not found.",
			type: "error",
		});
		return false;
	}
	return await axios
		.post(`/watched/${watchedId}/tag/${tag.id}`)
		.then((resp) => {
			console.log("tagWatched: Status:", resp.status);
			notify({ id: nid, text: `Tagged!`, type: "success" });
			return true;
		})
		.catch((err) => {
			console.error("tagWatched: Request failed!", err);
			notify({ id: nid, text: "Failed To Tag!", type: "error" });
			return false;
		});
}

export async function untagWatched(
	watchedId: number,
	tag: Tag,
): Promise<boolean> {
	const nid = notify({ text: `Untagging`, type: "loading" });
	if (typeof watchedId !== "number") {
		notify({
			id: nid,
			text: "Failed To Untag! Watched entry not found.",
			type: "error",
		});
		return false;
	}
	return await axios
		.delete(`/watched/${watchedId}/tag/${tag.id}`)
		.then((resp) => {
			console.log("untagWatched: Status:", resp.status);
			notify({ id: nid, text: `Untagged!`, type: "success" });
			return true;
		})
		.catch((err) => {
			console.error("untagWatched: Request failed!", err);
			notify({ id: nid, text: "Failed To Untag!", type: "error" });
			return false;
		});
}
