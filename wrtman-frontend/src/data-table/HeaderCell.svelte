<script lang="ts">
  import type ColumnModel from "./ColumnModel";
  import { mdiMenuDown, mdiMenuUp } from "@mdi/js";
  import Icon from "mdi-svelte";
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  export let column: ColumnModel<any>;
  export let sortKey: string;
  export let hasFilters: boolean;

  function sortAscClicked() {
    dispatch("sortKeyChanged", column.key);
  }

  function sortDescClicked() {
    dispatch("sortKeyChanged", "-" + column.key.toString());
  }
</script>

<th colspan={column?.component?.length || 1} class:no-border={hasFilters}>
  <div class="label-header">
    <div>{column.label}</div>
    {#if column.sortFunc}
      <div class="sort-icons">
        <span
          class="sort-icon"
          on:click={sortAscClicked}
          class:active={sortKey === column.key}><Icon path={mdiMenuUp} /></span
        ><span
          class="sort-icon"
          on:click={sortDescClicked}
          class:active={sortKey === "-" + column.key.toString()}
          ><Icon path={mdiMenuDown} /></span
        >
      </div>
    {/if}
  </div>
</th>

<style>

  .label-header {
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    justify-content: space-between;
    align-items: center;
  }
  .sort-icons {
    color: var(--background-dimmer);
  }
  .sort-icon:hover {
    color: #eee;
  }
  .sort-icon {
    cursor: pointer;
  }
  .sort-icon.active {
    color: var(--primary);
    cursor: initial;
  }
</style>
