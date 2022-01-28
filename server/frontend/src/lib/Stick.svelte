<script>
  import { onMount } from "svelte";

  export let X_VALUE = 0;
  export let Y_VALUE = 0;
  let nob;
  let VX = 0;
  let VY = 0;

  $: {
    let angle = Math.atan2(-Y_VALUE, X_VALUE);
    let magnitude = Math.sqrt(X_VALUE * X_VALUE + Y_VALUE * Y_VALUE);
    if (magnitude > 100) magnitude = 100;
    VX = magnitude * Math.cos(angle);
    VY = magnitude * Math.sin(angle);
  }
  let drag = {
    isMouseDown: false,
    target: null,
  };
  onMount(() => {
    document.addEventListener("mouseup", mouseup);
    document.addEventListener("mousemove", mousemove);
    return () => {
      document.removeEventListener("mouseup", mouseup);
      document.removeEventListener("mousemove", mousemove);
    };
  });
  function mousedown(e) {
    e.preventDefault();
    drag.isMouseDown = true;
    drag.target = e.target;
    return false;
  }
  function mouseup() {
    drag.isMouseDown = false;
    VX = 0;
    VY = 0;
  }
  function mousemove(e) {
    if (drag.isMouseDown == true) {
      let parent = drag.target.parentNode;
      let rect = parent.getBoundingClientRect();
      let w = rect.width / 2;
      let h = rect.height / 2;
      let x = ((e.clientX - rect.x - w) * 1.2) / w;
      let y = ((e.clientY - rect.y - h) * 1.2) / h;
      //if (Math.abs(x) > 1) x = x / Math.abs(x);
      //if (Math.abs(y) > 1) y = y / Math.abs(y);
      let angle = Math.atan2(y, x);
      let magnitude = Math.sqrt(x * x + y * y);
      if (magnitude > 1.0) magnitude = 1.0;
      VX = magnitude * Math.cos(angle) * 100;
      VY = magnitude * Math.sin(angle) * 100;
      X_VALUE = x * 100;
      Y_VALUE = y * -100;
    }
  }
</script>

<svg
  version="1.1"
  xmlns="http://www.w3.org/2000/svg"
  viewBox="-120 -120 240 240"
>
  <g>
    <circle class="st0" cx="0" cy="0" r="100" />
    <path class="st1" d="M-100,0 L100,0" />
    <path class="st1" d="M0,-100 L0,100" />
  </g>
  <circle on:mousedown={mousedown} class="st2" cx={VX} cy={VY} r="10" />
</svg>

<style>
  .st0 {
    stroke: black;
    stroke-width: 1;
    fill: #f0f0f0;
  }
  .st1 {
    stroke: black;
    stroke-width: 1;
  }
  .st2 {
    stroke: black;
    stroke-width: 3;
    fill: white;
    cursor: move;
  }
</style>
