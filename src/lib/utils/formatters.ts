import type { SystemMetrics } from '@lib/types';

export function formatPercent(value: number, decimals: number = 1): string {
  return `${value.toFixed(decimals)}%`;
}

export function formatBytes(gb: number, decimals: number = 2): string {
  if (gb < 1 && gb > 0) {
    return `${(gb * 1024).toFixed(decimals)} MB`;
  }
  return `${gb.toFixed(decimals)} GB`;
}

export function formatTemperature(temp: number | null): string {
  if (temp === null) return 'N/A';
  return `${temp.toFixed(1)}°C`;
}

export function formatMetrics(metrics: SystemMetrics): string {
  const cpu = formatPercent(metrics.cpu.usage_percent);
  const gpu = formatPercent(metrics.gpu.usage_percent);
  const ram = formatPercent(metrics.memory.usage_percent);

  return `CPU: ${cpu} | GPU: ${gpu} | RAM: ${ram}`;
}

export function getMetricsArray(metrics: SystemMetrics): string[] {
  return [
    metrics.cpu.usage_percent.toFixed(0),
    metrics.gpu.usage_percent.toFixed(0),
    metrics.cpu.temperature?.toFixed(0) || '0',
    metrics.gpu.temperature?.toFixed(0) || '0',
    metrics.memory.usage_percent.toFixed(0),
    metrics.storage.usage_percent.toFixed(0),
  ];
}
