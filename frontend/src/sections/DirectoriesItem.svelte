<script lang='ts'>
  import { actionPublisher } from '..//actions'

  import type { DirectoryStore } from '../stores/directories'

  export let selected: boolean = false
  export let directory: DirectoryStore
</script>

<li class:selected={selected} on:click on:dblclick={()=>actionPublisher.publish('view-directory-add', $directory.RealDir.UUID)}>
  <span class='title'>{$directory.RealDir.Path}</span>
  <span class='count'>{$directory.Entries.length}</span>
</li>

<style>
  li {
    list-style: none;
    margin: 0;
    padding: 0;
    padding: .5em .5em;
    overflow: auto;
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto;
    align-items: center;
    cursor: pointer;
  }
  li.selected {
    background-color: var(--primary);
  }
  .count {
    font-size: 60%;
  }
</style>