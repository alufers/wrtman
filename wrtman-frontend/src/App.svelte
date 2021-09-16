<script lang="ts">
  import config from "./config";
  import DHCPLeasesTable from "./DHCPLeasesTable.svelte";

  import Logo from "./Logo.svelte";
  import type IDevice from "./model/IDevice";
  import NetworkDeviceCard from "./NetworkDeviceCard.svelte";

  const fetchNetworkDevices = fetch(config.baseURL + "/api/devices").then(
    (res) => res.json() as Promise<IDevice[]>
  );
</script>

<main>
  <Logo />
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
  <footer>Made by alufers.</footer>
</main>

<style>
  :global(body) {
    background: black;
    color: #ddd;
    font-family: monospace;
    --primary: lightseagreen;
    --secondary: lightpink;
    --background-dim: #333;
    --background-dimmer: #555;
    --foreground-secondary: #bbb;
  }
  :global(a) {
    color: var(--primary);
  }
  :global(input) {
    border: none;
    outline: none;
    background: transparent;
    color: #fff;
    border: 1px solid var(--foreground-secondary);
    padding: 4px;
    font-family: monospace;
  }
  main {
    padding: 0.5em;
    padding-left: 5em;
    padding-right: 5em;
    margin: 0 auto;
  }

  @media (max-width: 640px) {
    main {
      padding-left: 0.5em;
      padding-right: 0.5em;
      max-width: none;
    }
  }

  .cards-row {
    display: flex;
    flex-direction: row;
  }
  footer {
    margin-top: 128px;
    color: var(--foreground-secondary);
  }

  :global(button) {
    outline: none;
    background: transparent;
    color: var(--primary);
    border: 1px solid var(--primary);
    font-family: monospace;
    cursor: pointer;
    padding: 4px;
  }
  :global(button):hover {
    background: var(--primary);
    color: #ddd;
  }
  :global(button):active {
    transform: scale(0.9);
  }
</style>
