<script lang='ts'>
  import { lib } from '../../wailsjs/go/models'
  import { actionPublisher } from '../actions'

  import type { DirectoryView, TagsView } from '../models/views'

  export let view: DirectoryView
  export let directory: lib.Directory
  export let tree: Object

  $: selectedFiles = view.selected

  let showBox: boolean = false
  let box: [number, number, number, number] = [0,0,0,0]
  let mainElement: HTMLElement
  $: boxRepresentation = [
    box[0]<box[2]?box[0]:box[2],
    box[1]<box[3]?box[1]:box[3],
    box[0]<box[2]?box[2]:box[0],
    box[1]<box[3]?box[3]:box[1],
  ]

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

  async function fileMousedown(e: MouseEvent, d: lib.DirectoryEntry) {
    let hits = [d.Path]
    if (e.shiftKey) {
      selectedFiles = [...new Set([...hits,...selectedFiles])]
    } else if (e.ctrlKey) {
      selectedFiles = selectedFiles.filter(v=>hits.includes(v))
    } else {
      selectedFiles = hits
    }
    actionPublisher.publish('view-select-files', {
      uuid: view.uuid,
      selected: selectedFiles,
    })
  }

  function viewMousedown(e: MouseEvent) {
    if (e.button !== 0) return
    e.stopPropagation()
    e.preventDefault()
    let b = mainElement.getBoundingClientRect()
    showBox = true
    box[0] = box[2] = e.clientX - b.left
    box[1] = box[3] = e.clientY - b.top
    box = [...box]
    let viewMousemove = (e: MouseEvent) => {
      box[2] = e.clientX - b.left
      box[3] = e.clientY - b.top
      box = [...box]
    }
    let viewMouseup = (e: MouseEvent) => {
      let nodes = mainElement.getElementsByTagName('LI')
      let els = Array.from(nodes).filter(v=>v.getAttribute('data-file-id'))
      let hits = []
      for (let el of els) {
        let elb = el.getBoundingClientRect()
        let x = elb.left - b.left
        let y = elb.top - b.top
        let isOverlapping = (x1min: number, x1max: number, x2min: number, x2max: number, y1min: number, y1max: number, y2min: number, y2max: number) => (x1min < x2max && x2min < x1max && y1min < y2max && y2min < y1max)
        // Check for constrain as well as contain
        if (isOverlapping(x, x+elb.width, boxRepresentation[0], boxRepresentation[2], y, y+elb.height, boxRepresentation[1], boxRepresentation[3])) {
          hits.push(el.getAttribute('data-file-id'))
        }
      }
      // TODO: Shift to add to selection, Control to remove
      if (e.shiftKey) {
        selectedFiles = [...new Set([...hits,...selectedFiles])]
      } else if (e.ctrlKey) {
        selectedFiles = selectedFiles.filter(v=>hits.includes(v))
      } else {
        selectedFiles = hits
      }

      actionPublisher.publish('view-select-files', {
        uuid: view.uuid,
        selected: selectedFiles,
      })

      showBox = false
      window.removeEventListener('mouseup', viewMouseup)
      window.removeEventListener('mousemove', viewMousemove)
    }
    window.addEventListener('mouseup', viewMouseup)
    window.addEventListener('mousemove', viewMousemove)
  }
</script>

<main bind:this={mainElement} on:mousedown={viewMousedown}>
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
    <li data-folder-id={key} on:dblclick={async ()=>await travel(key)}>
      <div class="item folder">
        <span>folder</span>
        <span class='title'>
          {key}
        </span>
      </div>
    </li>
  {/each}
  {#each files as [key, file] }
    <li data-file-id={key} class:selected={selectedFiles.includes(key)}>
      <div on:mousedown|stopPropagation={e=>fileMousedown(e, file)} class="item file">
        <span class='icon'>file</span>
        <span class='title'>
          {key}
        </span>
      </div>
    </li>
  {/each}
  {#if showBox}
    <div class='box' style="left: {boxRepresentation[0]}px; top: {boxRepresentation[1]}px; width: {boxRepresentation[2]-boxRepresentation[0]}px; height: {boxRepresentation[3]-boxRepresentation[1]}px"></div>
  {/if}
</main>

<style>
  main {
    position: relative;
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
  li.selected .title {
    background: var(--primary-light);
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
  }
  .icon {
    width: 2em;
    height: 2em;
  }
  .title {
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .box {
    position: absolute;
    background: var(--primary-light);
    opacity: 0.5;
    pointer-events: none;
  }
</style>