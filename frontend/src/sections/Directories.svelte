<script lang='ts'>
  import type { lib } from '../../wailsjs/go/models'

  export let directories: lib.Directory[] = []
</script>

<ul>
  {#each directories as directory}
    <li>{directory.Path}</li>
  {/each}
</ul>

<style>
  ul {
    list-style: none;
    margin: 0;
    padding: .5em;
    text-align: left;
  }
  li {
    list-style: none;
    margin: 0;
    padding: 0;
  }
</style>