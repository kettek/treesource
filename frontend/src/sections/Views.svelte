<script lang='ts'>
  import { DirectoryView, TagsView } from '../models/views'
  import type { lib } from 'wailsjs/go/models'
  import { actionPublisher } from '../actions'
  import { SplitPath, Find } from '@kettek/filepaths-to-tree'
  import DirectoryViewV from './DirectoryViewV.svelte'

  export let directories: lib.Directory[] = []
  export let directoryTrees: any = {}

  export let tagsViews: TagsView[]
  export let directoryViews: DirectoryView[]
  export let selectedView: string | number[]  = ''

  $: directoryView = directoryViews.find(v=>v.uuid===selectedView)
  $: directory = directoryView ? directories.find(v => v.UUID === directoryView.directory) : undefined

  function getDirectoryTitle(uuid: string): string {
    let d = directories.find(v=> String(v.UUID) === uuid)
    if (d) {
      let s = SplitPath(d.Path)
      return s[s.length-1]
    }
    return uuid
  }

</script>

<main>
  <nav>menu of some sort</nav>
  <nav class='tabs'>
    {#each directoryViews as tab}
      <li class:selected={selectedView===tab.uuid} on:click={()=>actionPublisher.publish('view-select', tab.uuid)}>
        <span class='title'>d: {getDirectoryTitle(String(tab.directory))}</span>
        <button class='close' on:click={()=>actionPublisher.publish('view-directory-remove', tab.uuid)}>x</button>
      </li>
    {/each}
    {#each tagsViews as tab}
      <li class:selected={selectedView===tab.uuid} on:click={()=>actionPublisher.publish('view-select', tab.uuid)}>
        <span class='title'>t: {tab.uuid}</span>
        <button class='close' on:click={()=>actionPublisher.publish('view-tags-remove', tab.uuid)}>x</button>
      </li>
    {/each}
  </nav>
  <section>
    {#if directory}
      <DirectoryViewV view={directoryView} directory={directory} tree={Find(directoryTrees[String(directory.UUID)], directoryView.wd)}></DirectoryViewV>
    {/if}
  </section>
</main>

<style>
  main {
    display: grid;
    grid-template-rows: auto auto minmax(0, 1fr);
  }
  section {
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    grid-template-rows: minmax(0, 1fr);
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