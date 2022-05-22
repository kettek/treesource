<script lang='ts'>
  import SplitPane from '../components/SplitPane.svelte'

  import { onMount } from 'svelte'
  import { QueryFile, ReadFile } from '../../wailsjs/go/main/WApp'

  import { lib } from '../../wailsjs/go/models'
  import type { DirectoryView, TagsView } from '../models/views'
  import Throbber from '../components/Throbber.svelte'

  export let directories: lib.Directory[] = []

  export let tagsViews: TagsView[]
  export let directoryViews: DirectoryView[]

  export let selectedView: string | number[]  = ''

  let fileCache: {[key: string]: string} = {}

  $: directoryView = directoryViews.find(v=>v.uuid===selectedView)
  $: directory = directoryView ? directories.find(v => v.UUID === directoryView.directory) : undefined

  let fileInfo: any

  $: focusedEntry = directory?.Entries.find(v=>v.Path===directoryView?.focused)

  async function readFile(name: string): Promise<string> {
    if (fileCache[name]) {
      return fileCache[name]
    }
    let bytes = (await ReadFile(name)) as unknown as string
    fileCache[name] = bytes
    return fileCache[name]
  }

  onMount(() => {
  })
</script>

<main>
  {#if focusedEntry}
    <SplitPane type="vertical" pos={40}>
      <section class='preview__information' slot=a>
        {#await QueryFile(directory.Path, focusedEntry.Path)}
          <Throbber/>
        {:then fileInfo}
          <section class='preview'>
            {#if fileInfo.Mimetype.startsWith('image')}
              {#await readFile(fileInfo.Path)}
                <Throbber/>
              {:then data}
                <img src="data:{fileInfo.Mimetype};base64,{data}" alt={fileInfo.Name}>
              {:catch err}
                <span>ERROR: {err}</span>
              {/await}
            {:else if fileInfo.Mimetype.startsWith('text')}
              <span>text</span>
            {/if}
          </section>
          <section class='info'>
            <li class='path'>
              <label>{fileInfo.Name}</label>
            </li>
            <li class='type'>
              <label for='fileInfo__Type'>type</label>
              <input id='fileInfo__Type' type="text" disabled value={fileInfo.Mimetype}>
            </li>
            <li class='permissions'>
              <label for='fileInfo__Permissions'>permissions</label>
              <input id='fileInfo__Permissions' type="text" disabled value={fileInfo.Permissions}>
            </li>
          </section>
        {:catch error}
          <span>ERROR: {error}</span>
        {/await}
      </section>
      <section slot=b>
        META
      </section>
    </SplitPane>
  {/if}
</main>

<style>
  main {
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    grid-template-rows: minmax(0, 1fr);
    padding: .5em;
  }
  section.preview__information {
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    grid-template-rows: minmax(0, 1fr) auto;
  }
  section.preview {
    background: black;
    color: white;
    height: 100%;
    border-radius: .25em;
  }
  section.preview img {
    max-width: 100%;
    height: auto;
    max-height: 100%;
  }
  li {
    list-style: none;
    padding: 0;
    margin: 0;
    display: grid;
    grid-template-columns: 6.5em minmax(0, 1fr);
    grid-template-rows: minmax(0, 1fr);
    font-size: 10pt;
    padding: .5em 0;
  }
  li.path {
    grid-template-columns: minmax(0, 1fr);
  }
  li.path > label {
    text-align: center;
    word-wrap: break-word;
  }
  li > label {
    font-weight: bold;
    text-align: left;
    padding: 0 .5em;
  }
  li > input {
    width: 100%;
    background: transparent;
    border: none;
  }
</style>