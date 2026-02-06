<script lang="ts">
  export let title: string;
  export let value: string;
  export let subtitle: string = '';
  export let color: string = 'text-blue-500';
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg p-4 shadow">
  <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">{title}</h3>
  <p class="text-2xl font-bold {color} mt-1">{value}</p>
  {#if subtitle}
    <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">{subtitle}</p>
  {/if}
</div>
