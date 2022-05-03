<script lang='ts'>
  import MenuBar from '../menu/MenuBar.svelte'
  import MenuItem from '../menu/MenuItem.svelte'
  import MenuList from '../menu/MenuList.svelte'
  import Menus from '../menu/Menus.svelte'
  import MenuSplit from '../menu/MenuSplit.svelte'
  import type { lib } from '../../wailsjs/go/models'

  export let project: lib.Project
  export let changed: boolean
  export let undoable: boolean
  export let redoable: boolean
</script>

<Menus>
  <MenuBar>
    <MenuItem popup='file-menu'>
      File
      <MenuList popup='file-menu'>
        <MenuItem action='file-new'>
          New Treesource Project
        </MenuItem>
        <MenuItem action='file-new' args='withDir'>
          Open Folder as New Project
        </MenuItem>
        <MenuSplit />
        <MenuItem action='file-open'>
          Open Treesource Project
        </MenuItem>
        <MenuItem subpopup='file-menu-open-recent'>
          Open Recent Project...
        </MenuItem>
        <MenuList subpopup='file-menu-open-recent'>
          <MenuItem action='file-open' args='some/file.trsrc'>
            HOT DOG
          </MenuItem>
        </MenuList>
        <MenuSplit />
        <MenuItem action='file-save' disabled={!project || !changed}>
          Save
        </MenuItem>
        <MenuSplit />
        <MenuItem action='file-close' disabled={!project}>
          Close
        </MenuItem>
        <MenuSplit />
        <MenuItem action='quit'>
          Quit
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
    <MenuItem action='undo' disabled={!undoable}>Undo</MenuItem>
    <MenuItem action='redo' disabled={!redoable}>Redo</MenuItem>
  </MenuBar>
</Menus>
