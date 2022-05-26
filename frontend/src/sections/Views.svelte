<script lang='ts'>
  import { views as viewsStore } from '../stores/views'
  import ViewsTab from './ViewsTab.svelte'
  import ViewsContent from './ViewsContent.svelte'

</script>

<main>
  <nav>menu of some sort</nav>
  <nav class='tabs'>
    {#each $viewsStore.views as tab}
      <ViewsTab selected={$viewsStore.selected===tab} view={tab}/>
    {/each}
  </nav>
  <section>
    {#if $viewsStore.selected}
      <ViewsContent view={$viewsStore.selected}/>
    {/if}
  </section>
</main>

<style>
  main {
    display: grid;
    grid-template-rows: auto auto minmax(0, 1fr);
  }
  section {
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    grid-template-rows: minmax(0, 1fr);
  }
  nav.tabs {
    display: flex;
    align-items: center;
    justify-content: flex-start;
    background: var(--secondary-dark);
    grid-column-gap: .5em;
    height: 2em;
  }
</style>