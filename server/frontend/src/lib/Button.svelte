<script>
  import { onMount } from "svelte";

  let klass = "";
  export { klass as class };
  export let label;
  export let value;

  let classes;
  $: if (value) {
    classes = ["btn", "active"];
    classes = classes.concat(klass.split(" "));
  } else {
    classes = ["btn"];
    classes = classes.concat(klass.split(" "));
  }

  function down(ev) {
    ev.preventDefault();
    value = true;
  }
  function up(ev) {
    value = false;
  }
</script>

<button
  class={classes.join(" ")}
  on:touchstart={down}
  on:touchend={up}
  on:touchcancel={up}
  on:mousedown={down}
  on:mouseup={up}
  on:mouseleave={up}>{label}</button
>
