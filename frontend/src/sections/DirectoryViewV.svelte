<script lang='ts'>
  import { lib } from '../../wailsjs/go/models'

  import type { DirectoryView, TagsView } from '../models/views'

  export let view: DirectoryView
  export let directory: lib.Directory
  export let tree: Object

  $: entries = Object.entries(tree)
  $: folders = entries.filter(v=>!(v[1] instanceof lib.DirectoryEntry))
  $: files = entries.filter(v=>v[1] instanceof lib.DirectoryEntry)

  async function travel(to) {
    let t = [view.wd, to].filter(v=>v!=='').join('/')
    console.log('travel', t)
  }
</script>

<main>
  {#each folders as [key, folder] }
    <div class="item folder" on:dblclick={async ()=>await travel(key)}>
      <span>folder</span>
      <span class='title'>
        {key}
      </span>
    </div>
  {/each}
  {#each files as [key, file] }
    <div class="item file">
      <span>file</span>
      <span class='title'>
        {key}
      </span>
    </div>
  {/each}
</main>

<style>
  main {
    display: flex;
    flex-wrap: wrap;
    align-content: flex-start;
    justify-content: flex-start;
    overflow-y: scroll;
  }
  .item {
    display: inline-flex;
    width: 6em;
    height: 6em;
    padding: .5em;
    margin: .5em;
    display: grid;
    grid-template-rows: auto minmax(0, 1fr);
    overflow: hidden;
    user-select: none;
    -webkit-user-select: none;
  }
  .item.folder {
    background: green;
  }
  .item.file {
    background: red;
  }
  .title {
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>