<script lang='ts'>
  import { getContext } from 'svelte'
  import { POPUPS } from './Menus.svelte'
  import { clickOutside } from './clicker'

  export let popup: string = ''
  export let subpopup: string = ''
  const { currentPopup, openPopup, currentSubPopup, openSubPopup } = getContext(POPUPS)
</script>

{#if popup && $currentPopup.name === popup}
  <nav style="left: {$currentPopup.x}px; top: {$currentPopup.y}px;" use:clickOutside on:outclick={()=>openPopup('', 0, 0)}>
    <slot></slot>
  </nav>
{:else if subpopup && $currentSubPopup.name === subpopup}
  <nav style="left: {$currentSubPopup.x}px; top: {$currentSubPopup.y}px;" use:clickOutside on:outclick={()=>openSubPopup('', 0, 0)}>
    <slot></slot>
  </nav>
{/if}

<style>
  nav {
    position: fixed;
    display: grid;
    grid-template-columns: minmax(0, 1fr);
  }
</style>