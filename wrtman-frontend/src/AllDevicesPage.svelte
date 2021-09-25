<script lang="ts">
  import config from "./config";
  import type ColumnModel from "./data-table/ColumnModel";
  import DataTable from "./data-table/DataTable.svelte";
  import DurationUntillCell from "./data-table/DurationUntillCell.svelte";
import SimpleTextCell from "./data-table/SimpleTextCell.svelte";
  import { sortDates, sortText } from "./data-table/sortFuncs";
  import type ISavedDevice from "./model/ISavedDevice";
import VendorLogo from "./VendorLogo.svelte";
  type IMappedSavedDevice = ISavedDevice & {
    createdAt: Date;
    updatedAt: Date;
    deletedAt: Date | null;
    lastSeen: Date;
  };
  const fetchAllDevices = fetch(config.baseURL + "/api/all-devices")
    .then((res) => res.json() as Promise<ISavedDevice[]>)
    .then((devices) =>
      devices.map(
        (device) =>
          ({
            ...device,
            createdAt: new Date(device.createdAt),
            deletedAt: device.deletedAt ? new Date(device.deletedAt) : null,
            lastSeen: new Date(device.lastSeen),
            updatedAt: new Date(device.updatedAt),
          } as IMappedSavedDevice)
      )
    );
  const columns: ColumnModel<IMappedSavedDevice>[] = [
    {
      key: "hostname",
      label: "Hostname",
      filterable: true,
      sortFunc: sortText,
    },
    {
      key: "macAddress",
      label: "MAC Address",
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
      key: "createdAt",
      label: "First Seen",
      component: [DurationUntillCell],
      sortFunc: sortDates,
    },
    {
      key: "lastSeen",
      label: "Last Seen",
      component: [DurationUntillCell],
      sortFunc: sortDates,
    },
    {
      key: "wirelessNetwork",
      label: "Wireless network",
      filterable: true,
      sortFunc: sortText,
    },
    {
      key: "wirelessAPName",
      label: "Access point",
      filterable: true,
      sortFunc: sortText,
    },
    {
      key: "note",
      label: "Note",
      filterable: true,
      sortFunc: sortText,
    },
  ];
</script>

<div>
  <section>
    <h3>All devices</h3>
    <DataTable
      {columns}
      dataPromise={fetchAllDevices}
      defaultSortKey="lastSeen"
    />
  </section>
</div>

<style>
  .cards-row {
    display: flex;
    flex-direction: row;
  }
</style>
