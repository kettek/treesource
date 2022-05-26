<script lang='ts'>
  import Menus from '../menu/Menus.svelte'
  import MenuBar from '../menu/MenuBar.svelte'
  import MenuItem from '../menu/MenuItem.svelte'
  import DirectoriesItem from './DirectoriesItem.svelte'

  import { directories as directoriesStore } from '../stores/directories'

  export let disabled: boolean
  let selectedDirectoryIndex: number = 0
  $: selectedDirectory = $directoriesStore.Directories[selectedDirectoryIndex]
</script>

<section>
  <Menus>
    <MenuBar>
      <MenuItem action='directory-add' disabled={disabled}>add</MenuItem>
      <MenuItem action='directory-remove' args={$selectedDirectory?.RealDir?.UUID} disabled={disabled||selectedDirectoryIndex<0||selectedDirectoryIndex>=$directoriesStore.Directories.length}>remove selected</MenuItem>
    </MenuBar>
  </Menus>
  <ul>
    {#each $directoriesStore.Directories as directory, i}
      <DirectoriesItem selected={selectedDirectoryIndex===i} directory={directory} on:click={()=>selectedDirectoryIndex=i}/>
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
    user-select: none;
    -webkit-user-select: none;
  }
</style>