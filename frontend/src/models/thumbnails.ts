import type { lib } from '../../wailsjs/go/models'

import { GenerateThumbnail } from '../../wailsjs/go/main/WApp'

let thumbnailCache: {[key: string]: lib.Thumbnail} = {}

export async function getThumbnail(paths: string[], opts: lib.ThumbnailOptions): Promise<lib.Thumbnail> {
  let path = paths.join('/')
  if (thumbnailCache[path]) {
    return thumbnailCache[path]
  }
  let thumbnail = await GenerateThumbnail(paths, opts)
  if (thumbnail instanceof Error) {
    let err = thumbnail
    thumbnail = {
      Format: 'unknown',
      Bytes: [],
    }
    thumbnailCache[path] = thumbnail
    //throw err
    return thumbnail
  }

  thumbnailCache[path] = thumbnail
  return thumbnail
}
