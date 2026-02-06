import { writable } from 'svelte/store';
import { invoke } from '@tauri-apps/api/core';
import type { SystemMetrics } from '@lib/types';

interface MetricsState {
  metrics: SystemMetrics | null;
  loading: boolean;
  error: string | null;
}

function createMetricsStore() {
  const { subscribe, set, update } = writable<MetricsState>({
    metrics: null,
    loading: false,
    error: null,
  });

  return {
    subscribe,
    refresh: async () => {
      update(state => ({ ...state, loading: true, error: null }));

      try {
        const metrics = await invoke<SystemMetrics>('collect_metrics');
        update(state => ({ ...state, metrics, loading: false }));
      } catch (error) {
        update(state => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to collect metrics'
        }));
      }
    },
  };
}

export const metrics = createMetricsStore();
