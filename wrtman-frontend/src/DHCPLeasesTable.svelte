<script lang="ts">
  import config from "./config";
  import type ColumnModel from "./data-table/ColumnModel";
  import DataTable from "./data-table/DataTable.svelte";
  import DurationUntillCell from "./data-table/DurationUntillCell.svelte";
  import SimpleTextCell from "./data-table/SimpleTextCell.svelte";
  import { sortDates, sortText } from "./data-table/sortFuncs";
  import type IDHCPLease from "./model/IDHCPLease";
  import VendorLogo from "./VendorLogo.svelte";

  type MappedDHCPLease = IDHCPLease & {
    expiryTime: Date;
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
      key: "expiryTime",
      label: "Expiry time",
      component: [DurationUntillCell],
      sortFunc: sortDates,
    },
  ];

  const fetchDHCPLeases = fetch(config.baseURL + "/api/dhcp-leases")
    .then((res) => res.json() as Promise<IDHCPLease[]>)
    .then((data) =>
      data
        .map((l) => ({ ...l, expiryTime: new Date(l.expiryTime) }))
        .sort((a, b) => b.expiryTime.getTime() - a.expiryTime.getTime())
    );
</script>

<section>
  <h3>DHCP Leases</h3>
  {#await fetchDHCPLeases}
    <p>...loading</p>
  {:then data}
    <p>{data.length} devices connected.</p>
  {:catch error}
    <p>An error occurred!</p>
  {/await}
  <DataTable {columns} dataPromise={fetchDHCPLeases} defaultSortKey="-expiryTime" />
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
