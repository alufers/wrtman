<script lang="ts">
  import config from "./config";
  import type ColumnModel from "./data-table/ColumnModel";
  import DataTable from "./data-table/DataTable.svelte";
  import DurationUntillCell from "./data-table/DurationUntillCell.svelte";
  import SimpleTextCell from "./data-table/SimpleTextCell.svelte";
  import { sortDates, sortNumbers, sortText } from "./data-table/sortFuncs";
  import lineReader from "./lineReader";
  import type IDHCPLease from "./model/IDHCPLease";
  import PingCell from "./PingCell.svelte";
  import type { PingResult } from "./PingResult";
  import VendorLogo from "./VendorLogo.svelte";

  type MappedDHCPLease = IDHCPLease & {
    expiryTime: Date;
    ping: PingResult | null;
  };

  const columns: ColumnModel<MappedDHCPLease>[] = [
    {
      key: "hostname",
      label: "Hostname",
      filterable: true,
      sortFunc: sortText,
    },
    {
      key: "ipAddress",
      label: "IP address",
      filterable: true,
      sortFunc: sortText,
    },
    {
      key: "vendor",
      label: "Vendor",
      filterable: true,
      component: [VendorLogo, SimpleTextCell],
      sortFunc: sortText,
    },
    {
      key: "macAddress",
      filterable: true,
      label: "MAC address",
      sortFunc: sortText,
    },
    {
      key: "ping",
      label: "Ping",
      component: [PingCell],
    },
    {
      key: "expiryTime",
      label: "Expiry time",
      component: [DurationUntillCell],
      sortFunc: sortDates,
    },
    {
      key: "ssid",
      label: "Wireless Network",
      filterable: true,
      sortFunc: sortText,
    },
    {
      key: "signalStrength",
      label: "RSSI",
      sortFunc: sortNumbers,
    },
    {
      key: "apHostname",
      label: "Access Point",
      filterable: true,
      sortFunc: sortText,
    },
  ];

  let fetchDHCPLeases = fetch(config.baseURL + "/api/dhcp-leases")
    .then((res) => res.json() as Promise<IDHCPLease[]>)
    .then((data) =>
      data
        .map((l) => ({ ...l, expiryTime: new Date(l.expiryTime), ping: null }))
        .sort((a, b) => b.expiryTime.getTime() - a.expiryTime.getTime())
    );

  const doPing = async (timeout?: number) => {
    const leases = await fetchDHCPLeases;
    const resp = await fetch("/api/ping", {
      method: "POST",
      headers: {
        "Content-type": "application/json",
      },
      body: JSON.stringify({
        addresses: leases.filter((l) => !l?.ping?.time).map((l) => l.ipAddress),
        timeout,
      }),
    });
    for await (const line of lineReader(resp.body.getReader())) {
      const res = JSON.parse(line);
      fetchDHCPLeases = fetchDHCPLeases.then((leases) =>
        leases.map((l) => {
          if (l.ipAddress === res.address) {
            return { ...l, ping: res };
          }
          return l;
        })
      );
      console.log(res);
    }
  };
  const pingPromise = (async () => {
    await doPing(0.5);
    await doPing(2);
    await doPing(5);
  })();
</script>

<section>
  <h3>DHCP Leases</h3>
  <!-- svelte-ignore empty-block -->
  <!-- {#await pingPromise}{:catch error}
    <p>Failed to ping: {error.message}</p>
  {/await} -->
  {#await fetchDHCPLeases}
    <p>...loading</p>
  {:then data}
    <p>{data.length} devices connected.</p>
  {:catch error}
    <p>An error occurred!</p>
  {/await}
  <DataTable
    {columns}
    dataPromise={fetchDHCPLeases}
    defaultSortKey="-expiryTime"
  />
</section>

<style>
  table {
    border-collapse: collapse;
  }
  td {
    padding-left: 8px;
    padding-right: 8px;
    border-bottom: 1px solid var(--background-dim);
  }
  .expiring {
    color: #aaa;
  }
</style>
