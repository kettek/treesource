export class DirectoryView {
  uuid: number[] | string
  directory: number[]
  wd: string

  constructor(o: any) {
    if (o.uuid) {
      this.uuid = o.uuid
    }
    if (o.directory) {
      this.directory = o.directory
    }
    if (o.wd) {
      this.wd = o.Wd
    } else {
      this.wd = ""
    }
  }
}

export class TagsView {
  uuid: number[] | string
  tags: string[]
  constructor(o: any) {
    if (o.uuid) {
      this.uuid = o.uuid
    }
    if (o.tags) {
      this.tags = o.tags
    }
  }
}