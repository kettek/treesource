<script type='ts'>
  import { onMount } from 'svelte'
  import type { lib } from '../../wailsjs/go/models'

  import { getIcon } from '../models/icons'
  import Throbber from './Throbber.svelte'

  export let paths: string[]
  let icon: lib.Icon = null

  onMount(async () => {
    icon = await getIcon(paths, { MaxWidth: 256, MaxHeight: 256, Method: "ApproxBiLinear" })

    return () => {

    }
  })
</script>

{#if icon}
  {#if icon.Format}
    <img src="data:{icon.Format};base64,{icon.Bytes}" alt="{icon.Format} icon">
  {:else}
    <span>DUNNO</span>
  {/if}
{:else}
  <Throbber/>
{/if}

<style>
  img {
    max-width: 100%;
    max-height: 100%;
  }
</style>