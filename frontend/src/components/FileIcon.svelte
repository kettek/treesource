<script type='ts'>
  import { onMount } from 'svelte'
  import type { lib } from '../../wailsjs/go/models'

  import { getIcon } from '../models/icons'
  import Throbber from './Throbber.svelte'

  export let paths: string[]
  let icon: lib.Icon = null
  let error: Error

  onMount(async () => {
    try {
      icon = await getIcon(paths, { MaxWidth: 256, MaxHeight: 256, Method: "ApproxBiLinear" })
    } catch(e) {
      error = e
      console.log(e)
    }
  })
</script>

{#if icon}
  {#if icon.Format}
    <img src="data:{icon.Format};base64,{icon.Bytes}" alt="{icon.Format} icon">
  {:else}
    <span>DUNNO</span>
  {/if}
{:else if error}
  <span>FAIL</span>
{:else}
  <Throbber/>
{/if}

<style>
  img {
    max-width: 100%;
    max-height: 100%;
  }
</style>