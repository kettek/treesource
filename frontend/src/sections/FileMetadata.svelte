<script lang='ts'>
  import { actionPublisher } from '../actions'

  import type { lib } from '../../wailsjs/go/models'
  import type { DirectoryView, TagsView } from '../models/views'

  export let directories: lib.Directory[] = []

  export let tagsViews: TagsView[]
  export let directoryViews: DirectoryView[]

  export let selectedView: string | number[]  = ''

  $: directoryView = directoryViews.find(v=>v.uuid===selectedView)
  $: directory = directoryView ? directories.find(v => v.UUID === directoryView.directory) : undefined

  $: focusedEntry = directory?.Entries.find(v=>v.Path===directoryView?.focused)

  function removeTag(tag: string) {
    actionPublisher.publish('entry-remove-tag', {
      uuid: selectedView,
      tag: tag,
    })
  }
</script>

<main>
  {#if focusedEntry}
    <section>
      <header>rating</header>
      <article class='rating'>
        * * * *
      </article>
    </section>
    <section>
      <header>tags</header>
      <article class='tags'>
        {#each focusedEntry.Tags as tag}
        <div class='tag'>
          <div class='tag__name'>
            {tag}
          </div>
          <div class='tag__remove' on:click={_=>removeTag(tag)}> x </div>
        </div>
        {/each}
      </article>
    </section>
  {/if}
</main>

<style>
  main {
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    grid-template-rows: auto minmax(0, 1fr);
    overflow: hidden;
  }
  section {
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    grid-template-rows: auto minmax(0, 1fr);
  }
  article.tags {
    display: flex;
    flex-wrap: wrap;
    flex-direction: row;
    align-items: flex-start;
    justify-content: flex-start;
    padding: .5em;
    overflow: auto;
  }
  .tag {
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto;
    grid-template-rows: minmax(0, 1fr);
    border: 1px solid white;
    border-radius: 2px;
  }
  .tag__name {
    border-right: 1px solid white;
  }
  .tag__remove {
    border-left: 1px solid gray;
  }
</style>