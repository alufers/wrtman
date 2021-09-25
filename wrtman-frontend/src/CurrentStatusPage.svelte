<script lang="ts">
  import config from "./config";
  import DHCPLeasesTable from "./DHCPLeasesTable.svelte";

  import type IDevice from "./model/IDevice";
  import NetworkDeviceCard from "./NetworkDeviceCard.svelte";

  const fetchNetworkDevices = fetch(config.baseURL + "/api/devices").then(
    (res) => res.json() as Promise<IDevice[]>
  );
</script>

<div>
  <section>
    <h3>Network devices</h3>
    {#await fetchNetworkDevices}
      <p>...loading</p>
    {:then data}
      <div class="cards-row">
        {#each data as device}
          <NetworkDeviceCard {device} />
        {/each}
      </div>
    {:catch error}
      <p>An error occurred!</p>
    {/await}
  </section>
  <DHCPLeasesTable />
</div>

<style>
  .cards-row {
    display: flex;
    flex-direction: row;
  }
</style>
