<script context='module'>
  export const POPUPS = Symbol()
</script>

<script lang='ts'>
  import { setContext, onDestroy } from 'svelte'
  import { Writable, writable } from 'svelte/store'

  interface Popup {
    name: string;
    x: number;
    y: number;
  }

  const currentPopup: Writable<Popup> = writable({name: '', x: 0, y: 0})
  const currentSubPopup: Writable<Popup> = writable({name: '', x: 0, y: 0})

  setContext(POPUPS, {
    openPopup: (popup: string, x: number, y: number) => {
      x = Math.floor(x)
      y = Math.floor(y)
      currentPopup.set({name: popup, x, y})
      currentSubPopup.set({name: '', x, y})
    },
    currentPopup,
    openSubPopup: (popup: string, x: number, y: number) => {
      x = Math.floor(x)
      y = Math.floor(y)
      currentSubPopup.set({name: popup, x, y})
    },
    currentSubPopup,
  })

</script>

<slot></slot>