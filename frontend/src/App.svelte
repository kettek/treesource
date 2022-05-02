<script lang="ts">
  import folderIcon from './assets/breeze-icons/icons/places/64/folder.svg'
  import noneIcon from './assets/breeze-icons/icons/mimetypes/64/none.svg'
  import { HasProject, ListFiles, NewProject, GetProject, CloseProjectFile, LoadProjectFile } from '../wailsjs/go/main/WApp.js'
  import { EventsOn, EventsOff, EventsOnMultiple, Quit } from '../wailsjs/runtime/runtime'
  import { lib } from '../wailsjs/go/models'
  import * as Dialog from '../wailsjs/go/main/Dialog'
  import SplitPane from './components/SplitPane.svelte'
  import { actionPublisher } from './actions'
  import { onMount } from 'svelte'
  import Menu from './sections/Menu.svelte'
  import Directories from './sections/Directories.svelte'

  let currentFiles: (lib.DirEntry[]|Error) = []
  let path: string

  let project: lib.Project
  let changed: boolean

  let directories: lib.Directory[] = []

  $: title = project ? project.Title : ''

  $: inRoot = path === '/' || path === ''

  async function getFiles() {
    try {
      currentFiles = await ListFiles(path)
    } catch(err) {
      console.log(err)
    }
  }
  async function navigate(to: lib.DirEntry) {
    if (!to.Info.IsDir) return
    // TODO: Send request to move to a directory, then wait for a response.
    path = path+"/" + to.Name
    await getFiles()
  }
  async function up() {
    // TODO: Send request to move up a directory, then wait for a response.
  }

  onMount(() => {
    let subs = []
    
    subs.push(actionPublisher.subscribe('file-new', async ({sourceTopic, message}) => {
      if (await HasProject()) return

      let dir: string = ''
      let ignoreDot: boolean = true
      if (message === 'withDir') {
        // Request source directory.
        let result = await Dialog.OpenDirectory({
          Title: "Choose source Folder",
          CanCreateDirectories: true,
        })
        if (result instanceof Error) {
          throw result
        }
        if (result === '') {
          // Canceled
          return
        }
        dir = result

        result = await Dialog.Message({
          Type: "question",
          Title: "File Parsing",
          Message: "Ignore '.' prefixed files and folders?",
        })
        if (result !== "Yes") {
          ignoreDot = false
        }
      }
      // Request project location.
      let result = await Dialog.SaveFile({
        Title: "Choose a location to save the treesource project",
        DefaultFilename: "project.trsrc",
        CanCreateDirectories: true,
        Filters: [{DisplayName: "treesource", Pattern: "*.trsrc;*.treesource"}]
      })
      if (result instanceof Error) {
        throw result
      }
      if (result === '') {
        // Canceled
        return
      }
      let error = await NewProject(result, dir, ignoreDot)
      if (error !== null) {
        throw error
      }
      // Should be good.
    }))
    subs.push(actionPublisher.subscribe('file-open', async ({sourceTopic, message}) => {
      let result = await Dialog.OpenFile({
        Title: "Open a treesource project",
        Filters: [{DisplayName: "treesource", Pattern: "*.trsrc;*.treesource"}]
      })
      if (result instanceof Error) {
        throw result
      }
      if (result === '') {
        return
      }
      let error = await LoadProjectFile(result, false)
      if (error !== null) {
        throw error
      }
    }))
    subs.push(actionPublisher.subscribe('file-close', async ({sourceTopic}) => {
      let error = await CloseProjectFile(false)
      if (error) {
        throw(error)
      }
    }))
    subs.push(actionPublisher.subscribe('quit', async ({message}) => {
      if (!message) {
        if (project && changed) {
          let result = await Dialog.Message({
            Type: "warning",
            Message: "The project has changes and is not saved. Quit anyway?",
          })
          if (result !== "Yes") {
            return
          }
        }
      }
      Quit()
    }))
    
    // Set up runtime event receival.
    EventsOnMultiple('project-load', async (data: any) => {
      directories = []
      project = await GetProject()
      console.log("project load", project.Directories)
    }, -1)

    //
    EventsOnMultiple('project-unload', (data: any) => {
      console.log("project unload", data)
      project = undefined
      changed = false
      directories = []
    }, -1)

    EventsOnMultiple('project-changed', (data: boolean) => {
      changed = data
    }, -1)

    EventsOnMultiple('directories', (data: any) => {
      console.log('directories', data)
    }, -1)
    EventsOnMultiple('directory', (data: any) => {
      data.Path = data.Name // FIXME
      directories = [...directories, new lib.Directory(data)]
      if (!directories[directories.length-1].Entries) {
        directories[directories.length-1].Entries = []
      }
    }, -1)
    EventsOnMultiple('directory-sync', (data: any) => {
      console.log('directory-sync', data)
    }, -1)
    EventsOnMultiple('directory-synced', (data: any) => {
      console.log('directory-synced', data)
    }, -1)
    EventsOnMultiple('directory-entry', (data: any) => {
      let d = directories.find(v=>v.UUID===data.UUID)
      if (!d) return
      let e = new lib.DirectoryEntry(data.Entry)
      if (!e.Tags) {
        e.Tags = []
      }
      d.Entries.push(e)
      directories = [...directories]
    }, -1)
    EventsOnMultiple('directory-entry-add', (data: any) => {
      console.log('entry-add', data)
    }, -1)
    EventsOnMultiple('directory-entry-remove', (data: any) => {
      console.log('entry-remove', data)
    }, -1)
    EventsOnMultiple('directory-entry-missing', (data: any) => {
      console.log('entry-missing', data)
    }, -1)
    EventsOnMultiple('directory-entry-found', (data: any) => {
      console.log('entry-found', data)
    }, -1)

    return () => {
      subs.forEach(v=>actionPublisher.unsubscribe(v))
    }
  })
