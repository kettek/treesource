<script lang='ts'>
  // Imports
  import { get } from 'svelte/store'
  import { actionPublisher } from '../actions'

  // Types
  import type { DirectoryViewStore } from '../stores/views'

  // Stores
  import { directories as directoriesStore } from '../stores/directories'

  // Properties
  export let view: DirectoryViewStore

  // Reactive Vars
  $: directory = directoriesStore.getByUUID($view?.directory)
  $: focusedEntry = $directory?.Entries.find(v=>get(v).Path===$view.focused)

  // Functions
  function removeTag(tag: string) {
    actionPublisher.publish('entry-remove-tag', {
      uuid: $view.uuid,
      tag: tag,
    })
  }

  function setRating(value: number) {
    actionPublisher.publish('entry-set-rating', {
      uuid: $directory.RealDir.UUID,
      path: $focusedEntry.Path,
      entry: {...$focusedEntry, Rating: value},
    })
  }
</script>

<main>
  {#if focusedEntry}
    <section>
      <header>rating</header>
      <article class='rating'>
        {#each [1, 2, 3, 4, 5] as rating}
          <span class='rating__star' class:selected={rating<=$focusedEntry.Rating} on:click={_=>setRating(rating)}>
            {rating<=$focusedEntry.Rating?'✦':'✧'}
          </span>
        {/each}
      </article>
    </section>
    <section>
      <header>tags</header>
      <article class='tags'>
        {#if $focusedEntry.Tags}
          {#each $focusedEntry.Tags as tag}
            <div class='tag'>
              <div class='tag__name'>
                {tag}
              </div>
              <div class='tag__remove' on:click={_=>removeTag(tag)}> x </div>
            </div>
          {/each}
        {/if}
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
  .rating__star {
    cursor: pointer;
  }
  .rating__star.selected {
    color: yellow;
  }
</style>