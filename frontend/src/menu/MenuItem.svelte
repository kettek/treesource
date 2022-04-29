<script lang='ts'>
  import { createEventDispatcher, getContext } from 'svelte'
  import { POPUPS } from './Menus.svelte'
  import { actionPublisher } from '../actions'

  const dispatch = createEventDispatcher()

  const { openPopup, closePopup, openSubPopup } = getContext(POPUPS)

  export let popup = ''
  export let subpopup = ''
  export let action = ''
  export let args: any = undefined
  export let disabled = false
  let self: Element

  function click(e: MouseEvent) {
    if (e.target !== self) return

    if (action) {
      actionPublisher.publish(action, args)
      closePopup()
      return
    }

    let rect = self.getBoundingClientRect()
    if (popup) {
      openPopup(popup, rect.x, rect.y + rect.height)
    } else if (subpopup) {
      openSubPopup(subpopup, rect.x + rect.width, rect.y)
    } else {
      dispatch('click', e)
    }
  }
</script>

<button class:disabled={disabled} bind:this={self} on:click={click}>
  <slot></slot>
</button>

<style>
  button {
    position: relative;
    margin: 0;
    border: 0;
    border-right: 1px solid teal;
    background: #666;
    text-align: left;
  }
  button.disabled {
    pointer-events: none;
  }
</style>