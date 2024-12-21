<script lang="ts">
  import Modal from "@/lib/Modal.svelte";
  import Spinner from "@/lib/Spinner.svelte";
  import Setting from "@/lib/settings/Setting.svelte";
  import SettingsList from "@/lib/settings/SettingsList.svelte";
  import { toRelativeTime } from "@/lib/util/helpers";
  import { notify } from "@/lib/util/notify";
  import type { AllTasksResponse } from "@/types";
  import axios from "axios";
  import { onMount } from "svelte";

  export let onClose: () => void;

  let formDisabled = false;
</script>

<Modal
  title="Tasks Schedule"
  desc="Want a routine task to occur more or less frequently? Configure it below."
  {onClose}
>
  <SettingsList>
    <Setting
      title="Trusted Header Authentication"
      desc="Name of the authentication header for proxy authentication. Only set this if Watcharr is running behind a trusted proxy"
    >
      <!-- TODO make this a modal like tasks schedule, we will have a few related settings that can be grouped into the modal -->
      <input
        type="text"
        placeholder="X-User"
        on:blur={() => {
          proxyHeaderDisabled = true;
          updateServerConfig("PROXY_AUTH_HEADER", serverConfig.PROXY_AUTH_HEADER, () => {
            proxyHeaderDisabled = false;
          });
        }}
        disabled={proxyHeaderDisabled}
        bind:value={serverConfig.PROXY_AUTH_HEADER}
      />
    </Setting>
  </SettingsList>
</Modal>

<style lang="scss">
  input {
    width: 200px;
    margin-top: 5px;
  }
</style>
