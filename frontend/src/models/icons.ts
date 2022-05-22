import type { lib } from '../../wailsjs/go/models'

import { GenerateIcon } from '../../wailsjs/go/main/WApp'

let iconCache: {[key: string]: lib.Icon} = {}

export async function getIcon(paths: string[], opts: lib.IconOptions): Promise<lib.Icon> {
  let path = paths.join('/')
  if (iconCache[path]) {
    return iconCache[path]
  }
  let icon = await GenerateIcon(paths, opts)
  if (icon instanceof Error) {
    let err = icon
    icon = {
      Format: 'unknown',
      Bytes: [],
    }
    iconCache[path] = icon
    //throw err
    return icon
  }

  iconCache[path] = icon
  return icon
}