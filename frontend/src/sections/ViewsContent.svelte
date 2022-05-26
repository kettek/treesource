<script lang='ts'>
  import { actionPublisher } from '../actions'

  import { Find } from '@kettek/filepaths-to-tree'

  import type { DirectoryEntryStore } from '../stores/directories'

  import type { DirectoryViewStore } from '../stores/views'
  import { directories as directoriesStore } from '../stores/directories'
  import ViewFile from './views/ViewFile.svelte'

  export let view: DirectoryViewStore

  $: directory = directoriesStore.getByUUID($view.directory)

  $: selectedFiles = $view.selected
  $: focusedFile = $view.focused

  $: entries = Object.entries(Find($directory?.Tree||{}, $view.wd)||[]) as [string, DirectoryEntryStore|any]
  $: folders = entries.filter(v=>!v[1].update)
  $: files = entries.filter(v=>v[1].update)

  // travel/locations
  $: crumbs = $view?.wd.split($directory?.RealDir?.Separator).map((v, i, a) => {
    return [v, a.slice(0, i+1).join($directory.RealDir.Separator)]
  })

  async function travel(to) {
    actionPublisher.publish('view-directory-navigate', {
      uuid: $view.uuid,
      path: to,
    })
  }

  // box/selection
  let showBox: boolean = false
  let box: [number, number, number, number] = [0,0,0,0]
  let mainElement: HTMLElement
  $: boxRepresentation = [
    box[0]<box[2]?box[0]:box[2],
    box[1]<box[3]?box[1]:box[3],
    box[0]<box[2]?box[2]:box[0],
    box[1]<box[3]?box[3]:box[1],
  ]

  function viewMousedown(e: MouseEvent) {
    if (e.button !== 0) return
    e.stopPropagation()
    e.preventDefault()
    let b = mainElement.getBoundingClientRect()
    showBox = true
    box[0] = box[2] = e.clientX - b.left + mainElement.scrollLeft
    box[1] = box[3] = e.clientY - b.top + mainElement.scrollTop
    box = [...box]
    let viewMousemove = (e: MouseEvent) => {
      box[2] = e.clientX - b.left + mainElement.scrollLeft
      box[3] = e.clientY - b.top + mainElement.scrollTop
      box = [...box]
    }
    let viewMouseup = (e: MouseEvent) => {
      let nodes = mainElement.getElementsByTagName('LI')
      let els = Array.from(nodes).filter(v=>v.getAttribute('data-file-id'))
      let hits = []
      for (let el of els) {
        let elb = el.getBoundingClientRect()
        let x = elb.left - b.left + mainElement.scrollLeft
        let y = elb.top - b.top + mainElement.scrollTop
        let isOverlapping = (x1min: number, x1max: number, x2min: number, x2max: number, y1min: number, y1max: number, y2min: number, y2max: number) => (x1min < x2max && x2min < x1max && y1min < y2max && y2min < y1max)
        // Check for constrain as well as contain
        if (isOverlapping(x, x+elb.width, boxRepresentation[0], boxRepresentation[2], y, y+elb.height, boxRepresentation[1], boxRepresentation[3])) {
          hits.push(el.getAttribute('data-file-id'))
        }
      }
      // TODO: Shift to add to selection, Control to remove
      if (e.shiftKey) {
        selectedFiles = [...new Set([...hits,...selectedFiles])]
        focusedFile = hits[hits.length-1]
      } else if (e.ctrlKey) {
        selectedFiles = selectedFiles.filter(v=>!hits.includes(v))
        if (!selectedFiles.includes(focusedFile)) {
          focusedFile = selectedFiles[selectedFiles.length-1]
        }
      } else {
        selectedFiles = hits
        focusedFile = hits[hits.length-1]
      }
      // Find our closest file to our mouse coordinate to use for focus.
      let closest = 99999
      for (let el of els) {
        let attr = el.getAttribute('data-file-id')
        if (!selectedFiles.includes(attr)) continue
        let elb = el.getBoundingClientRect()
        let x = elb.left - b.left + mainElement.scrollLeft
        let y = elb.top - b.top + mainElement.scrollTop
        let x1 = x + elb.width/2
        let y1 = y + elb.height/2
        let d = Math.sqrt(Math.pow(box[2]-x1, 2)+Math.pow(box[3]-y1, 2))
        if (d < closest) {
          closest = d
          focusedFile = attr
        }
      }

      view.select(focusedFile, selectedFiles)

      showBox = false
      window.removeEventListener('mouseup', viewMouseup)
      window.removeEventListener('mousemove', viewMousemove)
    }
    window.addEventListener('mouseup', viewMouseup)
    window.addEventListener('mousemove', viewMousemove)
  }
</script>

<main bind:this={mainElement} on:mousedown={viewMousedown}>
  <nav>
    {#if !crumbs[0] || crumbs[0][0] !== ''}
      <li on:click={async ()=>await travel("/")}></li>
    {/if}
    {#each crumbs as crumb}
      <li class:focused={$view.wd===crumb[1]} on:click={async ()=>await travel("/"+crumb[1])} title={crumb[1]}>{crumb[0]}</li>
    {/each}
  </nav>
  <section>
    {#if $view.wd != ""}
      <li on:click={async ()=>await travel("..")}>
        <div class="item folder">
          <span class='icon'>folder</span>
          <span class='title'>
            ..
          </span>
        </div>
      </li>
    {/if}
    {#each folders as [key, folder] (key)}
      <li data-folder-id={key} on:click={async ()=>await travel(key)}>
        <div class="item folder">
          <span>folder</span>
          <span class='title'>
            {key}
          </span>
        </div>
      </li>
    {/each}
    {#each files as [key, file] (key)}
      <ViewFile key={key} file={file} directory={directory} view={view} />
    {/each}
    {#if showBox}
      <div class='box' style="left: {boxRepresentation[0]}px; top: {boxRepresentation[1]}px; width: {boxRepresentation[2]-boxRepresentation[0]}px; height: {boxRepresentation[3]-boxRepresentation[1]}px"></div>
    {/if}
  </section>
</main>

<style>
  main {
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    grid-template-rows: auto minmax(0, 1fr);
    overflow: hidden;
  }
  nav {
    display: flex;
    align-items: flex-start;
  }
  nav li {
    position: relative;
    display: block;
    height: 1.5em;
    background: var(--primary-darker);
    padding-left: 1em;
  }
  nav li:after {
    display: block;
    content: ' ';
    width: 0;
    height: 0;
    position: absolute;
    bottom: 0;
    left: 100%;
    border-style: solid;
    border-width: .75em 0em .75em 1em;
    border-color: transparent transparent transparent var(--primary-darker);
    z-index: 2;
  }
  nav li:nth-child(even) {
    background: var(--primary-dark);
  }
  nav li:nth-child(even):after {
    border-color: transparent transparent transparent var(--primary-dark);
  }
  nav li.focused {
    background: var(--primary);
  }
  nav li.focused:after {
    border-color: transparent transparent transparent var(--primary);
  }
  section {
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
    border: 1px solid transparent;
  }
  li .title {
    border: 1px solid transparent;
  }
  li.selected .title {
    background: var(--primary-light);
  }
  li.focused .title {
    border: 1px solid var(--secondary-light);
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
  .item.folder {
    background: green;
  }
  .item.file {
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
  .box {
    position: absolute;
    background: var(--primary-light);
    opacity: 0.5;
    pointer-events: none;
  }
</style>