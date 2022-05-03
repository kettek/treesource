<script lang='ts'>
  import Menus from '../menu/Menus.svelte'
  import MenuBar from '../menu/MenuBar.svelte'
  import MenuItem from '../menu/MenuItem.svelte'
  import type { lib } from '../../wailsjs/go/models'

  export let disabled: boolean
  export let directories: lib.Directory[] = []
  let selectedDirectoryIndex: number = 0
  $: selectedDirectory = directories[selectedDirectoryIndex]
</script>

<section>
  <Menus>
    <MenuBar>
      <MenuItem action='directory-add' disabled={disabled}>add</MenuItem>
      <MenuItem action='directory-remove' args={selectedDirectory?.UUID} disabled={disabled||selectedDirectoryIndex<0||selectedDirectoryIndex>=directories.length}>remove selected</MenuItem>
    </MenuBar>
  </Menus>
  <ul>
    {#each directories as directory, i}
      <li class:selected={selectedDirectoryIndex===i} on:click={()=>selectedDirectoryIndex=i}>{directory.Path}</li>
    {/each}
  </ul>
</section>

<style>
  section {
    display: grid;
    grid-template-rows: auto minmax(0, 1fr);
    grid-template-columns: minmax(0, 1fr);
    overflow: hidden;
  }
  ul {
    list-style: none;
    margin: 0;
    padding: 0;
    text-align: left;
    background-color: var(--primary-dark);
    overflow-y: scroll;
  }
  li {
    list-style: none;
    margin: 0;
    padding: 0;
    padding: .5em .5em;
    overflow: auto;
  }
  li.selected {
    background-color: var(--primary);
  }
</style>