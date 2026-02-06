import { writable } from 'svelte/store';
import { invoke } from '@tauri-apps/api/core';
import type { DivoomDevice } from '@lib/types';

interface DevicesState {
  devices: DivoomDevice[];
  loading: boolean;
  error: string | null;
}

function createDevicesStore() {
  const { subscribe, set, update } = writable<DevicesState>({
    devices: [],
    loading: false,
    error: null,
  });

  return {
    subscribe,
    discover: async () => {
      update(state => ({ ...state, loading: true, error: null }));

      try {
        const devices = await invoke<DivoomDevice[]>('discover_devices_command');
        update(state => ({ ...state, devices, loading: false }));
      } catch (error) {
        update(state => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to discover devices'
        }));
      }
    },
  };
}

export const devices = createDevicesStore();
