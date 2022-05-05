<script lang='ts'>
  import type { DirectoryView, TagsView } from 'src/models/views'
  import { actionPublisher } from '../actions'

  export let tagsViews: TagsView[]
  export let directoryViews: DirectoryView[]
  let selectedView: string | number[]  = ''
</script>

<main>
  <nav>menu of some sort</nav>
  <nav>
    {#each directoryViews as tab}
      <li on:click={()=>selectedView=tab.uuid}>
        <span class='title'> dir: {tab.uuid} </span>
        <button class='close' on:click={()=>actionPublisher.publish('view-directory-remove', tab.uuid)}>x</button>
      </li>
    {/each}
    {#each tagsViews as tab}
      <li on:click={()=>selectedView=tab.uuid}>
        <span class='title'> tags: {tab.uuid} </span>
        <button class='close' on:click={()=>actionPublisher.publish('view-tags-remove', tab.uuid)}>x</button>
      </li>
    {/each}
  </nav>
  <section>content of current view</section>
</main>

<style>
  li {
    list-style: none;
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto;
  }
</style>