<script lang="ts">
  import type { PingResult } from "./PingResult";
  import Icon from "mdi-svelte";
  import { mdiLoading, mdiCheck, mdiClose } from "@mdi/js";
  export let row: { ping: PingResult };
  export let column: any;

  function formatTime(time: number) {
    if (time <= 9) {
      return time.toFixed(2);
    }
    if (time < 99) {
      return time.toFixed(1);
    }
    return Math.floor(time).toFixed(0);
  }
</script>

<div>
  {#if row.ping === null}
    <Icon path={mdiLoading} spin />
  {:else if row.ping.time != null}
    <span class="ping-ok">
      <Icon path={mdiCheck} />
      <span>{formatTime(row.ping.time)} ms</span>
    </span>
  {:else}
    <span class="ping-error">
      <Icon path={mdiClose} />
      <span>{row.ping.error}</span>
    </span>
  {/if}
</div>

<style>
  .ping-error {
    color: red;
  }
  .ping-ok {
    color: green;
  }
</style>
