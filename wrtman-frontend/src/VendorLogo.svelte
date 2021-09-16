<script lang="ts">
  import type IDHCPLease from "./model/IDHCPLease";
  import { getVendorLogoURL } from "./vendorLogos";

  export let row: { vendor: string; hostname: string };
  console.log(row);
  let vendorLogo = getVendorLogoURL(row);
</script>

<div>
  {#if vendorLogo}
    <img src={vendorLogo} alt="{row.vendor} logo" width="24" />
  {:else}
    <div class="spacer">&nbsp;</div>
  {/if}
</div>

<style>
  .spacer {
    width: 24px;
    height: 24px;
  }
</style>
