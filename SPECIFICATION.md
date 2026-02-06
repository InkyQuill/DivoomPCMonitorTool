# Divoom PC Companion - Full Specification

## Project Overview

**Divoom PC Companion** is a cross-platform desktop application for monitoring and displaying PC system metrics on Divoom devices (Pixoo64, TimeGate, etc.) over the local network. This is a complete rewrite of the existing DivoomPCMonitor using Rust + Tauri architecture.

**Existing implementations:**
- Windows: .NET Framework 4.8 with LibreHardwareMonitor (C#/WinForms)
- Linux: Python 3 with psutil, lm-sensors, and GPU tools

**Target improvements:**
- Single codebase for all platforms
- Better performance and lower resource usage
- Modern UI with smoother user experience
- Better error handling and recovery
- Smaller application size

---

## Architecture

### Technology Stack

**Backend (Rust):**
- Tauri 2.0 - window management, system tray, auto-start
- sysinfo - cross-platform system metrics
- serde/serde_json - serialization
- tokio - async runtime
- reqwest - HTTP client
- NVML wrapper (optional) - NVIDIA GPU monitoring

**Frontend (Svelte):**
- Svelte 5 - reactive components
- Svelte-i18n - internationalization
- Tailwind CSS - styling
- Chart.js or similar - optional charts

### Project Structure

```
divoom-pc-companion/
├── src-tauri/                  # Rust backend
│   ├── src/
│   │   ├── main.rs            # Entry point
│   │   ├── lib.rs             # Library exports
│   │   ├── core/              # Core functionality
│   │   │   ├── mod.rs
│   │   │   ├── errors.rs      # Error types
│   │   │   �   └── metrics.rs  # Metric data structures
│   │   ├── system/            # System monitoring
│   │   │   ├── mod.rs
│   │   │   ├── collector.rs   # Main collector
│   │   │   ├── cpu.rs
│   │   │   ├── memory.rs
│   │   │   ├── gpu/
│   │   │   │   ├── mod.rs
│   │   │   │   ├── nvidia.rs
│   │   │   │   ├── amd.rs
│   │   │   │   ├── intel.rs
│   │   │   │   └── generic.rs
│   │   │   └── storage.rs
│   │   ├── devices/           # Device communication
│   │   │   ├── mod.rs
│   │   │   ├── discovery.rs   # Device discovery
│   │   │   ├── protocol.rs    # Divoom protocol
│   │   │   ├── models.rs      # Device types
│   │   │   └── sender.rs      # Data sender
│   │   ├── config/            # Configuration management
│   │   │   ├── mod.rs
│   │   │   ├── schema.rs
│   │   │   ├── storage.rs
│   │   │   └── defaults.rs
│   │   ├── ui/                # Window management
│   │   │   ├── mod.rs
│   │   │   ├── tray.rs
│   │   │   ├── windows.rs
│   │   │   └── autostart.rs
│   │   └── i18n/              # Internationalization
│   │       ├── mod.rs
│   │       ├── detector.rs
│   │       └── loader.rs
│   ├── Cargo.toml
│   ├── tauri.conf.json
│   └── build.rs
├── src/                        # Svelte frontend
│   ├── lib/
│   │   ├── components/        # Reusable components
│   │   │   ├── MetricCard.svelte
│   │   │   ├── DeviceList.svelte
│   │   │   └── StatusIndicator.svelte
│   │   ├── stores/            # Svelte stores
│   │   │   ├── metrics.ts
│   │   │   ├── devices.ts
│   │   │   └── config.ts
│   │   ├── i18n/
│   │   │   └── locales/
│   │   │       ├── en.json
│   │   │       ├── ru.json
│   │   │       ├── zh.json
│   │   │       └── ...
│   │   └── utils/
│   │       ├── formatters.ts
│   │       └── constants.ts
│   ├── windows/               # Window components
│   │   ├── Main.svelte        # Main widget window
│   │   ├── Settings.svelte
│   │   └── DeviceSelector.svelte
│   ├── App.svelte
│   └── main.ts
├── public/
│   ├── icons/
│   └── styles/
├── locales/                   # Backend locale files
├── docs/
└── package.json
```

---

## Divoom Protocol

### Discovery

**Endpoint:** `GET http://app.divoom-gz.com/Device/ReturnSameLANDevice`

**Response:**
```json
{
  "DeviceList": [
    {
      "DeviceId": 12345678,
      "Hardware": 400,
      "DeviceName": "TimeGate",
      "DevicePrivateIP": "192.168.1.100",
      "DeviceMac": "AA:BB:CC:DD:EE:FF"
    }
  ]
}
```

**Device Types:**
- `Hardware == 400` - TimeGate (5 screens)
- Other values - Pixoo64 and similar (1 screen)

### TimeGate Screen Selection

For TimeGate devices (Hardware 400), additional API calls are needed:

**Get 5-LCD Info:**
```
GET http://app.divoom-gz.com/Channel/Get5LcdInfoV2?DeviceType=LCD&DeviceId={DeviceId}
```

**Response:**
```json
{
  "LcdIndependence": 1,
  "ChannelType": 5,
  "ClockId": 625
}
```

**Select Screen:**
```
POST http://{device_ip}:80/post
Content-Type: application/json

{
  "LcdIndependence": 1,
  "Command": "Channel/SetClockSelectId",
  "LcdIndex": 1,
  "ClockId": 625
}
```

### Send System Metrics

**Endpoint:** `POST http://{device_ip}:80/post`

**Request:**
```json
{
  "Command": "Device/UpdatePCParaInfo",
  "ScreenList": [
    {
      "LcdId": 0,
      "DispData": [
        "45%",    // [0] CPU Usage
        "72%",    // [1] GPU Usage
        "55°C",   // [2] CPU Temperature
        "65°C",   // [3] GPU Temperature
        "68%",    // [4] Memory Usage
        "42°C"    // [5] Disk Usage or Temperature
      ]
    }
  ]
}
```

**Important:** The Windows version actually sends disk temperature for index 5, not disk usage. The Linux version sends disk usage. Should be configurable.

**LcdId values:**
- TimeGate: 0-4 (5 screens)
- Other devices: 0 (single screen)

---

## System Monitoring Module

### CPU Metrics

**Collected Data:**
- Total usage percentage
- Per-core usage (optional)
- Temperature (maximum of all cores)
- Frequency (current, max - optional)
- Power consumption (watts - optional)

**Implementation by Platform:**

**Windows:**
```rust
// Primary: LibreHardwareMonitor equivalent (need Rust alternative)
// Alternatives:
// 1. WMI queries (winapi crate)
// 2. Performance Counters (windows crate)
// 3. sysinfo crate (cross-platform)
```

**Linux:**
```rust
// Primary: /proc/stat for usage
// Temperature methods (in order of priority):
// 1. psutil-sensors wrapper
// 2. /sys/class/thermal/thermal_zone*/temp
// 3. lm-sensors via subprocess
```

**macOS:**
```rust
// IOKit/SMC via crate
// Alternative: sysinfo crate
```

**Rust Crates to Consider:**
- `sysinfo` - basic metrics (cross-platform)
- `wmi` (Windows) - WMI queries
- `sensors-sys` - lm-sensors bindings

### GPU Metrics

**Collected Data:**
- Usage percentage
- Temperature
- Memory usage (used/total)
- Core/Memory frequency (optional)

**Vendor-Specific Detection:**

**NVIDIA:**
```rust
// Primary: nvml-wrapper crate (NVML API)
// Fallback: nvidia-smi subprocess parsing
```

**AMD:**
```rust
// Linux: rocm-smi subprocess parsing
// Windows: ADLX SDK (if available) or LibreHardwareMonitor alternative
```

**Intel:**
```rust
// Linux: intel_gpu_top subprocess parsing
// Windows: oneAPI Level Zero or alternative
```

### Memory Metrics

**Collected Data:**
- Total physical memory
- Used memory
- Available memory
- Swap usage (total/used)

**Implementation:**
- `sysinfo` crate provides cross-platform support
- Direct system calls for more details

### Storage Metrics

**Collected Data:**
- Disk usage percentage (primary partition)
- Read/Write speeds (optional)
- Temperature (optional, varies by platform)

**Implementation:**
- `sysinfo` for usage
- Platform-specific for temperature

### Network Metrics (Optional)

**Collected Data:**
- Upload speed
- Download speed
- Total transferred

---

## Configuration Management

### Config File Location

**Windows:** `%APPDATA%\DivoomPCCompanion\config.json`
**Linux:** `~/.config/divoom-pc-companion/config.json`
**macOS:** `~/Library/Application Support/DivoomPCCompanion/config.json`

### Configuration Schema

```json
{
  "version": 1,
  "general": {
    "language": "auto",
    "theme": "system",
    "start_minimized": false,
    "minimize_to_tray": true,
    "check_for_updates": true
  },
  "widget": {
    "position": { "x": 100, "y": 100 },
    "size": { "width": 350, "height": 250 },
    "always_on_top": true,
    "show_title_bar": false,
    "opacity": 100,
    "display_mode": "compact"
  },
  "monitoring": {
    "update_interval": 1000,
    "metrics": {
      "cpu": true,
      "cpu_temp": true,
      "gpu": true,
      "gpu_temp": true,
      "memory": true,
      "disk": true
    },
    "temperature_unit": "celsius",
    "disk_metric_type": "usage"
  },
  "device": {
    "auto_connect": true,
    "auto_discovery": true,
    "selected_device": {
      "id": null,
      "name": null,
      "ip": null,
      "hardware_type": null
    },
    "screen_number": 1,
    "send_interval": 1000,
    "connection_timeout": 5000,
    "retry_attempts": 3,
    "retry_delay": 2000
  },
  "alerts": {
    "enabled": true,
    "cpu_temp_warning": 80,
    "cpu_temp_critical": 90,
    "gpu_temp_warning": 85,
    "gpu_temp_critical": 95,
    "show_notifications": true
  },
  "autostart": {
    "enabled": false,
    "start_minimized": false,
    "delay_seconds": 0
  },
  "advanced": {
    "log_level": "info",
    "debug_mode": false
  }
}
```

---

## Windows/Views Specification

### 1. Main Widget Window

**Purpose:** Display real-time system metrics in a compact, draggable widget.

**Characteristics:**
- Frameless window (no title bar)
- Transparent background (optional)
- Draggable by mouse
- Always on top (optional, configurable)
- Click-through when minimized (optional)
- Resizable (within min/max bounds)
- Minimize to tray on close

**Dimensions:**
- Minimum: 300x200px
- Default: 350x250px
- Maximum: 600x400px

**Layout:**
```
┌─────────────────────────────────────┐
│ [≡]             Divoom PC Monitor   │ <- Optional title bar
├─────────────────────────────────────┤
│                                     │
│  ┌─────────────┐  ┌─────────────┐   │
│  │   CPU: 45%  │  │   GPU: 72%  │   │
│  │   55°C      │  │   65°C      │   │
│  └─────────────┘  └─────────────┘   │
│                                     │
│  ┌─────────────┐  ┌─────────────┐   │
│  │   RAM: 68%  │  │   HDD: 42%  │   │
│  │  16/24 GB   │  │   250 GB    │   │
│  └─────────────┘  └─────────────┘   │
│                                     │
│  Status: ● Connected to TimeGate    │
└─────────────────────────────────────┘
```

**Components:**
- **MetricCard** - Reusable card for each metric
  - Icon
  - Label
  - Value (large text)
  - Secondary value (optional)
  - Color indicator (green/yellow/red based on thresholds)

**Context Menu (Right-click):**
- Show/Hide
- Settings
- Select Device
- Refresh Device List
- Disconnect
- About
- Exit

### 2. Settings Window

**Purpose:** Configure all application settings.

**Structure:** Tabbed interface with categories.

**Tabs:**

#### General Tab
- Language dropdown (auto + supported languages)
- Theme dropdown (light/dark/system)
- Start minimized checkbox
- Minimize to tray checkbox
- Check for updates checkbox

#### Widget Tab
- Always on top checkbox
- Show title bar checkbox
- Opacity slider (50-100%)
- Display mode dropdown (compact/detailed/graphs)
- Reset position button

#### Metrics Tab
- Update interval slider (1-10 seconds)
- Enable/disable metrics:
  - [x] CPU Usage
  - [x] CPU Temperature
  - [x] GPU Usage
  - [x] GPU Temperature
  - [x] Memory Usage
  - [x] Disk Usage/Temperature
- Temperature unit (Celsius/Fahrenheit)
- Disk metric type (Usage/Temperature)
- Alert thresholds:
  - CPU Temp Warning: __°C
  - CPU Temp Critical: __°C
  - GPU Temp Warning: __°C
  - GPU Temp Critical: __°C

#### Device Tab
- Auto-connect on startup
- Auto-discovery interval
- Send interval (1-5 seconds)
- Connection timeout
- Retry attempts
- Advanced (expandable):
  - Display number selector (1-5 for TimeGate)
  - Brightness (if supported)
  - Test connection button

#### Autostart Tab
- Enable autostart checkbox
- Start mode: Normal/Minimized
- Start with elevated privileges (Windows only)
- Delay startup (seconds)

#### About Tab
- Application icon and name
- Version: x.x.x
- Links:
  - Documentation
  - GitHub Repository
  - Issue Tracker
  - License
- Check for updates button

### 3. Device Selector Window

**Purpose:** Select and configure Divoom device.

**Layout:**
```
┌──────────────────────────────────────────────┐
│  Select Device                [Refresh]      │
├──────────────────────────────────────────────┤
│                                              │
│  Searching for devices...                    │
│                                              │
│  ┌────────────────────────────────────────┐  │
│  │ TimeGate                    192.168.1.100│  │
│  │ Hardware: 400          [Connect]       │  │
│  └────────────────────────────────────────┘  │
│  ┌────────────────────────────────────────┐  │
│  │ Pixoo64                     192.168.1.101│  │
│  │ Hardware: 204          [Connect]       │  │
│  └────────────────────────────────────────┘  │
│                                              │
│              [Cancel]            [Auto Select]│
└──────────────────────────────────────────────┘
```

**Features:**
- Auto-refresh every 30 seconds
- Manual refresh button
- Auto-select if only one device
- Show connection status
- For TimeGate: show screen selector

**TimeGate Screen Selector:**
```
┌──────────────────────────────────────────────┐
│  Select Display Screen (TimeGate)            │
├──────────────────────────────────────────────┤
│                                              │
│  [1]  [2]  [3]  [4]  [5]                     │
│   ●    ○    ○    ○    ○                      │
│                                              │
│              [Cancel]       [Connect]        │
└──────────────────────────────────────────────┘
```

---

## System Tray

### Tray Icon

**States:**
- Normal: Blue/green icon
- Connected: Green with indicator
- Disconnected: Orange/yellow
- Error: Red

**Tooltip:**
```
Divoom PC Monitor
Status: Connected to TimeGate
CPU: 45% (55°C)
GPU: 72% (65°C)
```

### Context Menu Items

- **Show/Hide** - Toggle main window
- **--- Separator ---**
- **Quick Settings** (submenu)
  - Always on Top
  - Update Interval (1s, 2s, 5s, 10s)
- **--- Separator ---**
- **Device** (submenu)
  - Refresh Devices
  - Disconnect
  - Select Device...
- **--- Separator ---**
- **Settings...**
- **About...**
- **--- Separator ---**
- **Exit**

---

## Auto-start Implementation

### Windows

**Method:** Task Scheduler (preferred for elevated privileges)

```rust
use windows::Win32::System::TaskScheduler::*;

// Create task with:
// - Trigger: Logon
// - Action: Execute application
// - Principal: Run with highest privileges
```

**Alternative:** Registry Startup folder (non-elevated)
```
HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run
```

### Linux

**Method:** .desktop file in autostart

```
~/.config/autostart/divoom-pc-companion.desktop
```

```ini
[Desktop Entry]
Type=Application
Name=Divoom PC Companion
Exec=/usr/bin/divoom-pc-companion
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
```

### macOS

**Method:** Login Items via LaunchServices

```rust
use core_foundation::base::TCFType;
use core_foundation::url::CFURL;
```

---

## Internationalization (i18n)

### Supported Languages

Based on existing implementation:
1. English (en) - Default
2. Russian (ru)
3. Chinese Simplified (zh-CN)
4. Chinese Traditional (zh-TW)
5. German (de)
6. French (fr)
7. Japanese (ja)
8. Korean (ko)
9. Spanish (es)
10. Italian (it)
11. Portuguese (pt-BR)
12. Polish (pl)
13. Ukrainian (uk)
14. Dutch (nl)
15. Swedish (sv)
16. Finnish (fi)
17. Czech (cs)
18. Hungarian (hu)
19. Arabic (ar)

### Implementation

**Frontend (Svelte-i18n):**
```javascript
// src/lib/i18n/locales/en.json
{
  "app": {
    "title": "Divoom PC Companion",
    "version": "Version {{version}}"
  },
  "metrics": {
    "cpu_usage": "CPU Usage",
    "cpu_temp": "CPU Temperature",
    "gpu_usage": "GPU Usage",
    "gpu_temp": "GPU Temperature",
    "memory_usage": "Memory Usage",
    "disk_usage": "Disk Usage"
  },
  "device": {
    "connected": "Connected to {{name}}",
    "disconnected": "Disconnected",
    "searching": "Searching for devices..."
  }
}
```

**Backend (Rust):**
```rust
// locales/en.json
{
  "tray": {
    "show": "Show",
    "hide": "Hide",
    "exit": "Exit"
  }
}
```

**Auto-detection:**
```rust
use sys_locale::get_locale;

let locale = get_locale().unwrap_or_else(|| String::from("en"));
```

---

## Error Handling and Logging

### Error Types

```rust
pub enum AppError {
    DeviceDiscoveryFailed,
    DeviceConnectionFailed,
    DataSendFailed,
    MetricCollectionFailed(String),
    ConfigError(ConfigError),
    IoError(std::io::Error),
    HttpError(reqwest::Error),
    JsonError(serde_json::Error),
}
```

### Logging

**Levels:** ERROR, WARN, INFO, DEBUG, TRACE

**Log File Location:**
- Windows: `%APPDATA%\DivoomPCCompanion\logs\`
- Linux: `~/.local/state/divoom-pc-companion/logs/`
- macOS: `~/Library/Logs/DivoomPCCompanion/`

**Rotation:** Daily, keep last 7 days

---

## Notifications

### Toast Notifications (Desktop)

**Triggered by:**
- Device connection lost
- Device reconnected
- High temperature warning
- Data send failure (persistent)
- Update available

**Implementation:**
- Windows: WinRT notifications
- Linux: libnotify
- macOS: NSUserNotification

---

## Build and Release

### Target Platforms

**Windows:**
- Windows 10 (1903+) - x64
- Windows 11 - x64

**Linux:**
- Ubuntu 20.04+ - x64
- Debian 11+ - x64
- Fedora 35+ - x64
- Arch Linux - x64

**macOS:**
- macOS 11+ - x64
- macOS 11+ - aarch64 (Apple Silicon)

### Build Configuration

```toml
# Cargo.toml
[profile.release]
opt-level = "z"     # Optimize for size
lto = true          # Link-time optimization
codegen-units = 1   # Better optimization
strip = true        # Remove symbols
panic = "abort"     # Reduce binary size
```

### Packaging

**Windows:**
- NSIS installer
- Portable ZIP

**Linux:**
- .deb (Debian/Ubuntu)
- .rpm (Fedora/RHEL)
- AppImage (universal)
- AUR package (Arch)

**macOS:**
- .dmg
- .app bundle

---

## Performance Requirements

### Resource Limits

- **RAM Usage:** < 50 MB in background
- **CPU Usage:** < 1% idle, < 5% during update
- **Disk Space:** < 15 MB installed
- **Startup Time:** < 3 seconds
- **Update Interval:** 1-10 seconds (configurable)

### Optimization Strategies

1. **Lazy Loading:** Load GPU monitoring only if available
2. **Caching:** Cache device list and system info
3. **Async Operations:** Non-blocking I/O
4. **Connection Pooling:** Reuse HTTP connections
5. **Batch Updates:** Send metrics in single request

---

## Security Considerations

1. **Network Security:** Only communicate on local network
2. **Input Validation:** Validate all device responses
3. **Error Messages:** Don't expose sensitive information
4. **Permissions:** Request minimal permissions
5. **Code Signing:** Sign executables (Windows/macOS)

---

## Testing Strategy

### Unit Tests

- Metric collection for each platform
- Protocol serialization/deserialization
- Configuration management
- Error handling

### Integration Tests

- Device discovery
- Data transmission
- Auto-start functionality

### Manual Testing Checklist

- [ ] All metrics display correctly
- [ ] Device discovery works
- [ ] Device connection persists
- [ ] Settings save/load correctly
- [ ] Tray menu functions
- [ ] Auto-start works on all platforms
- [ ] Language switching works
- [ ] Error recovery works
- [ ] No memory leaks (24h test)

---

## Migration from Existing Implementations

### Windows (C#) → Rust

**Breaking Changes:**
- Config location changes
- Registry keys removed
- Different UI framework

**Migration Path:**
1. Import existing config on first run
2. Detect and convert Registry settings
3. Prompt user for device re-selection

### Linux (Python) → Rust

**Breaking Changes:**
- Config file format changes
- systemd service needs update

**Migration Path:**
1. Import `~/.divoom_config.json`
2. Update systemd service file
3. Maintain compatibility with existing config

---

## Roadmap

### Phase 1: Core (Week 1-2)
- [ ] Project setup (Tauri + Svelte)
- [ ] Basic metric collection (CPU, Memory)
- [ ] Config management
- [ ] Basic UI framework

### Phase 2: Device Integration (Week 3)
- [ ] Device discovery
- [ ] Protocol implementation
- [ ] Data transmission
- [ ] Error handling

### Phase 3: Complete Metrics (Week 4)
- [ ] GPU monitoring (all vendors)
- [ ] Temperature sensors
- [ ] Storage metrics
- [ ] Network metrics (optional)

### Phase 4: UI Polish (Week 5)
- [ ] Main widget
- [ ] Settings window
- [ ] Device selector
- [ ] System tray
- [ ] Animations

### Phase 5: System Integration (Week 6)
- [ ] Auto-start (all platforms)
- [ ] Internationalization
- [ ] Notifications
- [ ] Update checker

### Phase 6: Testing & Release (Week 7-8)
- [ ] Cross-platform testing
- [ ] Performance optimization
- [ ] Packaging
- [ ] Documentation
- [ ] Release v1.0

---

## Open Questions and Research Needed

1. **Divoom Protocol:**
   - Are there other device types besides 400?
   - What's the ClockId magic number (625)?
   - Are there authentication requirements?

2. **GPU Monitoring:**
   - Is there a cross-platform Rust crate for all GPU vendors?
   - Should we use subprocess parsing or native APIs?

3. **Temperature Monitoring:**
   - What's the most reliable method per platform?
   - How to handle systems without sensors?

4. **UI Framework:**
   - Should we use Svelte or consider alternatives?
   - Do we need charting libraries?

5. **Distribution:**
   - Should we use GitHub Releases or other distribution?
   - Auto-update mechanism?

---

## References

**Existing Implementations Analyzed:**
- Windows: `/Windows/DivoomPCMonitorTool/`
  - C# .NET Framework 4.8
  - LibreHardwareMonitor for metrics
  - Windows Task Scheduler for auto-start
  - Registry for config

- Linux: `/Linux/divoom.py`
  - Python 3 with psutil
  - lm-sensors for CPU temp
  - Vendor-specific GPU tools
  - JSON config file

**Related Projects:**
- LibreHardwareMonitor - Hardware monitoring
- sysinfo - Rust system info crate
- Tauri - Rust desktop app framework

---

## Appendix: Protocol Examples

### Complete Discovery Flow

```bash
# 1. Discover devices
curl http://app.divoom-gz.com/Device/ReturnSameLANDevice

# Response:
{
  "DeviceList": [
    {
      "DeviceId": 12345678,
      "Hardware": 400,
      "DeviceName": "TimeGate",
      "DevicePrivateIP": "192.168.1.100",
      "DeviceMac": "AA:BB:CC:DD:EE:FF"
    }
  ]
}

# 2. For TimeGate, get screen info
curl "http://app.divoom-gz.com/Channel/Get5LcdInfoV2?DeviceType=LCD&DeviceId=12345678"

# Response:
{
  "LcdIndependence": 1,
  "ChannelType": 5,
  "ClockId": 625
}

# 3. Select screen
curl -X POST http://192.168.1.100:80/post \
  -H "Content-Type: application/json" \
  -d '{
    "LcdIndependence": 1,
    "Command": "Channel/SetClockSelectId",
    "LcdIndex": 1,
    "ClockId": 625
  }'

# 4. Send metrics
curl -X POST http://192.168.1.100:80/post \
  -H "Content-Type: application/json" \
  -d '{
    "Command": "Device/UpdatePCParaInfo",
    "ScreenList": [{
      "LcdId": 0,
      "DispData": ["45%", "72%", "55°C", "65°C", "68%", "42°C"]
    }]
  }'
```

---

**Document Version:** 1.0
**Last Updated:** 2025-01-XX
**Status:** Draft
