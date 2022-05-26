<script type='ts'>
  import { onMount } from 'svelte'
  import type { lib, xdgicons } from '../../wailsjs/go/models'
  import { GetIcon } from '../../wailsjs/go/xdgicons/Theme';

  import mime from 'mime'
  import Throbber from './Throbber.svelte'

  let icon: xdgicons.Icon = {}
  let iconMimetype: string = ''

  onMount(async () => {
    try {
      icon = await GetIcon("places/folder", 64, 1)
      iconMimetype = mime.getType(icon.Ext)
    } catch(e) {
      console.log('getIcon error', e)
      // For now...
      icon = await GetIcon("mimetypes/application-octet-stream", 64, 1)
      iconMimetype = mime.getType(".svg")
    }
  })
</script>

{#if icon}
  <img src="data:{iconMimetype};base64,{icon.Bytes}" alt="{iconMimetype} thumbnail">
{:else}
  <Throbber/>
{/if}

<style>
  img {
    max-width: 100%;
    max-height: 100%;
    box-shadow: 1px 1px 3px var(--neutral-darker);
  }
</style>