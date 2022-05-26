<script lang='ts'>
  import { actionPublisher } from '../actions'
  import type { DirectoryViewStore } from '../stores/views'
  import { directories as directoriesStore } from '../stores/directories'
  import { get } from 'svelte/store'
  import { SplitPath } from '@kettek/filepaths-to-tree'

  export let selected: boolean = false
  export let view: DirectoryViewStore

  let title: string = ''
  $: {
    title = get(directoriesStore.getByUUID($view.directory))?.RealDir?.Path
    if (title) {
      let s = SplitPath(title)
      title = s[s.length-1]
    } else {
      title = $view.directory as unknown as string
    }
  }
</script>

<li class:selected on:click={()=>actionPublisher.publish('view-select', $view.uuid)}>
  <span class='title'>d: {title}</span>
  <button class='close' on:click={()=>actionPublisher.publish('view-directory-remove', $view.uuid)}>x</button>
</li>

<style>
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