</script>

<main>
  <section class='menu'>
    <Menu project={project} changed={changed}></Menu>
  </section>
  <section class='view'>
    <SplitPane type="horizontal" pos=80>
      <section slot=a>
        <SplitPane type="horizontal" pos=20>
          <section slot=a class='view__dirs'>
            <SplitPane type="vertical" pos=50>
              <div slot=a class='view__dirs__dirs'>
                <Directories bind:directories={directories}></Directories>
              </div>
              <div slot=b class='view__dirs__tags'>
                tags
              </div>
            </SplitPane>
          </section>
          <section slot=b class='view__view'>
            <div class='view__view__top'>tabs & viewmenu</div>
            <div class='view__view__items'>files</div>
            <div class='view__view__controls'>controls</div>
          </section>
        </SplitPane>
      </section>
      <section slot=b class='view__info'>
        <SplitPane type="vertical" pos=40>
          <div slot=a class='view__info__preview'>preview</div>
          <div slot=b class='view__info__data'>metadata</div>
        </SplitPane>
      </section>
    </SplitPane>
  </section>
</main>

<style>
  main {
    width: 100%;
    height: 100%;
    display: grid;
    grid-template-rows: auto minmax(0, 1fr);
    grid-template-columns: minmax(0, 1fr);
    background: var(--neutral-dark);
    color: var(--primary-text);
  }
  section.menu {
    background: var(--neutral);
  }
  section.view {
    display: grid;
    grid-template-rows: minmax(0, 1fr);
    grid-template-columns: minmax(0, 1fr);
  }
  section.view__dirs {
    display: grid;
    grid-template-rows: minmax(0, 1fr);
  }
  div.view__dirs__dirs {
    border: 1px solid black;
  }
  div.view__dirs__tags {
    border: 1px solid black;
  }
  section.view__view {
    display: grid;
    grid-template-rows: auto minmax(0, 1fr) auto;
    background: var(--secondary-dark);
  }
  section.view__info {
    display: grid;
    grid-template-rows: minmax(0, 1fr);
    border: 1px solid yellow;
  }
</style>
