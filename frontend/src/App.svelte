<script lang="ts">
  import folderIcon from './assets/breeze-icons/icons/places/64/folder.svg'
  import noneIcon from './assets/breeze-icons/icons/mimetypes/64/none.svg'
  import { ListFiles } from '../wailsjs/go/lib/App.js'
  import type { lib } from '../wailsjs/go/models'
  import MenuBar from './menu/MenuBar.svelte'
  import MenuItem from './menu/MenuItem.svelte'
  import MenuList from './menu/MenuList.svelte'
  import Menus from './menu/Menus.svelte'
  import MenuSplit from './menu/MenuSplit.svelte'
  import { actionPublisher } from './actions'
  import { onMount } from 'svelte'

  let currentFiles: (lib.DirEntry[]|Error) = []
  let path: string

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
    
    subs.push(actionPublisher.subscribe('file-new', async ({sourceTopic}) => {
      console.log('TODO: Create file')
    }))
    subs.push(actionPublisher.subscribe('file-open', async ({sourceTopic, message}) => {
      if (message === undefined) {
        console.log('TODO: Open file dialog')
      } else {
        console.log('TODO: Open file', message)
      }
    }))
    subs.push(actionPublisher.subscribe('file-close', async ({sourceTopic}) => {
      console.log('TODO: Close file')
    }))

    return () => {
      subs.forEach(v=>actionPublisher.unsubscribe(v))
    }
  })
</script>

<main>
  <section class='menu'>
    <Menus>
      <MenuBar>
        <MenuItem popup='file-menu'>
          File
          <MenuList popup='file-menu'>
            <MenuItem action='file-new'>
              New Treesource File
            </MenuItem>
            <MenuSplit />
            <MenuItem action='file-open'>
              Open Treesource File
            </MenuItem>
            <MenuItem subpopup='file-menu-open-recent'>
              Open Recent...
            </MenuItem>
            <MenuList subpopup='file-menu-open-recent'>
              <MenuItem action='file-open' args='some/file.trsrc'>
                HOT DOG
              </MenuItem>
            </MenuList>
            <MenuSplit />
            <MenuItem action='file-close' disabled>
              Close
            </MenuItem>
          </MenuList>
        </MenuItem>
        <MenuItem popup='help-menu'>
          Help
          <MenuList popup='help-menu'>
            <MenuItem disabled>Get Started</MenuItem>
            <MenuSplit />
            <MenuItem disabled>View License</MenuItem>
            <MenuSplit />
            <MenuItem disabled>Check For Updates</MenuItem>
            <MenuSplit />
            <MenuItem disabled >About</MenuItem>
          </MenuList>
        </MenuItem>
      </MenuBar>
    </Menus>
  </section>
  <section class='view'>
    <section class='view__dirs'>
      <div class='view__dirs__dirs'>
        dirs
      </div>
      <div class='view__dirs__tags'>
        tags
      </div>
    </section>
    <section class='view__view'>
      <div class='view__view__top'>tabs & viewmenu</div>
      <div class='view__view__items'>files</div>
      <div class='view__view__controls'>controls</div>
    </section>
    <section class='view__info'>
      <div class='view__info__preview'>preview</div>
      <div class='view__info__data'>metadata</div>
    </section>
  </section>
</main>

<style>
  main {
    width: 100%;
    height: 100%;
    display: grid;
    grid-template-rows: auto minmax(0, 1fr);
    grid-template-columns: minmax(0, 1fr);
  }
  section.menu {
  }
  section.view {
    display: grid;
    grid-template-rows: minmax(0, 1fr);
    grid-template-columns: auto minmax(0, 1fr) auto;
  }
  section.view__dirs {
    display: grid;
    grid-template-rows: minmax(0, 1fr) minmax(0, 1fr);
  }
  section.view__view {
    display: grid;
    grid-template-rows: auto minmax(0, 1fr) auto;
  }
  section.view__info {
    display: grid;
    grid-template-rows: minmax(0, 1fr) auto;
  }
</style>
