# Divoom PC Companion - Project Scaffold

## Overview

A complete Rust + Tauri 2.0 + Svelte 5 project scaffold for cross-platform PC monitoring application targeting Divoom devices (Pixoo64, TimeGate).

## Project Structure

```
DivoomPCMonitorTool/
├── src-tauri/           # Rust backend
│   ├── src/
│   │   ├── core/       # Core types and errors
│   │   ├── system/     # System metrics collection
│   │   ├── devices/    # Divoom device communication
│   │   ├── config/     # Configuration management
│   │   ├── ui/         # UI utilities (tray, windows)
│   │   └── i18n/       # Internationalization
│   ├── Cargo.toml      # Rust dependencies
│   └── tauri.conf.json # Tauri configuration
├── src/                # Frontend
│   ├── lib/
│   │   ├── components/ # Svelte components
│   │   ├── stores/     # Svelte stores
│   │   ├── i18n/       # Translations
│   │   └── utils/      # Utilities
│   ├── windows/        # Window components
│   └── main.ts         # Entry point
├── package.json        # Node dependencies
├── vite.config.ts      # Vite configuration
├── tailwind.config.js  # Tailwind CSS
└── tsconfig.json       # TypeScript config
```

## Implemented Features

### Rust Backend ✅
- **Core Module**: Error types, metrics data structures
- **Config Module**: Load/save config, platform-specific paths
- **Device Module**: Discovery, protocol models, data sending
- **System Module**: Basic metrics collection (CPU, Memory, Storage)
  - GPU monitoring stubs (TODO: vendor-specific)
  - Temperature detection stubs (TODO: platform-specific)
- **Tauri Commands**: discover_devices, send_metrics, get_config, save_config, collect_metrics

### Frontend ✅
- **Svelte 5 + TypeScript**: Full type safety
- **Tailwind CSS**: Utility-first styling
- **Svelte Stores**: metrics, devices, config
- **Components**: MetricCard, DeviceList, StatusIndicator
- **Windows**: Main (metrics display), Settings, DeviceSelector
- **i18n**: English, Russian, Chinese locales

## Status

- ✅ Cargo check passes
- ✅ npm run check passes
- ⏳ npm run tauri:dev - Ready to test
- ⏳ npm run tauri:build - Ready to build

## Next Steps

### Immediate (Post-Scaffold)
1. Run `npm run tauri:dev` to verify the application works
2. Create proper icons (currently using placeholder blue squares)
3. Test device discovery on local network
4. Implement system tray

### Core Features (TODO)
1. **CPU Temperature**: Platform-specific detection
   - Linux: Read from `/sys/class/thermal/`
   - Windows: WMI queries
   - macOS: IOKit

2. **GPU Monitoring**: Vendor-specific implementations
   - NVIDIA: NVML
   - AMD: ADLX
   - Intel: oneAPI Level Zero

3. **Settings UI**: Complete all settings tabs
4. **Autostart**: Platform-specific implementations
5. **Error Handling**: User-friendly error notifications
6. **Performance**: Optimize metrics collection interval

### Advanced Features
1. **Additional Languages**: 13 more languages (total 16)
2. **Tray Menu**: Quick actions, status indicator
3. **Device Reconnection**: Auto-reconnect with retry
4. **Multiple Devices**: Support monitoring multiple devices
5. **Custom Themes**: User-defined color schemes
6. **Data Export**: Metrics history export

## Development Commands

```bash
# Install dependencies
npm install

# Development
npm run tauri:dev

# Type checking
npm run check
cd src-tauri && cargo check

# Production build
npm run tauri:build

# Output: src-tauri/target/release/bundle/
```

## Configuration

Config file location (auto-created on first run):
- **Linux**: `~/.config/divoom-pc-companion/config.json`
- **macOS**: `~/Library/Application Support/DivoomPCCompanion/config.json`
- **Windows**: `%APPDATA%\DivoomPCCompanion\config.json`

## API Reference

### Tauri Commands

```typescript
// Discover Divoom devices on LAN
const devices = await invoke<DivoomDevice[]>('discover_devices_command');

// Send metrics to device
await invoke('send_metrics_command', {
  ip: '192.168.1.100',
  lcd_id: 0,
  metrics: ['50', '60', '45', '70', '75', '50']
});

// Get current configuration
const config = await invoke<AppConfig>('get_config');

// Save configuration
await invoke('save_config_command', { config: newConfig });

// Collect system metrics
const metrics = await invoke<SystemMetrics>('collect_metrics');
```

## Device Protocol

### Discovery
```
GET http://app.divoom-gz.com/Device/ReturnSameLANDevice
```

### Send Metrics
```
POST http://{device_ip}:80/post
{
  "Command": "Device/UpdatePCParaInfo",
  "TextSec": 20,
  "LcdId": 0,
  "TextContent": ["CPU%", "GPU%", "CPU°", "GPU°", "RAM%", "Disk%"]
}
```

## Requirements Met

- ✅ Tauri 2.0 with Svelte 5 + TypeScript
- ✅ Cross-platform (Linux, macOS, Windows)
- ✅ sysinfo for basic metrics
- ✅ reqwest for HTTP client
- ✅ serde/serde_json for serialization
- ✅ Tailwind CSS for styling
- ✅ Stub implementations for GPU/temperature
- ✅ i18n support (3/16 languages)
- ✅ Frameless window 350x250, always on top
- ✅ Configuration management
- ✅ Device discovery
- ✅ Metrics collection every second

## License

[To be determined]
