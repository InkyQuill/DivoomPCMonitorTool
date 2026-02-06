<script lang="ts">
  export let connected = false;
  export let loading = false;

  let color = connected ? 'bg-green-500' : 'bg-gray-400';
  let pulse = loading ? 'animate-pulse' : '';
</script>

<div
  class="w-3 h-3 rounded-full {color} {pulse}"
  title={connected ? 'Connected' : 'Disconnected'}
></div>
