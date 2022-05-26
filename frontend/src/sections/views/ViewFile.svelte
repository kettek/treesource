<script lang='ts'>
  import type { DirectoryEntryStore, DirectoryStore } from '../../stores/directories'
  import type { DirectoryViewStore } from '../../stores/views'
  import { OpenFile } from '../../../wailsjs/go/main/WApp.js'
  import FileIcon from '../../components/FileIcon.svelte'

  export let key: string = ''
  export let directory: DirectoryStore
  export let view: DirectoryViewStore
  export let file: DirectoryEntryStore
</script>

<li data-file-id={$file.Path} class:selected={$view.selected.includes($file.Path)} class:focused={$file.Path===$view.focused} on:dblclick={async ()=>await OpenFile([$directory.RealDir.Path, $file.Path])}>
  <div class="item file">
    <span class='icon'>
      <FileIcon paths={[$directory.RealDir.Path, $file.Path]}/>
    </span>
    <span class='title'>
      {key}
    </span>
  </div>
</li>

<style>
  li {
    list-style: none;
    border: 1px solid transparent;
    margin: .1em;
  }
  li.selected {
    background: var(--primary-light);
  }
  .item {
    display: inline-flex;
    flex-direction: column;
    align-items: center;
    justify-content: flex-start;
    width: 8em;
    height: 12em;
    overflow: hidden;
    user-select: none;
    -webkit-user-select: none;
  }
  .icon {
    width: 80%;
    height: 65%;
    display: flex;
    align-items: flex-end;
    justify-content: center;
  }
  .title {
    overflow: hidden;
    text-overflow: ellipsis;
    word-break: break-all;
  }
</style>