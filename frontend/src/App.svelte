<script lang="ts">
  import folderIcon from './assets/breeze-icons/icons/places/64/folder.svg'
  import noneIcon from './assets/breeze-icons/icons/mimetypes/64/none.svg'
  import {ListFiles} from '../wailsjs/go/main/App.js'
  import type {main} from '../wailsjs/go/models'

  let currentFiles: (main.DirEntry[]|Error) = []
  let path: string

  $: inRoot = path === '/' || path === ''

  async function getFiles() {
    try {
      currentFiles = await ListFiles(path)
    } catch(err) {
      console.log(err)
    }
  }
  async function navigate(to: main.DirEntry) {
    if (!to.Info.IsDir) return
    // TODO: Send request to move to a directory, then wait for a response.
    path = path+"/" + to.Name
    await getFiles()
  }
  async function up() {
    // TODO: Send request to move up a directory, then wait for a response.
  }
</script>

<main>
  <div class="result" id="resultPath">
    <ul class='entries'>
      {#if !inRoot}
        <li class='entry' on:click={up}>
          <span class='entry__image'>
            <img src="{folderIcon}" alt='folder'>
          </span>
          <span class='entry__name'> Up </span>
        </li>
      {/if}
      {#if Array.isArray(currentFiles)}
        {#each currentFiles as entry}
          <li class='entry' on:click={()=>navigate(entry)}>
            <span class='entry__image'>
              {#if entry.Info.IsDir}
                <img src="{folderIcon}" alt='folder'>
              {:else}
                <img src="{noneIcon}" alt='file'>
              {/if}
            </span>
            <span class='entry__name'> {entry.Name} </span>
          </li>
        {/each}
      {/if}
    </ul>
  </div>
  <div class="input-box" id="input">
    <input autocomplete="off" bind:value={path} class="input" id="path" type="text" on:keyup={e=>e.key==='Enter'?getFiles():null} />
    <button class="btn" on:click={getFiles}>List Files</button>
  </div>
</main>

<style>
  .result {
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    overflow: auto;
  }
  main {
    width: 100%;
    height: 100%;
    display: grid;
    grid-template-rows: minmax(0, 1fr) auto;
    grid-template-columns: minmax(0, 1fr);
  }
  .entries {
    list-style: none;
  }
  .entry {
    display: grid;
    grid-template-columns: auto minmax(0, 1fr);
    align-items: center;
  }
  .entry__image img {
    width: 64px;
    height: 64px;
  }
</style>
