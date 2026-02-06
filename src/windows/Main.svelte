<script lang="ts">
  import { onMount } from 'svelte';
  import { metrics } from '@stores/metrics';
  import MetricCard from '@components/MetricCard.svelte';
  import StatusIndicator from '@components/StatusIndicator.svelte';
  import { formatPercent, formatBytes } from '@utils/formatters';

  let connected = false;
  let intervalId: number | null = null;

  onMount(() => {
    // Refresh metrics every second
    metrics.refresh();
    intervalId = setInterval(() => {
      metrics.refresh();
    }, 1000);

    return () => {
      if (intervalId) {
        clearInterval(intervalId);
      }
    };
  });
</script>

<div class="main-window p-4 bg-gray-100 dark:bg-gray-900 min-h-screen">
  <header class="flex justify-between items-center mb-4">
    <h1 class="text-xl font-bold">Divoom PC Companion</h1>
    <StatusIndicator {connected} />
  </header>

  {#if $metrics.loading}
    <p class="text-gray-500">Loading metrics...</p>
  {:else if $metrics.error}
    <p class="text-red-500">{$metrics.error}</p>
  {:else if $metrics.metrics}
    <div class="grid grid-cols-2 gap-3">
      <MetricCard
        title="CPU"
        value={formatPercent($metrics.metrics.cpu.usage_percent)}
        subtitle={$metrics.metrics.cpu.temperature
          ? `Temp: ${$metrics.metrics.cpu.temperature.toFixed(1)}°C`
          : ''}
        color="text-blue-500"
      />
      <MetricCard
        title="GPU"
        value={formatPercent($metrics.metrics.gpu.usage_percent)}
        subtitle={$metrics.metrics.gpu.temperature
          ? `Temp: ${$metrics.metrics.gpu.temperature.toFixed(1)}°C`
          : ''}
        color="text-green-500"
      />
      <MetricCard
        title="Memory"
        value={formatPercent($metrics.metrics.memory.usage_percent)}
        subtitle={`${formatBytes($metrics.metrics.memory.used_gb)} / ${formatBytes(
          $metrics.metrics.memory.total_gb
        )}`}
        color="text-purple-500"
      />
      <MetricCard
        title="Storage"
        value={formatPercent($metrics.metrics.storage.usage_percent)}
        subtitle={`${formatBytes($metrics.metrics.storage.used_gb)} / ${formatBytes(
          $metrics.metrics.storage.total_gb
        )}`}
        color="text-orange-500"
      />
    </div>
  {/if}
</div>
