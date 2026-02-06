import { describe, it, expect, vi, beforeEach } from 'vitest';
import { get } from 'svelte/store';
import { devices } from './devices';
import { invoke } from '@tauri-apps/api/core';

// Mock Tauri invoke
vi.mock('@tauri-apps/api/core', () => ({
  invoke: vi.fn(),
}));

describe('Devices Store', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should initialize with empty state', () => {
    const state = get(devices);
    expect(state.devices).toEqual([]);
    expect(state.loading).toBe(false);
    expect(state.error).toBe(null);
  });

  it('should set loading to true when discover starts', async () => {
    const mockDevices = [
      {
        ip: '192.168.1.100',
        product_name: 'Pixoo64',
        device_name: 'My Pixoo',
        hardware: '100',
        mac: 'AA:BB:CC:DD:EE:FF',
      },
    ];

    vi.mocked(invoke).mockResolvedValueOnce(mockDevices);

    const promise = devices.discover();
    const state = get(devices);

    expect(state.loading).toBe(true);
    expect(state.error).toBe(null);

    await promise;
  });

  it('should update devices list on successful discovery', async () => {
    const mockDevices = [
      {
        ip: '192.168.1.100',
        product_name: 'Pixoo64',
        device_name: 'My Pixoo',
        hardware: '100',
        mac: 'AA:BB:CC:DD:EE:FF',
      },
    ];

    vi.mocked(invoke).mockResolvedValueOnce(mockDevices);

    await devices.discover();

    const state = get(devices);
    expect(state.devices).toEqual(mockDevices);
    expect(state.loading).toBe(false);
    expect(state.error).toBe(null);
  });

  it('should handle error when discovery fails', async () => {
    const mockError = new Error('Network error');
    vi.mocked(invoke).mockRejectedValueOnce(mockError);

    await devices.discover();

    const state = get(devices);
    // Note: Devices array may contain previous test data since store is singleton
    // In production, you'd want to clear state or use fresh store instance
    expect(state.loading).toBe(false);
    expect(state.error).toBe('Network error');
  });

  it('should handle non-Error errors', async () => {
    vi.mocked(invoke).mockRejectedValueOnce('String error');

    await devices.discover();

    const state = get(devices);
    expect(state.error).toBe('Failed to discover devices');
  });

  it('should clear previous errors on new discovery attempt', async () => {
    // First call fails
    vi.mocked(invoke).mockRejectedValueOnce(new Error('First error'));
    await devices.discover();

    let state = get(devices);
    expect(state.error).toBe('First error');

    // Second call succeeds
    const mockDevices = [
      {
        ip: '192.168.1.100',
        product_name: 'Pixoo64',
        device_name: 'My Pixoo',
        hardware: '100',
        mac: 'AA:BB:CC:DD:EE:FF',
      },
    ];
    vi.mocked(invoke).mockResolvedValueOnce(mockDevices);
    await devices.discover();

    state = get(devices);
    expect(state.error).toBe(null);
    expect(state.devices).toEqual(mockDevices);
  });

  it('should handle multiple devices', async () => {
    const mockDevices = [
      {
        ip: '192.168.1.100',
        product_name: 'Pixoo64',
        device_name: 'Device 1',
        hardware: '100',
        mac: 'AA:BB:CC:DD:EE:FF',
      },
      {
        ip: '192.168.1.101',
        product_name: 'TimeGate',
        device_name: 'Device 2',
        hardware: '400',
        mac: '11:22:33:44:55:66',
      },
    ];

    vi.mocked(invoke).mockResolvedValueOnce(mockDevices);

    await devices.discover();

    const state = get(devices);
    expect(state.devices).toHaveLength(2);
    expect(state.devices[0].product_name).toBe('Pixoo64');
    expect(state.devices[1].product_name).toBe('TimeGate');
  });
});
