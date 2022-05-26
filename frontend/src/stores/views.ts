import type { DirectoryView } from '../models/views'
import { writable, get, Subscriber, Writable } from 'svelte/store'
import type { lib } from '../../wailsjs/go/models'
import type Views from 'src/sections/Views.svelte'

interface DirectoryViewStore extends Writable<DirectoryView> {
  setWorkingDir: (wd: string) => void
  selectFiles: (focused: string, selected: string[]) => void
}

function createDirectoryViewStore(v: DirectoryView): DirectoryViewStore {
  const { subscribe, set, update } = writable(v)

  return {
    subscribe,
    set,
    update,
    setWorkingDir: (wd: string) => {
      let v = get({subscribe})
      v.wd = wd
      set(v)
      console.log('set wd', wd)
    },
    selectFiles: (focused: string, selected: string[]) => {
      let v = get({subscribe})
      v.focused = focused
      v.selected = selected
      set(v)
      console.log('select', v)
    },
  }
}

interface ViewsStore extends Writable<DirectoryViewStore[]> {
  clear: () => void
  add: (DirectoryView) => void
  remove: (uuid: number[]|string) => void
  get: (uuid: number[]|string) => DirectoryViewStore
}

export const views: ViewsStore = ((): ViewsStore => {
  const views: DirectoryViewStore[] = []

  const { subscribe, set, update } = writable(views)

  return {
    subscribe,
    set,
    update,
    clear: () => {
      set([])
    },
    add: (d: DirectoryView) => {
      let vs = get({subscribe})
      if (vs.find(v=>get(v).uuid===d.uuid)) return
      vs.push(createDirectoryViewStore(d))
      set(vs)
      console.log('add dv', d)
    },
    remove: (uuid: number[]|string) => {
      let vs = get({subscribe})
      vs = vs.filter(v=>get(v).uuid!==uuid)
      set(vs)
      console.log('remvoe dv')
    },
    get: (uuid: number[]|string) => {
      let vs = get({subscribe})
      return vs.find(v=>get(v).uuid===uuid)
    },
  }
})()