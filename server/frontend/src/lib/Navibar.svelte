<script>
  import { onMount } from "svelte";

  export let devices = ["No Device"];
  export let device = 0;
  export let connection = null;
  export let profiles = [];
  export let profile = 0;

  let uri;
  let disabled = false;
  let cancel = null;
  onMount(() => {
    let protocol = "ws://";
    if (location.protocol === "https:") {
      protocol = "wss://";
    }
    uri = protocol + location.host + "/api/ws";
  });
  function change(ev) {
    console.log(ev);
    disabled = true;
    if (ev.target.checked) {
      console.log("connecting:", uri);
      if (cancel != null) {
        cancel();
      }
      let conn = new WebSocket(uri);
      conn.onclose = () => {
        console.log("disconnected");
        connection = null;
        disabled = false;
        cancel = null;
      };
      conn.onopen = () => {
        console.log("connected");
        connection = conn;
        cancel = () => {
          conn.close();
          cancel = null;
        };
        disabled = false;
      };
      conn.onerror = () => {
        console.log("error");
        connection = null;
        ev.target.checked = false;
        disabled = false;
      };
    } else {
      if (cancel != null) {
        console.log("disconnecting");
        cancel();
      }
    }
  }
</script>

<header class="navbar">
  <section class="navbar-section">
    <a href="/" class="navbar-brand">ProCon Emulator</a>
    <div class="input-group" style="padding-left: 1em">
      <span class="input-group-addon">device</span>
      <select class="form-select" bind:value={device}>
        {#each devices as d, i}
          <option value={i}>{d}</option>
        {/each}
      </select>
    </div>
    <div class="input-group" style="padding-left: 1em">
      <span class="input-group-addon">profile</span>
      <select class="form-select" bind:value={profile}>
        {#each profiles as p, i}
          <option value={i}>{p.name}</option>
        {/each}
      </select>
    </div>
  </section>
  <section class="navbar-section">
    <div class="form-group">
      <label class="form-switch form-inline">
        <input type="checkbox" on:change={change} {disabled} />
        <i class="form-icon" /> Connect
      </label>
    </div>
  </section>
</header>

<style>
  .navbar {
    box-shadow: lightgrey 2px 2px 2px;
    padding: 1rem;
  }
</style>
