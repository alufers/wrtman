<script lang="ts">
  import Logo from "./Logo.svelte";
  import { Router, Link, Route } from "svelte-navigator";
  import CurrentStatusPage from "./CurrentStatusPage.svelte";
  import AllDevicesPage from "./AllDevicesPage.svelte";

  function getProps({ location, href, isPartiallyCurrent, isCurrent }) {
    const isActive = href === "/" ? isCurrent : isPartiallyCurrent || isCurrent;

    // The object returned here is spread on the anchor element's attributes
    if (isActive) {
      return { class: "active" };
    }
    return {};
  }
</script>

<Router>
  <main>
    <Logo />
    <nav>
      <Link to="/" {getProps}>Status</Link>
      <Link to="all-devices" {getProps}>All Devices</Link>
    </nav>
    <Route path="/" component={CurrentStatusPage} />
    <Route path="/all-devices" component={AllDevicesPage} />
    <footer>Made by alufers.</footer>
  </main>
</Router>

<style>
  :global(body) {
    background: black;
    color: #ddd;
    font-family: monospace;
    --primary: lightseagreen;
    --primary-dim: seagreen;
    --secondary: lightpink;
    --background-dim: #333;
    --background-dimmer: #555;
    --foreground-secondary: #bbb;
  }
  :global(a) {
    color: var(--primary);
  }
  :global(a.active) {
    text-decoration: none;
    opacity: 0.8;
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

  :global(select) {
    border: none;
    outline: none;
    background: transparent;
    color: #fff;
    border: 1px solid var(--foreground-secondary);
    padding: 4px;
    font-family: monospace;
  }

  :global(select option) {
    background: black;
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

  footer {
    margin-top: 128px;
    color: var(--foreground-secondary);
  }
</style>
