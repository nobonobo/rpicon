<script>
  import { onMount } from "svelte";
  import Navibar from "$lib/Navibar.svelte";
  import Stick from "$lib/Stick.svelte";
  import Button from "$lib/Button.svelte";

  import { getMapping } from "gamepad-api-mappings";

  const padValue = {
    L_STICK: {
      PRESSED: false,
      X_VALUE: 0,
      Y_VALUE: 0,
    },
    R_STICK: {
      PRESSED: false,
      X_VALUE: 0,
      Y_VALUE: 0,
    },
    DPAD_UP: false,
    DPAD_DOWN: false,
    DPAD_LEFT: false,
    DPAD_RIGHT: false,
    L: false,
    ZL: false,
    R: false,
    ZR: false,
    PLUS: false,
    MINUS: false,
    HOME: false,
    CAPTURE: false,
    Y: false,
    X: false,
    B: false,
    A: false,
  };
  let deviceIndex = 0;
  let devices = [];
  let connection = null;
  const profiles = [
    { name: "Pro Controller", func: buildForProCon, deadzone: 0.2 },
    {
      name: "Pro Controller for RaceGame",
      func: buildForRaceGame,
      deadzone: 0.2,
    },
    {
      name: "for Thrustmaster T80 Wheel",
      func: buildForWheel,
      deadzone: 0.15686,
    },
  ];
  let profile = 0;

  function normalize(s, deadzone) {
    let l = Math.abs(s);
    if (l < deadzone) {
      return 0;
    }
    let normal = (l - deadzone) / (1 - deadzone);
    if (s < 0) return -normal;
    return normal;
  }

  function radial(coord, deadzone = 0) {
    const angle = Math.atan2(coord.y, coord.x);
    let magnitude = Math.sqrt(coord.x * coord.x + coord.y * coord.y);
    if (magnitude <= deadzone) {
      return { x: 0, y: 0 };
    }
    if (magnitude > 1) {
      magnitude = 1;
    }
    return {
      x: Math.cos(angle) * normalize(magnitude, deadzone),
      y: Math.sin(angle) * normalize(magnitude, deadzone),
    };
  }

  function buildForProCon(pad) {
    padValue.L_STICK.PRESSED = pad.buttons[10].pressed;
    padValue.L_STICK.X_VALUE = (pad.axes[0] * 100) | 0;
    padValue.L_STICK.Y_VALUE = (pad.axes[1] * -100) | 0;
    padValue.R_STICK.PRESSED = pad.buttons[11].pressed;
    padValue.R_STICK.X_VALUE = (pad.axes[2] * 100) | 0;
    padValue.R_STICK.Y_VALUE = (pad.axes[3] * -100) | 0;
    padValue.A = pad.buttons[0].pressed;
    padValue.B = pad.buttons[1].pressed;
    padValue.X = pad.buttons[2].pressed;
    padValue.Y = pad.buttons[3].pressed;
    padValue.L = pad.buttons[4].pressed;
    padValue.R = pad.buttons[5].pressed;
    padValue.ZL = pad.buttons[6].pressed;
    padValue.ZR = pad.buttons[7].pressed;
    padValue.MINUS = pad.buttons[8].pressed;
    padValue.PLUS = pad.buttons[9].pressed;
    if (pad.buttons.length > 12) {
      padValue.DPAD_UP = pad.buttons[12].pressed;
      padValue.DPAD_DOWN = pad.buttons[13].pressed;
    }
    if (pad.buttons.length > 14) {
      padValue.DPAD_LEFT = pad.buttons[14].pressed;
      padValue.DPAD_RIGHT = pad.buttons[15].pressed;
    }
    if (pad.buttons.length > 16) {
      padValue.HOME = pad.buttons[16].pressed;
    }
    if (pad.buttons.length > 17) {
      padValue.CAPTURE = pad.buttons[17].pressed;
    }
  }

  function buildForRaceGame(pad) {
    let x = pad.axes[0] * 1.2;
    padValue.L_STICK.PRESSED = pad.buttons[10].pressed;
    padValue.L_STICK.X_VALUE = (x * x * x * 100) | 0;
    padValue.L_STICK.Y_VALUE = (pad.axes[1] * -100) | 0;
    padValue.R_STICK.PRESSED = pad.buttons[11].pressed;
    padValue.R_STICK.X_VALUE = (pad.axes[2] * 100) | 0;
    padValue.R_STICK.Y_VALUE = (pad.axes[3] * -100) | 0;
    padValue.A = pad.buttons[0].pressed;
    padValue.B = pad.buttons[1].pressed;
    padValue.X = pad.buttons[2].pressed;
    padValue.Y = pad.buttons[3].pressed;
    padValue.L = pad.buttons[4].pressed;
    padValue.R = pad.buttons[5].pressed;
    padValue.ZL = pad.buttons[6].pressed;
    padValue.ZR = pad.buttons[7].pressed;
    padValue.MINUS = pad.buttons[8].pressed;
    padValue.PLUS = pad.buttons[9].pressed;
    if (pad.buttons.length > 12) {
      padValue.DPAD_UP = pad.buttons[12].pressed;
      padValue.DPAD_DOWN = pad.buttons[13].pressed;
    }
    if (pad.buttons.length > 14) {
      padValue.DPAD_LEFT = pad.buttons[14].pressed;
      padValue.DPAD_RIGHT = pad.buttons[15].pressed;
    }
    if (pad.buttons.length > 16) {
      padValue.HOME = pad.buttons[16].pressed;
    }
    if (pad.buttons.length > 17) {
      padValue.CAPTURE = pad.buttons[17].pressed;
    }
  }

  function buildForWheel(pad) {
    let x = normalize(pad.axes[0], 0.157);
    padValue.L_STICK.X_VALUE = (x * 100) | 0;
    padValue.L_STICK.Y_VALUE = 0;
    padValue.R_STICK.X_VALUE = 0;
    padValue.R_STICK.Y_VALUE = ((pad.axes[5] - pad.axes[4]) * 50) | 0;
    padValue.Y = pad.buttons[0].pressed;
    padValue.B = pad.buttons[1].pressed;
    padValue.A = pad.buttons[2].pressed;
    padValue.X = pad.buttons[3].pressed;
    padValue.L = pad.buttons[4].pressed;
    padValue.R = pad.buttons[5].pressed;
    padValue.ZL = pad.buttons[6].pressed;
    padValue.ZR = pad.buttons[7].pressed;
    padValue.MINUS = pad.buttons[8].pressed;
    padValue.PLUS = pad.buttons[9].pressed;
    padValue.CAPTURE = pad.buttons[10].pressed;
    padValue.HOME = pad.buttons[11].pressed;
  }

  onMount(() => {
    var frame = null;
    //var old = JSON.parse(JSON.stringify(padValue));
    console.log("mounted");

    function update() {
      var gamepads = navigator.getGamepads();
      for (var i = 0; i < gamepads.length; i++) {
        var pad = gamepads[i];
        if (pad.index != deviceIndex) continue;
        if (pad) {
          profiles[profile].func(pad);
          if (connection != null && connection.readyState == WebSocket.OPEN)
            connection.send(JSON.stringify(padValue));
        }
      }
      frame = window.requestAnimationFrame(update);
      /*
      if (JSON.stringify(padValue) != JSON.stringify(old)) {
        old = JSON.parse(JSON.stringify(padValue));
        console.log(old);
      }
      */
    }
    async function connect(ev) {
      console.log(ev.gamepad.id, getMapping(ev.gamepad.id, ev.gamepad.mapping));
      console.log("gamepadconnected:", ev.gamepad);
      devices[ev.gamepad.index] = ev.gamepad.id;
    }
    async function disconnect(ev) {
      if (ev.gamepad.index != 0) return;
      console.log("gamepaddisconnected:", ev.gamepad);
    }
    window.addEventListener("gamepadconnected", connect);
    window.addEventListener("gamepaddisconnected", disconnect);
    update();
    return () => {
      window.removeEventListener("gamepadconnected", connect);
      window.removeEventListener("gamepaddisconnected", disconnect);
      cancelAnimationFrame(frame);
    };
  });
