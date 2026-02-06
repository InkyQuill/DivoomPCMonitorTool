<script lang="ts">
  import type { DivoomDevice } from '@lib/types';

  export let devices: DivoomDevice[] = [];
  export let loading = false;
  export let onSelect: (device: DivoomDevice) => void;

  function handleSelect(device: DivoomDevice) {
    onSelect(device);
  }
</script>

<div class="device-list">
  {#if loading}
    <p class="text-gray-500">Searching for devices...</p>
  {:else if devices.length === 0}
    <p class="text-gray-500">No devices found</p>
  {:else}
    <ul class="space-y-2">
      {#each devices as device (device.ip)}
        <li
          class="bg-white dark:bg-gray-800 rounded-lg p-3 shadow cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700"
          on:click={() => handleSelect(device)}
          role="button"
          tabindex="0"
        >
          <h4 class="font-semibold">{device.device_name || device.product_name}</h4>
          <p class="text-sm text-gray-500">{device.ip}</p>
          <p class="text-xs text-gray-400">Hardware: {device.hardware}</p>
        </li>
      {/each}
    </ul>
  {/if}
</div>
