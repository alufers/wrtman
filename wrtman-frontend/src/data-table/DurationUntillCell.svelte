<script lang="ts">
  import { onDestroy } from "svelte";

  import { formatDuration } from "../util";

  import type ColumnModel from "./ColumnModel";

  export let row: any;
  export let column: ColumnModel<any>;
 
  let currentDiff = formatDuration(
    Math.floor((row[column.key].getTime() - new Date().getTime()) / 1000)
  );
  const interval = setInterval(
    () =>
      (currentDiff = formatDuration(
        Math.floor((row[column.key].getTime() - new Date().getTime()) / 1000)
      )),
    1000
  );

  onDestroy(() => clearInterval(interval));

 
</script>

{currentDiff}