</script>

<Navibar
  {devices}
  bind:device={deviceIndex}
  {profiles}
  bind:profile
  bind:connection
/>
<main class="container">
  <div class="columns" style="width: 100%">
    <div class="column col-3 col-md-12" style="margin-bottom: 1em">
      <Button
        class="float-left"
        label={"LS"}
        bind:value={padValue.L_STICK.PRESSED}
      />
      <Button class="float-right" label={"Cap"} bind:value={padValue.CAPTURE} />
      <Stick
        bind:X_VALUE={padValue.L_STICK.X_VALUE}
        bind:Y_VALUE={padValue.L_STICK.Y_VALUE}
      />
    </div>
    <div class="column col-3 col-md-12" style="margin-bottom: 1em">
      <div class="columns">
        <div class="btn-group column col-12">
          <Button label={"ZL"} bind:value={padValue.ZL} />
          <Button label={"L"} bind:value={padValue.L} />
          <Button label={"-"} bind:value={padValue.MINUS} />
        </div>
      </div>
      <div class="columns" style="padding-top: 1em;">
        <div class="column col-3" />
        <Button
          class="column col-6"
          label={"↑"}
          bind:value={padValue.DPAD_UP}
        />
        <div class="column col-3" />
      </div>
      <div class="columns" style="padding-top: 1em;">
        <div class="btn-group column col-12">
          <Button label={"←"} bind:value={padValue.DPAD_LEFT} />
          <Button label={"→"} bind:value={padValue.DPAD_RIGHT} />
        </div>
      </div>
      <div class="columns" style="padding-top: 1em;">
        <div class="column col-3" />
        <Button
          class="column col-6"
          label={"↓"}
          bind:value={padValue.DPAD_DOWN}
        />
        <div class="column col-3" />
      </div>
    </div>
    <div class="column col-3 col-md-12" style="margin-bottom: 1em">
      <Button
        class="float-left"
        label={"RS"}
        bind:value={padValue.R_STICK.PRESSED}
      />
      <Button class="float-right" label={"Home"} bind:value={padValue.HOME} />
      <Stick
        bind:X_VALUE={padValue.R_STICK.X_VALUE}
        bind:Y_VALUE={padValue.R_STICK.Y_VALUE}
      />
    </div>
    <div class="column col-3 col-md-12" style="margin-bottom: 1em">
      <div class="columns">
        <div class="btn-group column col-12">
          <Button label={"+"} bind:value={padValue.PLUS} />
          <Button label={"R"} bind:value={padValue.R} />
          <Button label={"ZR"} bind:value={padValue.ZR} />
        </div>
      </div>
      <div class="columns" style="padding-top: 1em;">
        <div class="column col-3" />
        <Button class="column col-6" label={"X"} bind:value={padValue.X} />
        <div class="column col-3" />
      </div>
      <div class="columns" style="padding-top: 1em;">
        <div class="btn-group column col-12">
          <Button label={"Y"} bind:value={padValue.Y} />
          <Button label={"A"} bind:value={padValue.A} />
        </div>
      </div>
      <div class="columns" style="padding-top: 1em;">
        <div class="column col-3" />
        <Button class="column col-6" label={"B"} bind:value={padValue.B} />
        <div class="column col-3" />
      </div>
    </div>
  </div>
</main>

<style>
  .container {
    padding: 1rem;
  }
</style>
