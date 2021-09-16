<script lang="ts">
  import type IDevice from "./model/IDevice";
  import Icon from "mdi-svelte";
  import { mdiRouterNetwork, mdiAccessPoint } from "@mdi/js";
 
  export let device: IDevice;
  import { onDestroy } from "svelte";
  import { formatDuration } from "./util";
  let uptimeSeconds = device.uptimeSeconds;

  const interval = setInterval(() => (uptimeSeconds += 1), 1000);

  onDestroy(() => clearInterval(interval));

  function addrWithStrippedPort() {
    return device.address.split(":")[0];
  }
</script>

<div class="card">
  <div class="header">
    <div class="icon">
      <Icon path={device.hasDHCP ? mdiRouterNetwork : mdiAccessPoint} />
    </div>
    <div class="name">{device.hostname}</div>
  </div>
  <div class="spacer" />
  <table class="details">
    {#if device.vendor}
      <tr>
        <td class="detail-name"> vendor </td>
        <td>
          <div class="vendor">
            {device.vendor}
          </div>
        </td>
      </tr>
    {/if}
    <tr>
      <td class="detail-name"> uptime </td>
      <td>
        {formatDuration(uptimeSeconds)}
      </td>
    </tr>
    <tr>
      <td class="detail-name"> address </td>
      <td>
        <a href="http://{addrWithStrippedPort()}">{addrWithStrippedPort()}</a>
      </td>
    </tr>
  </table>
</div>

<style>
  .card {
    height: 100px;
    width: 200px;
    background: var(--background-dim);
    margin-right: 16px;
    padding: 8px;
    display: flex;
    flex-direction: column;
  }
  .card .name {
    font-weight: bold;
    font-size: 1.1rem;
    padding-top: 1px;
  }
  .card .header {
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    align-items: center;
  }
  .spacer {
    flex-grow: 1;
  }
  .icon {
    margin-right: 4px;
  }
  .details {
    width: 100%;
  }
  td {
    text-align: right;
  }
  .detail-name {
    padding-right: 4px;
    /* text-align: right; */
  }
  .vendor {
    text-overflow: ellipsis;
    white-space: nowrap;
    width: 130px;
    overflow: hidden;
    text-align: right;
  }
</style>
