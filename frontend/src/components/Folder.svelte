<script lang='ts'>
  import { actionPublisher } from '../actions'

  export let key: string = ''
  export let title: string = ''

  import FolderIcon from './FolderIcon.svelte'
</script>

<li data-folder-id={key} on:click>
  <div class="item folder">
    <span class='icon'>
      <FolderIcon/>
    </span>
    <span class='title'>
      {title}
    </span>
  </div>
</li>


<style>
  li {
    list-style: none;
    border: 1px solid transparent;
  }
  li .title {
    border: 1px solid transparent;
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
    align-items: center;
    justify-content: center;
  }
  .title {
    overflow: hidden;
    text-overflow: ellipsis;
    word-break: break-all;
  }
</style>