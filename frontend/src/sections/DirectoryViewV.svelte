<script lang='ts'>
  import { lib } from '../../wailsjs/go/models'
  import { actionPublisher } from '../actions'

  import type { DirectoryView, TagsView } from '../models/views'

  export let view: DirectoryView
  export let directory: lib.Directory
  export let tree: Object

  $: entries = Object.entries(tree)
  $: folders = entries.filter(v=>!(v[1] instanceof lib.DirectoryEntry))
  $: files = entries.filter(v=>v[1] instanceof lib.DirectoryEntry)

  async function travel(to) {
    //let t = [view.wd, to].filter(v=>v!=='').join('/')
    actionPublisher.publish('view-directory-navigate', {
      uuid: view.uuid,
      path: to,
    })
  }
</script>

<main>
  {#if view.wd != ""}
    <li on:dblclick={async ()=>await travel("..")}>
      <div class="item folder">
        <span class='icon'>folder</span>
        <span class='title'>
          ..
        </span>
      </div>
    </li>
  {/if}
  {#each folders as [key, folder] }
    <li on:dblclick={async ()=>await travel(key)}>
      <div class="item folder">
        <span>folder</span>
        <span class='title'>
          {key}
        </span>
      </div>
    </li>
  {/each}
  {#each files as [key, file] }
    <li>
      <div class="item file">
        <span class='icon'>file</span>
        <span class='title'>
          {key}
        </span>
      </div>
    </li>
  {/each}
</main>

<style>
  main {
    display: flex;
    flex-wrap: wrap;
    align-content: flex-start;
    justify-content: flex-start;
    overflow-y: scroll;
    font-size: 80%;
  }
  li {
    list-style: none;
    padding: .5em;
  }
  .item {
    display: inline-flex;
    width: 6em;
    height: 6em;
    padding: .5em;
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
  .icon {
    width: 2em;
    height: 2em;
  }
  .title {
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>