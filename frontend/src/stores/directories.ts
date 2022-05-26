import { writable, get, Subscriber, Writable } from 'svelte/store'
import type { lib } from '../../wailsjs/go/models'
import * as ftt from '@kettek/filepaths-to-tree'

export interface DirectoryEntryStore extends Writable<lib.DirectoryEntry> {
}

function createEntryStore(f: lib.DirectoryEntry): DirectoryEntryStore {
  const { subscribe, set, update } = writable(f)

  return {
    subscribe,
    set,
    update,
  }
}

export interface DirectoryStoreData {
  Entries: DirectoryEntryStore[]
  Tree: any
  RealDir: lib.Directory
}

export interface DirectoryStore extends Writable<DirectoryStoreData> {
  addEntry: (e: lib.DirectoryEntry) => void
  getByPath: (p: string) => DirectoryEntryStore
  removeByPath: (p: string) => void
}

function createDirectoryStore(d: lib.Directory): DirectoryStore {
  const dir: DirectoryStoreData = {
    Entries: [],
    Tree: {},
    RealDir: d,
  }

  const { subscribe, set, update } = writable(dir)

  return {
    subscribe,
    set,
    update,
    getByPath: (p: string) => {
      let dir = get({subscribe})
      return dir.Entries.find(v=>get(v).Path===p)
    },
    removeByPath: (p: string) => {
      let dir = get({subscribe})
      dir.Entries = dir.Entries.filter(v=>get(v).Path!==p)
      ftt.Remove(dir.Tree, p)
      set(dir)
    },
    addEntry: (e: lib.DirectoryEntry) => {
      let dir = get({subscribe})
      if (dir.Entries.find(v=>get(v).Path===e.Path)) {
        return
      }
      dir.Entries.push(createEntryStore(e))
      ftt.Insert(dir.Tree, e.Path, dir.Entries[dir.Entries.length-1])
      set(dir)
    }
  }
}

export interface DirectoriesStoreData {
  Directories: DirectoryStore[]
  Selected: DirectoryStore 
}

export interface DirectoriesStore extends Writable<DirectoriesStoreData> {
  clear: () => void
  select: (uuid: string|number[]) => void
  addDirectory: (e: lib.Directory) => void
  getByUUID: (uuid: string|number[]) => DirectoryStore
  removeByUUID: (uuid: string|number[]) => void
}

export const directories: DirectoriesStore = ((): DirectoriesStore => {
  const dirstore: DirectoriesStoreData = {
    Directories: [],
    Selected: null,
  }
  const { subscribe, set, update } = writable(dirstore)

  return {
    subscribe,
    set,
    update,
    clear: () => {
      set({
        Directories: [],
        Selected: null,
      })
    },
    select: (uuid: string|number[]) => {
      let ds = get({subscribe})
      ds.Selected = ds.Directories.find(v=>get(v).RealDir.UUID === uuid)
      set(ds)
    },
    addDirectory: (e: lib.Directory) => {
      let ds = get({subscribe})
      if (ds.Directories.find(v=>get(v).RealDir.UUID === e.UUID)) {
        return
      }
      ds.Directories.push(createDirectoryStore(e))
      set(ds)
    },
    getByUUID: (uuid: string|number[]) => {
      let ds = get({subscribe})
      return ds.Directories.find(v=>get(v).RealDir.UUID === uuid)
    },
    removeByUUID: (uuid: string|number[]) => {
      let ds = get({subscribe})
      ds.Directories = ds.Directories.filter(v=>get(v).RealDir.UUID !== uuid)
      set(ds)
    },
  }
})()