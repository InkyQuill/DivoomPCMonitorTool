import { writable, get } from 'svelte/store';
import { invoke } from '@tauri-apps/api/core';
import type { AppConfig } from '@lib/types';

interface ConfigState {
  config: AppConfig | null;
  loading: boolean;
  error: string | null;
}

function createConfigStore() {
  const { subscribe, set, update } = writable<ConfigState>({
    config: null,
    loading: false,
    error: null,
  });

  return {
    subscribe,
    load: async () => {
      update((state) => ({ ...state, loading: true, error: null }));

      try {
        const config = await invoke<AppConfig>('get_config');
        update((state) => ({ ...state, config, loading: false }));
      } catch (error) {
        update((state) => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to load config',
        }));
      }
    },
    save: async (config: AppConfig) => {
      update((state) => ({ ...state, loading: true, error: null }));

      try {
        await invoke('save_config_command', { config });
        update((state) => ({ ...state, config, loading: false }));
      } catch (error) {
        update((state) => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to save config',
        }));
      }
    },
  };
}

export const config = createConfigStore();
