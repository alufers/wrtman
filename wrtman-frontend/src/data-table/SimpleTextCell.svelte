<script lang="ts">
  import type ColumnModel from "./ColumnModel";

  export let row: any;
  export let column: ColumnModel<any>;
</script>
{row[column.key]}
