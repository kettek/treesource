<script lang="ts">
  import { HasProject, NewProject, GetProject, CloseProjectFile, LoadProjectFile, AddProjectDirectory, RemoveProjectDirectory, Ready, Undoable, Redoable, Undo, Redo, Unsaved, SaveProject, UpdateEntry } from '../wailsjs/go/main/WApp.js'
  import { AddDirectoryView, RemoveDirectoryView, AddTagsView, RemoveTagsView, SelectView, NavigateDirectoryView, SelectViewFiles } from '../wailsjs/go/main/WApp.js'
  import { EventsOnMultiple, Quit } from '../wailsjs/runtime/runtime'
  import { lib } from '../wailsjs/go/models'
  import * as Dialog from '../wailsjs/go/main/Dialog'
  import SplitPane from './components/SplitPane.svelte'
  import { actionPublisher } from './actions'
  import { onMount } from 'svelte'
  import Menu from './sections/Menu.svelte'
  import Directories from './sections/Directories.svelte'
  import { DirectoryView } from './models/views'
  import Views from './sections/Views.svelte'
  import FilePreview from './sections/FilePreview.svelte'
  import FileMetadata from './sections/FileMetadata.svelte'
  import { settings } from './stores/settings'
  import { directories as directoriesStore } from './stores/directories'
  import { views as viewsStore } from './stores/views'

  let project: lib.Project
  let undoable: boolean = false
  let redoable: boolean = false
  let unsaved: boolean = false

  onMount(async () => {
    try {
      await settings.load()
    } catch(err: any) {
      try {
        await settings.save()
      } catch(err: any) {
        console.log('settings', err)
      }
    }
    console.log('loaded settings', $settings)

    let subs = []

    subs.push(actionPublisher.subscribe('undo', async ({sourceTopic, message}) => {
      Undo()
      await refresh()
    }))
    subs.push(actionPublisher.subscribe('redo', async ({sourceTopic, message}) => {
      Redo()
      await refresh()
    }))

    
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
    subs.push(actionPublisher.subscribe('file-save', async ({sourceTopic, message}) => {
      let result = await SaveProject(false)
      await refresh()
    }))
    subs.push(actionPublisher.subscribe('file-close', async ({sourceTopic}) => {
      let error = await CloseProjectFile(false)
      if (error) {
        throw(error)
      }
    }))
    subs.push(actionPublisher.subscribe('quit', async ({message}) => {
      if (!message) {
        if (project && unsaved) {
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
    // Directory-related actions.
    subs.push(actionPublisher.subscribe('directory-add', async({message}) => {
      let dir = ''
      let ignoreDot = true
      console.log('directory-add', message)
      let result = await Dialog.OpenDirectory({
        Title: "Choose Folder to add",
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

      console.log('calling add with', dir, ignoreDot)

      AddProjectDirectory(dir, ignoreDot)
    }))
    subs.push(actionPublisher.subscribe('directory-remove', async({message}) => {
      RemoveProjectDirectory(message)
    }))

    // Views
    subs.push(actionPublisher.subscribe('view-directory-add', async({message}) => {
      AddDirectoryView(message)
    }))
    subs.push(actionPublisher.subscribe('view-directory-remove', async({message}) => {
      RemoveDirectoryView(message)
    }))
    subs.push(actionPublisher.subscribe('view-directory-navigate', async({message}) => {
      NavigateDirectoryView(message.uuid, message.path)
    }))
    subs.push(actionPublisher.subscribe('view-tags-add', async({message}) => {
      AddTagsView(message)
    }))
    subs.push(actionPublisher.subscribe('view-tags-remove', async({message}) => {
      RemoveTagsView(message)
    }))
    subs.push(actionPublisher.subscribe('view-select', async({message}) => {
      SelectView(message)
    }))
    subs.push(actionPublisher.subscribe('view-select-files', async({message}) => {
      await SelectViewFiles(message.uuid, message.selected, message.focused)
    }))
    subs.push(actionPublisher.subscribe('entry-set-rating', async({message}) => {
      console.log('entry-set-rating', message)
      await UpdateEntry(message.uuid, message.path, message.entry)
    }))

    async function refresh() {
      undoable = await Undoable()
      redoable = await Redoable()
      unsaved = await Unsaved()
      console.log('refresh', undoable, redoable)
    }
    
    // Set up runtime event receival.
    EventsOnMultiple('project-load', async (data: any) => {
      directoriesStore.clear()
      viewsStore.clear()
      project = await GetProject()
      console.log("project load", project.Directories)
      await refresh()
    }, -1)

    //
    EventsOnMultiple('project-unload', async (data: any) => {
      console.log("project unload", data)
      project = undefined
      directoriesStore.clear()
      viewsStore.clear()
      await refresh()
    }, -1)

    EventsOnMultiple('project-changed', (data: boolean) => {
    }, -1)

    EventsOnMultiple('directories', (data: any) => {
      console.log('directories', data)
    }, -1)
    EventsOnMultiple('directory', async (data: any) => {
      directoriesStore.addDirectory(new lib.Directory(data))
      await refresh()
    }, -1)
    EventsOnMultiple('directory-add', async (data: any) => {
      directoriesStore.addDirectory(new lib.Directory(data))
      await refresh()
    }, -1)
    EventsOnMultiple('directory-remove', async (data: any) => {
      directoriesStore.removeByUUID(data.UUID)
      await refresh()
    }, -1)
    EventsOnMultiple('directory-sync', async (data: any) => {
      console.log('directory-sync', data)
      await refresh()
    }, -1)
    EventsOnMultiple('directory-synced', async (data: any) => {
      console.log('directory-synced', data)
      await refresh()
    }, -1)
    EventsOnMultiple('directory-entry', (data: any) => {
      let ds = directoriesStore.getByUUID(data.UUID)
      if (ds) {
        ds.addEntry(new lib.DirectoryEntry(data.Entry))
      }
    }, -1)
    EventsOnMultiple('directory-entry-add', async (data: any) => {
      /*let d = directories.find(v=>v.UUID===data.UUID)
      if (!d) return
      let e = new lib.DirectoryEntry(data.Entry)
      if (!e.Tags) {
        e.Tags = []
      }
      if (e.Missing === undefined) {
        e.Missing = false
      }
      d.Entries.push(e)
      directories = [...directories]
      await refresh()*/
    }, -1)
    EventsOnMultiple('directory-entry-update', async (data: any) => {
      let ds = directoriesStore.getByUUID(data.UUID)
      if (ds) {
        let e = ds.getByPath(data.Entry.Path)
        if (e) {
          e.set(new lib.DirectoryEntry(data.Entry))
        }
      }
      await refresh()
    }, -1)
    EventsOnMultiple('directory-entry-remove', async (data: any) => {
      let ds = directoriesStore.getByUUID(data.UUID)
      if (ds) {
        ds.removeByPath(data.Entry.Path)
      }
      await refresh()
    }, -1)
    EventsOnMultiple('directory-entry-missing', async (data: any) => {
      console.log('entry-missing', data)
      await refresh()
    }, -1)
    EventsOnMultiple('directory-entry-found', async (data: any) => {
      console.log('entry-found', data)
      await refresh()
    }, -1)
    // View stuff
    EventsOnMultiple('view-directory-add', async (data: any) => {
      viewsStore.add(new DirectoryView(data.View))
    }, -1)
    EventsOnMultiple('view-directory-remove', async (data: any) => {
      viewsStore.remove(data.View.uuid)
    }, -1)
    EventsOnMultiple('view-directory-navigate', async (data: any) => {
      let vs = viewsStore.get(data.UUID)
      if (vs) {
        vs._setWorkingDir(data.Path)
      }
    }, -1)
    EventsOnMultiple('view-tags-add', async (data: any) => {
      // TODO
    }, -1)
    EventsOnMultiple('view-tags-remove', async (data: any) => {
      // TODO
    }, -1)
    EventsOnMultiple('view-select', async (data: any) => {
      viewsStore.select(data.UUID)
    }, -1)
    EventsOnMultiple('view-select-files', async (data: any) => {
      let vs = viewsStore.get(data.UUID)
      if (vs) {
        vs._selectFiles(data.Focused, data.Selected)
      }
    }, -1)

    Ready()

    return () => {
      subs.forEach(v=>actionPublisher.unsubscribe(v))
    }
  })
</script>

<main>
  <section class='menu'>
    <Menu project={project} unsaved={unsaved} undoable={undoable} redoable={redoable}></Menu>
  </section>
  <section class='view'>
    <SplitPane type="horizontal" pos={80}>
      <section slot=a>
        <SplitPane type="horizontal" pos={20}>
          <section slot=a class='view__dirs'>
            <SplitPane type="vertical" pos={50}>
              <Directories slot=a disabled={!project}></Directories>
              <div slot=b class='view__dirs__tags'>
                tags
              </div>
            </SplitPane>
          </section>
          <section slot=b class='view__view'>
            <Views></Views>
            <div class='view__view__controls'>controls</div>
          </section>
        </SplitPane>
      </section>
      <section slot=b class='view__info'>
        <SplitPane type='vertical' pos={40}>
          <FilePreview view={$viewsStore.selected} slot=a></FilePreview>
          <FileMetadata slot=b view={$viewsStore.selected}></FileMetadata>
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
    background: var(--secondary-dark);
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
    grid-template-rows: minmax(0, 1fr) auto;
    background: var(--neutral-dark);
  }
  section.view__info {
    display: grid;
    grid-template-rows: minmax(0, 1fr);
  }
</style>
