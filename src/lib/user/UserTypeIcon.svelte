<script lang="ts">
  import { UserType } from "@/types";
  import Icon from "../Icon.svelte";

  export let type: UserType;
</script>

{#if !type}
  <span
    style="font-family: 'Rampart One'; font-weight: bold; font-size: 21px; line-height: 20px; user-select: none;"
  >
    W
  </span>
{:else if type == UserType.Jellyfin}
  <Icon i="jellyfin" wh={20} />
{:else if type == UserType.Plex}
  <Icon i="plex" wh={20} />
{:else if type == UserType.Proxy}
  <Icon i="lock-closed" wh={20} />
{:else}
  <!-- Unknown type -->
  <span
    style="font-family: 'Rampart One'; font-weight: bold; font-size: 21px; line-height: 20px; user-select: none;"
  >
    ?
  </span>
{/if}
