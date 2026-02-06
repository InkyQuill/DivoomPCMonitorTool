<script lang="ts">
  import { config } from '@stores/config';
  import { SUPPORTED_LANGUAGES, UPDATE_INTERVALS } from '@utils/constants';

  async function handleSave() {
    if ($config.config) {
      await config.save($config.config);
    }
  }
</script>

<div class="settings-window p-6 bg-white dark:bg-gray-800 min-h-screen">
  <h1 class="text-2xl font-bold mb-6">Settings</h1>

  {#if $config.config}
    <div class="space-y-6">
      <!-- General Settings -->
      <section>
        <h2 class="text-lg font-semibold mb-3">General</h2>

        <div class="space-y-3">
          <div>
            <label class="block text-sm font-medium mb-1">Language</label>
            <select
              class="w-full p-2 border rounded dark:bg-gray-700 dark:border-gray-600"
              bind:value={$config.config.general.language}
            >
              {#each SUPPORTED_LANGUAGES as lang}
                <option value={lang.code}>{lang.name}</option>
              {/each}
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium mb-1">Update Interval</label>
            <select
              class="w-full p-2 border rounded dark:bg-gray-700 dark:border-gray-600"
              bind:value={$config.config.general.update_interval}
            >
              {#each UPDATE_INTERVALS as interval}
                <option value={interval.value}>{interval.name}</option>
              {/each}
            </select>
          </div>

          <label class="flex items-center space-x-2">
            <input
              type="checkbox"
              class="rounded"
              bind:checked={$config.config.general.start_minimized}
            />
            <span class="text-sm">Start minimized</span>
          </label>

          <label class="flex items-center space-x-2">
            <input
              type="checkbox"
              class="rounded"
              bind:checked={$config.config.general.autostart}
            />
            <span class="text-sm">Autostart</span>
          </label>
        </div>
      </section>

      <!-- Metrics Settings -->
      <section>
        <h2 class="text-lg font-semibold mb-3">Metrics</h2>

        <div class="space-y-2">
          <label class="flex items-center space-x-2">
            <input
              type="checkbox"
              class="rounded"
              bind:checked={$config.config.metrics.enable_cpu}
            />
            <span class="text-sm">Enable CPU monitoring</span>
          </label>

          <label class="flex items-center space-x-2">
            <input
              type="checkbox"
              class="rounded"
              bind:checked={$config.config.metrics.enable_gpu}
            />
            <span class="text-sm">Enable GPU monitoring</span>
          </label>

          <label class="flex items-center space-x-2">
            <input
              type="checkbox"
              class="rounded"
              bind:checked={$config.config.metrics.enable_memory}
            />
            <span class="text-sm">Enable Memory monitoring</span>
          </label>

          <label class="flex items-center space-x-2">
            <input
              type="checkbox"
              class="rounded"
              bind:checked={$config.config.metrics.enable_storage}
            />
            <span class="text-sm">Enable Storage monitoring</span>
          </label>
        </div>
      </section>
    </div>

    <div class="mt-6 flex space-x-3">
      <button
        on:click={handleSave}
        class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        disabled={$config.loading}
      >
        {$config.loading ? 'Saving...' : 'Save'}
      </button>
    </div>
  {/if}
</div>
