<script lang="ts">
  import type ColumnModel from "./ColumnModel";
  import HeaderCell from "./HeaderCell.svelte";
  export let columns: ColumnModel<any>[];
  export let dataPromise: Promise<any[]>;
  export let defaultSortKey: string;

  let sortKey = defaultSortKey;

  let filters = {};

  let isMobile = false;

  $: filtered = dataPromise.then((d) => {
    let key = sortKey;
    let shouldReverse = false;
    if (key.startsWith("-")) {
      key = key.substring(1);
      shouldReverse = true;
    }

    let data = d.filter((row) => {
      if (Object.keys(filters).length === 0) return true;
      return Object.keys(filters).every((fKey) => {
        if (!row[fKey]) return filters[fKey] === "";
        return row[fKey]
          .toString()
          .toLowerCase()
          .includes(filters[fKey].toLowerCase());
      });
    });

    const sortColumn =
      columns.find((c) => c.key === key) || columns.find((c) => !!c.sortFunc);

    data.sort((a, b) => {
      return sortColumn.sortFunc(a, b, sortColumn);
    });

    if (shouldReverse) {
      data.reverse();
    }

    return data;
  });

  $: nColumns = columns.map((col) => col.component || null).flat().length;

  $: hasFilters = columns.some((c) => c.component);

  function clearFilters() {
    filters = {};
  }
</script>

<table>
  {#if isMobile}
    <thead>
      <tr><th>Order by</th></tr>
      <tr
        ><td>
          <select>
            {#each columns as column}
              {#if column.sortFunc}
                <option value={column.key}>{column.label}</option>
              {/if}
            {/each}
          </select>
        </td></tr
      >
      {#each columns as column}
        {#if column.filterable}
          <tr>
            <th>Filter {column.label}</th>
          </tr>
          <tr>
            <td>
              <input
                type="text"
                placeholder="Filter {column.label}"
                class="filter-input"
                bind:value={filters[column.key]}
              />
            </td>
          </tr>
        {/if}
      {/each}
    </thead>
    <tbody>
      {#await filtered}
        <tr>
          <td class="placeholder-cell"> loading... </td>
        </tr>
      {:then data}
        {#if data.length === 0}
          <tr>
            <td class="placeholder-cell">
              <div>No data available. <br /><br /><br /></div>
              <button on:click={clearFilters}>Clear filters</button>
            </td>
          </tr>
        {/if}
        {#each data as row}
          {#each columns as column}
            <tr>
              <th>
                {column.label}
              </th>
            </tr>
            <tr>
              {#if Array.isArray(column.component)}
                {#each column.component as component}
                  <td><svelte:component this={component} {row} {column} /></td>
                {/each}
              {:else if column.component}
                <td>
                  <svelte:component this={column.component} {row} {column} />
                </td>
              {:else}
                <td>
                  {#if row[column.key] != null}
                    {row[column.key]}
                  {/if}
                </td>
              {/if}
            </tr>
          {/each}
        {/each}
      {:catch error}
        <tr>
          <td class="placeholder-cell">
            An error has occured: {error.toString()}
          </td>
        </tr>
      {/await}
    </tbody>
  {:else}
    <!-- END MOBILE -->
    <thead>
      <tr>
        {#each columns as column}
          <HeaderCell
            {column}
            {hasFilters}
            {sortKey}
            on:sortKeyChanged={(ev) => (sortKey = ev.detail)}
          />
        {/each}
      </tr>

      {#if hasFilters}
        <tr>
          {#each columns as column}
            <th colspan={column?.component?.length || 1}>
              {#if column.filterable}
                <input
                  type="text"
                  placeholder="Filter {column.label}"
                  class="filter-input"
                  bind:value={filters[column.key]}
                />
              {/if}
            </th>
          {/each}
        </tr>
      {/if}
    </thead>
    <tbody>
      {#await filtered}
        <tr>
          <td colspan={nColumns} class="placeholder-cell"> loading... </td>
        </tr>
      {:then data}
        {#if data.length === 0}
          <tr>
            <td colspan={nColumns} class="placeholder-cell">
              <div>No data available. <br /><br /><br /></div>
              <button on:click={clearFilters}>Clear filters</button>
            </td>
          </tr>
        {/if}
        {#each data as row}
          <tr class="filter-hidden">
            {#each columns as column}
              {#if Array.isArray(column.component)}
                {#each column.component as component}
                  <td><svelte:component this={component} {row} {column} /></td>
                {/each}
              {:else if column.component}
                <td>
                  <svelte:component this={column.component} {row} {column} />
                </td>
              {:else}
                <td>
                  {#if row[column.key] != null}
                    {row[column.key]}
                  {/if}
                </td>
              {/if}
            {/each}
          </tr>
        {/each}
      {:catch error}
        <tr>
          <td colspan={nColumns} class="placeholder-cell">
            An error has occured: {error.toString()}
          </td>
        </tr>
      {/await}
    </tbody>
    {#await Promise.all([filtered, dataPromise]) then data}
      <tfoot>
        <tr>
          <td colspan={nColumns}>
            Showing <strong>{data[0].length}</strong> of
            <strong>{data[1].length}</strong> items
          </td>
        </tr>
      </tfoot>
    {/await}
  {/if}
</table>

<style>
  table {
    border-collapse: collapse;
    width: 100%;
  }
  th {
    text-align: left;
  }
  td,
  th {
    padding-left: 8px;
    padding-right: 8px;
    border-bottom: 1px solid var(--background-dim);
  }
  .placeholder-cell {
    padding: 32px;
    text-align: center;
  }
  .filter-input {
    margin-top: 8px;
    margin-bottom: 8px;
    width: 100%;
  }
  tfoot td {
    padding-top: 16px;
    color: var(--foreground-secondary);
  }
</style>
