import type { DirectoryView } from '../models/views'
import { writable, get, Subscriber, Writable } from 'svelte/store'
import type { lib } from '../../wailsjs/go/models'
import type Views from 'src/sections/Views.svelte'
import { actionPublisher } from '../actions'

interface DirectoryViewStore extends Writable<DirectoryView> {
  _setWorkingDir: (wd: string) => void
  _selectFiles: (focused: string, selected: string[]) => void
  select: (focused: string, selected: string[]) => void
}

function createDirectoryViewStore(v: DirectoryView): DirectoryViewStore {
  const { subscribe, set, update } = writable(v)

  return {
    subscribe,
    set,
    update,
    _setWorkingDir: (wd: string) => {
      let v = get({subscribe})
      v.wd = wd
      set(v)
    },
    _selectFiles: (focused: string, selected: string[]) => {
      let v = get({subscribe})
      v.focused = focused
      v.selected = selected
      set(v)
    },
    select: (focused: string, selected: string[]) => {
      let v = get({subscribe})
      actionPublisher.publish('view-select-files', {
        uuid: v.uuid,
        selected: selected,
        focused: focused,
      })
    },
  }
}

interface ViewStoreData {
  views: DirectoryViewStore[]
  selected: DirectoryViewStore
}

interface ViewsStore extends Writable<ViewStoreData> {
  clear: () => void
  add: (d: DirectoryView) => void
  remove: (uuid: number[]|string) => void
  get: (uuid: number[]|string) => DirectoryViewStore
  select: (uuid: number[]|string) => void
}

export const views: ViewsStore = ((): ViewsStore => {
  const viewStoreData: ViewStoreData = {
    views: [],
    selected: null,
  }

  const { subscribe, set, update } = writable(viewStoreData)

  return {
    subscribe,
    set,
    update,
    clear: () => {
      set({
        views: [],
        selected: null,
      })
    },
    add: (d: DirectoryView) => {
      let vs = get({subscribe})
      if (vs.views.find(v=>get(v).uuid===d.uuid)) return
      vs.views.push(createDirectoryViewStore(d))
      set(vs)
    },
    remove: (uuid: number[]|string) => {
      let vs = get({subscribe})
      vs.views = vs.views.filter(v=>get(v).uuid!==uuid)
      set(vs)
    },
    get: (uuid: number[]|string) => {
      let vs = get({subscribe})
      return vs.views.find(v=>get(v).uuid===uuid)
    },
    select: (uuid: number[]|string) => {
      let vs = get({subscribe})
      vs.selected = vs.views.find(v=>get(v).uuid===uuid)
      set(vs)
    },
  }
})()