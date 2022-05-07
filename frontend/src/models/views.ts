export class DirectoryView {
  uuid: number[] | string
  directory: number[]
  wd: string
  selected: string[]

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
    if (o.selected) {
      this.selected = o.selected
    } else {
      this.selected = []
    }
  }
}

export class TagsView {
  uuid: number[] | string
  tags: string[]
  selected: string[]
  constructor(o: any) {
    if (o.uuid) {
      this.uuid = o.uuid
    }
    if (o.tags) {
      this.tags = o.tags
    }
    if (o.selected) {
      this.selected = o.selected
    } else {
      this.selected = []
    }
  }
}