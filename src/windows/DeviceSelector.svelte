<script lang="ts">
  import { devices } from '@stores/devices';
  import DeviceList from '@components/DeviceList.svelte';
  import type { DivoomDevice } from '@lib/types';

  export let onSelect: (device: DivoomDevice) => void;

  function handleDiscover() {
    devices.discover();
  }

  function handleSelect(device: DivoomDevice) {
    onSelect(device);
  }
</script>

<div class="device-selector p-6 bg-white dark:bg-gray-800">
  <div class="flex justify-between items-center mb-4">
    <h2 class="text-xl font-bold">Select Device</h2>
    <button
      on:click={handleDiscover}
      class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      disabled={$devices.loading}
    >
      {$devices.loading ? 'Searching...' : 'Discover'}
    </button>
  </div>

  <DeviceList
    devices={$devices.devices}
    loading={$devices.loading}
    onSelect={handleSelect}
  />
</div>
