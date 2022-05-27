<script lang='ts'>
  // Imports
  import { get } from 'svelte/store'
  import { PeekFile, QueryFile, ReadFile } from '../../wailsjs/go/main/WApp'

  // Components
  import Throbber from '../components/Throbber.svelte'
  import AudioPlayer from '../components/AudioPlayer.svelte'

  // Stores
  import type { DirectoryViewStore } from '../stores/views'
  import { directories as directoriesStore } from '../stores/directories'

  // Properties
  export let view: DirectoryViewStore

  // Vars
  let fileCache: {[key: string]: string} = {}
  let peekCache: {[key: string]: string} = {}

  // Reactive Vars
  $: directory = directoriesStore.getByUUID($view?.directory)
  $: focusedEntry = $directory?.Entries.find(v=>get(v).Path===$view.focused)

  // Functions
  async function readFile(name: string): Promise<string> {
    if (fileCache[name]) {
      return fileCache[name]
    }
    let bytes = (await ReadFile(name)) as unknown as string
    fileCache[name] = bytes
    return fileCache[name]
  }

  async function peekFile(name: string): Promise<string> {
    if (peekCache[name]) {
      return peekCache[name]
    }
    let bytes = (await PeekFile(name, 200)) as unknown as string
    peekCache[name] = atob(bytes)
    return peekCache[name]
  }
</script>

<main>
  {#if focusedEntry}
    <section class='preview__information'>
      {#await QueryFile($directory.RealDir.Path, $focusedEntry.Path)}
        <Throbber/>
      {:then fileInfo}
        {console.log("eggs", fileInfo)}
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
            {#await peekFile(fileInfo.Path)}
              <Throbber/>
            {:then data}
              <pre>{data}</pre>
            {:catch err}
              <span>ERROR: {err}</span>
            {/await}
          {:else if fileInfo.Mimetype.startsWith('audio')}
            {#await readFile(fileInfo.Path)}
              <Throbber/>
            {:then data}
              <AudioPlayer src="data:{fileInfo.Mimetype};base64,{data}"/>
            {:catch err}
              <span>ERROR: {err}</span>
            {/await}
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
          {#if fileInfo.Mimetype.startsWith('image')}
            <li class='dimensions'>
              <label for='fileInfo__Dimensions'>dimensions</label>
              <input id='fileInfo__Dimensions' type="text" disabled value={fileInfo.Extra.Width + 'x' + fileInfo.Extra.Height}>
            </li>
            <li class='colormodel'>
              <label for='fileInfo__ColorModel'>color model</label>
              <input id='fileInfo__ColorModel' type="text" disabled value={fileInfo.Extra.ColorModel}>
            </li>
          {/if}
          <li class='permissions'>
            <label for='fileInfo__Permissions'>permissions</label>
            <input id='fileInfo__Permissions' type="text" disabled value={fileInfo.Permissions}>
          </li>
          <li class='modification'>
            <label for='fileInfo__ModificationTime'>modification</label>
            <input id='fileInfo__ModificationTime' type="text" disabled value={new Intl.DateTimeFormat('en', {dateStyle: 'medium', timeStyle: 'medium'}).format(new Date(fileInfo.ModTime))}>
          </li>
        </section>
      {:catch error}
        <span>ERROR: {error}</span>
      {/await}
    </section>
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
    overflow: hidden;
  }
  section.preview img {
    max-width: 100%;
    height: auto;
    max-height: 100%;
  }
  section.preview pre {
    max-width: 100%;
    max-height: 100%;
    text-align: left;
    word-wrap: break-word;
    white-space: pre-wrap;
  }
  li {
    list-style: none;
    padding: 0;
    margin: 0;
    display: grid;
    grid-template-columns: 6.5em minmax(0, 1fr);
    grid-template-rows: minmax(0, 1fr);
    font-size: 10pt;
    padding: .15em 0;
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
  audio {

  }
</style>