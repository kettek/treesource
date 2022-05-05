<script lang='ts'>
  import type { DirectoryView, TagsView } from 'src/models/views'
  import { actionPublisher } from '../actions'

  export let tagsViews: TagsView[]
  export let directoryViews: DirectoryView[]
  let selectedView: string | number[]  = ''
</script>

<main>
  <nav>menu of some sort</nav>
  <nav class='tabs'>
    {#each directoryViews as tab}
      <li class:selected={selectedView===tab.uuid} on:click={()=>selectedView=tab.uuid}>
        <span class='title'> dir: {tab.uuid} </span>
        <button class='close' on:click={()=>actionPublisher.publish('view-directory-remove', tab.uuid)}>x</button>
      </li>
    {/each}
    {#each tagsViews as tab}
      <li class:selected={selectedView===tab.uuid} on:click={()=>selectedView=tab.uuid}>
        <span class='title'> tags: {tab.uuid} </span>
        <button class='close' on:click={()=>actionPublisher.publish('view-tags-remove', tab.uuid)}>x</button>
      </li>
    {/each}
  </nav>
  <section>content of current view</section>
</main>

<style>
  main {
    display: grid;
    grid-template-rows: auto auto minmax(0, 1fr);
  }
  nav.tabs {
    display: flex;
    align-items: center;
    justify-content: flex-start;
    background: var(--secondary-dark);
    grid-column-gap: .5em;
    height: 2em;
  }
  li {
    list-style: none;
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto;
    grid-template-rows: minmax(0, 1fr);
    align-items: center;
    justify-content: flex-start;
    background: var(--secondary);
    height: 1.5em;
    overflow: hidden;
    padding: 0;
    margin: 0;
    cursor: default;
  }
  li.selected {
    background: var(--secondary-light);
  }
  li span {
    overflow: auto;
  }
  li button {
    margin: 0;
    padding: 0;
    width: 2em;
    height: 2em;
    border: 0;
    background: none;
    cursor: pointer;
  }
</style>