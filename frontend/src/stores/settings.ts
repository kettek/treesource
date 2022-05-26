import { SaveSettings, LoadSettings } from '../../wailsjs/go/main/WApp'
import { writable, get } from 'svelte/store'

interface Settings {
  thumbnailWidth: number
  thumbnailHeight: number
  thumbnailMethod: 'CatmullRom' | 'NearestNeighbor' | 'ApproxBiLinear'
  autoplayAudio: boolean
  autoplayVideo: boolean
}

const DefaultSettings: Settings = {
  thumbnailWidth: 200,
  thumbnailHeight: 200,
  thumbnailMethod: 'NearestNeighbor',
  autoplayAudio: true,
  autoplayVideo: false,
}

function createSettings() {
  const { subscribe, set, update } = writable({...DefaultSettings})

  return {
    subscribe,
    reset: () => set({...DefaultSettings}),
    json: () => JSON.stringify(get({subscribe})),
    save: () => SaveSettings(get({subscribe})),
    load: () => LoadSettings().then(r=>set({...DefaultSettings, ...r})),
  }
}

export const settings = createSettings()