<script lang="ts">

  export let src: string = ''
  export let autoplay: boolean = true

  let duration: number = 0
  let currentTime: number = 0
  let muted: boolean = false
  let paused: boolean = true

  function scrobble(e: MouseEvent) {
    console.log(e)
  }
</script>

<main>
  <nav>
    <button on:click={_=>currentTime=0}>rewind</button>
    <button on:click={_=>paused=!paused}>{paused?'play':'pause'}</button>
  </nav>
  <div class='slider' on:mousedown={scrobble}>
    <div class='rail'></div>
    <div class='train' style="left: calc({currentTime/duration*100}% - .5em)"><!-- choo choo --></div>
  </div>
  <audio autoplay={autoplay} src={src} bind:duration bind:currentTime bind:muted bind:paused/>
</main>

<style>
  main {
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    grid-template-rows: auto auto;
    font-size: 75%;
    padding: 1em;
  }
  audio {
    display: none
  }
  .slider {
    position: relative;
    width: 100%;
    height: 2em;
    background: white;
    border-radius: .5em;
  }
  .rail {
    position: absolute;
    width: 100%;
    top: 0.9em;
    height: .2em;
    background: gray;
    border-radius: .5em;
  }
  .train {
    position: absolute;
    left: -.5em;
    top: 0em;
    width: 1em;
    height: 2em;
    background: red;
    border-radius: 2em;
  }
</style>