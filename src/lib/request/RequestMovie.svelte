<script lang="ts">
	import { run } from "svelte/legacy";

	import axios from "axios";
	import Modal from "../Modal.svelte";
	import type {
		ArrRequestResponse,
		Media,
		RadarrSettings,
		RadarrTestResponse,
	} from "@/types";
	import { notify } from "../util/notify";
	import DropDown from "../DropDown.svelte";
	import Setting from "../settings/Setting.svelte";
	import Spinner from "../Spinner.svelte";

	interface Props {
		content: Media;
		onClose: (r: ArrRequestResponse | undefined) => void;
		approveMode?: boolean;
		originalRequest?: ArrRequestResponse | undefined;
	}

	let {
		content,
		onClose,
		approveMode = false,
		originalRequest = undefined,
	}: Props = $props();

	let servarrs: RadarrSettings[] = $state();
	let selectedServarrIndex: number = $state();
	let inputsDisabled = true;
	let selectedServerCfg: RadarrTestResponse | undefined = $state();
	let addRequestRunning = $state(false);

	async function getServers() {
		try {
			inputsDisabled = true;
			const r = await axios.get("/arr/rad");
			if (r.data?.length > 0) {
				servarrs = r.data;
				selectedServarrIndex = 0;
			} else {
				notify({ text: "No servers found", type: "error" });
			}
			inputsDisabled = false;
			processOriginalRequest();
		} catch (err) {
			console.error("Failed to get servers!", err);
			notify({ text: "Failed to load servers", type: "error" });
		}
	}

	async function getConfig(name: string) {
		try {
			inputsDisabled = true;
			const r = await axios.get<RadarrTestResponse>(`/arr/rad/config/${name}`);
			selectedServerCfg = r.data;
			inputsDisabled = false;
		} catch (err) {
			console.error("Failed to get config!", err);
			notify({ text: "Failed to load config", type: "error" });
		}
	}

	async function request() {
		let nid;
		try {
			if (!servarrs || !servarrs[selectedServarrIndex]) {
				notify({ text: "Must select a server", type: "error" });
				return;
			}
			if (!selectedServerCfg) {
				notify({ text: "No selected server config found", type: "error" });
				return;
			}
			addRequestRunning = true;
			nid = notify({ text: "Requesting", type: "loading" });
			const server = servarrs[selectedServarrIndex];
			const rootFolder = selectedServerCfg.rootFolders?.find(
				(f) => f.id === server.rootFolder,
			);
			if (!rootFolder) {
				console.error(
					"movie request.. no root folder found with id:",
					server.rootFolder,
					"rf:",
					rootFolder,
				);
				notify({ id: nid, text: "No Root Folder Found", type: "error" });
				return;
			}
			const resp = await axios.post<ArrRequestResponse>(
				`/arr/rad/request${approveMode && originalRequest ? `/approve/${originalRequest.id}` : ""}`,
				{
					serverName: server.name,
					title: content.name,
					year: content.releaseDate
						? new Date(content.releaseDate)?.getFullYear()
						: undefined,
					tmdbId: content.ids.tmdb,
					qualityProfile: server.qualityProfile,
					rootFolder: rootFolder.path,
				},
			);
			addRequestRunning = false;
			if (resp.data) {
				notify({ id: nid, text: "Request complete", type: "success" });
				onClose(resp.data);
			}
		} catch (err) {
			console.error("content request failed!", err);
			addRequestRunning = false;
			notify({ id: nid, text: "Request failed!", type: "error" });
		}
	}

	function processOriginalRequest() {
		if (!originalRequest) {
			return;
		}
		try {
			if (originalRequest.requestJson) {
				const ogr = JSON.parse(originalRequest.requestJson);
				if (!ogr) {
					console.info("processOriginalRequest: No json.", ogr);
					return;
				}
				if (ogr?.serverName) {
					console.debug(
						"processOriginalRequest: restoring server name:",
						ogr?.serverName,
					);
					const idx = servarrs?.findIndex((s) => s.name === ogr?.serverName);
					if (idx !== -1) {
						selectedServarrIndex = idx;
					}
				}
			} else {
				notify({
					type: "error",
					text: "Full original request could not be restored. You may continue, but prefilled settings may not be true to the original request.",
					time: 10000,
				});
			}
		} catch (err) {
			console.error("processOriginalRequest: Failed!", err);
			notify({
				text: "Failed when processing original request!",
				type: "error",
			});
		}
	}

	run(() => {
		if (typeof selectedServarrIndex !== "undefined" && servarrs?.length > 0) {
			const s = servarrs[selectedServarrIndex];
			if (!s) {
				selectedServerCfg = undefined;
			} else {
				getConfig(s.name);
			}
		}
	});

	getServers();
</script>

<Modal
	title={approveMode ? "Approve Request" : "Request"}
	desc={content.name}
	onClose={() => onClose(undefined)}
>
	<div class="req-ctr">
		{#if servarrs}
			{@const server = servarrs[selectedServarrIndex]}

			{#if servarrs?.length > 1}
				<Setting title="Select the server to use">
					<DropDown
						placeholder="Select a server"
						bind:active={selectedServarrIndex}
						options={servarrs?.length > 0
							? servarrs.map((s, i) => {
									return { id: i, value: s.name };
								})
							: []}
					/>
				</Setting>
			{/if}

			<button onclick={request} disabled={addRequestRunning}>
				{approveMode ? "Approve" : "Request"}
			</button>
		{:else}
			<Spinner />
		{/if}
	</div>
</Modal>

<style lang="scss">
	.req-ctr {
		display: flex;
		flex-flow: column;
		gap: 10px;
		height: 100%;

		button {
			margin-top: auto;
			margin-left: auto;
			width: max-content;
		}
	}
</style>
