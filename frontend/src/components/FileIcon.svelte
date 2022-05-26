<script type='ts'>
  import { onMount } from 'svelte'
  import type { lib, xdgicons } from '../../wailsjs/go/models'
  import { GetIcon } from '../../wailsjs/go/xdgicons/Theme';
  import { settings } from '../stores/settings'

  import mime from 'mime'
  import { getThumbnail } from '../models/thumbnails'
  import Throbber from './Throbber.svelte'

  export let paths: string[]
  let mimetype: string = ''
  let icon: xdgicons.Icon = {}
  let iconMimetype: string = ''
  let thumbnail: lib.Thumbnail = null
  let error: Error

  onMount(async () => {
    mimetype = mime.getType(paths[paths.length-1]) || 'application/octet-stream'
    try {
      icon = await GetIcon("mimetypes/"+mimetype.replace("/","-"), 64, 1)
      iconMimetype = mime.getType(icon.Ext)
    } catch(e) {
      console.log('getIcon error', e)
      // For now...
      icon = await GetIcon("mimetypes/application-octet-stream", 64, 1)
      iconMimetype = mime.getType(".svg")
    }
    try {
      thumbnail = await getThumbnail(paths, { MaxWidth: $settings.thumbnailWidth, MaxHeight: $settings.thumbnailHeight, Method: $settings.thumbnailMethod })
    } catch(e) {
      error = e
      console.log('getThumbnail error', e)
    }
  })
</script>

{#if thumbnail}
  {#if thumbnail.Format}
    <img src="data:{thumbnail.Format};base64,{thumbnail.Bytes}" alt="{thumbnail.Format} thumbnail">
  {:else}
    <img src="data:{iconMimetype};base64,{icon.Bytes}" alt="{iconMimetype} thumbnail">
  {/if}
{:else if error}
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