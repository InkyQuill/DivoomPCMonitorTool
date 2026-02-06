---
name: divoom-pc-companion-patterns
description: Coding patterns extracted from Divoom PC Companion rewrite (Rust + Tauri + Svelte)
version: 1.0.0
source: local-git-analysis
analyzed_commits: 8
---

# Divoom PC Companion Patterns

This skill captures the coding patterns, architectural decisions, and conventions used in the Divoom PC Companion project - a cross-platform desktop application for monitoring system metrics on Divoom devices.

## Project Overview

**Technology Stack:**
- **Backend**: Rust + Tauri 2.0
- **Frontend**: Svelte 5 + TypeScript + Tailwind CSS
- **Build**: Vite 5 + cargo
- **Platform**: Windows, Linux, macOS

**Key Features:**
- Real-time system monitoring (CPU, GPU, memory, storage)
- Divoom device discovery and communication
- Cross-platform native desktop app
- Internationalization (19 languages)
- System tray integration
- Auto-start functionality

---

## Code Architecture

### Modular Rust Backend

The backend follows a **domain-driven modular architecture** with clear separation of concerns:

```
src-tauri/src/
├── core/          # Core types and errors
├── system/        # System metrics collection
├── devices/       # Device communication
├── config/        # Configuration management
├── ui/            # Window/tray/autostart
└── i18n/          # Internationalization
```

**Pattern**: Each module has a `mod.rs` that re-exports public API, keeping the main `lib.rs` clean.

**Example**:
```rust
// src/core/mod.rs
pub mod errors;
pub mod metrics;

pub use errors::{AppError, Result};
pub use metrics::{SystemMetrics, CpuMetrics, GpuMetrics, MemoryMetrics, StorageMetrics};
```

### Frontend Structure

**Organized by responsibility**:
```
src/
├── lib/
│   ├── components/    # Reusable UI components
│   ├── stores/        # Svelte stores (state management)
│   ├── utils/         # Utility functions
│   └── types.ts       # Shared TypeScript types
├── windows/           # Window-level components
├── App.svelte         # Root component
└── main.ts            # Entry point
```

**Pattern**: Use **path aliases** for clean imports:
```typescript
// vite.config.ts
resolve: {
  alias: {
    '@': path.resolve(__dirname, './src'),
    '@components': path.resolve(__dirname, './src/lib/components'),
    '@stores': path.resolve(__dirname, './src/lib/stores'),
    '@utils': path.resolve(__dirname, './src/lib/utils'),
    '@windows': path.resolve(__dirname, './src/windows')
  }
}
```

---

## Rust Patterns

### Error Handling with `thiserror`

**Use `thiserror` for structured error types**:

```rust
use thiserror::Error;

#[derive(Error, Debug)]
pub enum AppError {
    #[error("Device discovery failed")]
    DeviceDiscoveryFailed,

    #[error("Device connection failed: {0}")]
    DeviceConnectionFailed(String),

    #[error("Metric collection failed: {0}")]
    MetricCollectionFailed(String),

    #[error("Configuration error: {0}")]
    ConfigError(String),

    #[error("IO error: {0}")]
    IoError(#[from] std::io::Error),

    #[error("HTTP error: {0}")]
    HttpError(#[from] reqwest::Error),

    #[error("JSON error: {0}")]
    JsonError(#[from] serde_json::Error),
}

pub type Result<T> = std::result::Result<T, AppError>;
```

**Benefits**:
- Automatic `From` implementations for wrapped errors
- Clean error messages with `Display`
- Type-safe error handling

### Serde for Configuration

**All config and data structures use `serde`**:

```rust
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AppConfig {
    pub general: GeneralConfig,
    pub metrics: MetricsConfig,
    pub device: DeviceConfig,
}

// JSON file storage
pub fn load_config() -> Result<AppConfig> {
    let content = fs::read_to_string(&config_path)?;
    let config: AppConfig = serde_json::from_str(&content)?;
    Ok(config)
}
```

**Pattern**: Use nested structs to group related settings:
```rust
pub struct AppConfig {
    pub general: GeneralConfig,      // Language, autostart
    pub metrics: MetricsConfig,      // What to monitor
    pub device: DeviceConfig,        // Connection settings
}
```

### Tauri Commands Pattern

**Commands follow this structure**:

```rust
#[tauri::command]
async fn command_name(
    input: String,
    state: State<'_, AppState>,
) -> Result<ResponseType, String> {
    // Validate input
    if input.is_empty() {
        return Err("Input cannot be empty".to_string());
    }

    // Do work
    do_something().map_err(|e| e.to_string())?;

    // Return result
    Ok(ResponseType { /* ... */ })
}
```

**Pattern**: Always convert errors to `String` for Tauri compatibility:
```rust
.map_err(|e| e.to_string())?
```

### Async with `tokio`

**Use async I/O for network operations**:

```rust
pub async fn send_metrics(device_ip: &str, lcd_id: i32, metrics: Vec<String>) -> Result<()> {
    let client = Client::new();
    let url = format!("http://{}:80/post", device_ip);

    client
        .post(&url)
        .timeout(std::time::Duration::from_secs(5))
        .json(&request)
        .send()
        .await
        .map_err(|e| AppError::DataSendFailed(e.to_string()))?;

    Ok(())
}
```

**Best practice**: Always set timeouts on HTTP requests.

### System Metrics Collection

**Use the `sysinfo` crate for cross-platform metrics**:

```rust
use sysinfo::System;

pub struct MetricsCollector {
    system: System,
}

impl MetricsCollector {
    pub fn new() -> Self {
        let mut system = System::new_all();
        system.refresh_all();
        Self { system }
    }

    pub fn collect(&mut self) -> Result<SystemMetrics> {
        self.system.refresh_all();

        let cpu_usage = self.system.global_cpu_usage();
        let total_memory = self.system.total_memory();
        let used_memory = self.system.used_memory();

        Ok(SystemMetrics { /* ... */ })
    }
}
```

**Pattern**: Refresh all metrics at once, then extract individual values.

---

## TypeScript/Svelte Patterns

### Svelte Stores for State Management

**Use custom stores with actions**:

```typescript
import { writable } from 'svelte/store';
import { invoke } from '@tauri-apps/api/core';

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
      update(state => ({ ...state, loading: true, error: null }));

      try {
        const config = await invoke<AppConfig>('get_config');
        update(state => ({ ...state, config, loading: false }));
      } catch (error) {
        update(state => ({
          ...state,
          loading: false,
          error: error instanceof Error ? error.message : 'Failed to load config'
        }));
      }
    },
    save: async (config: AppConfig) => {
      // Similar pattern
    }
  };
}

export const config = createConfigStore();
```

**Pattern**: Store state includes `loading` and `error` fields for async operations.

### Type Safety with TypeScript

**All Tauri commands are typed**:

```typescript
import type { SystemMetrics } from '@lib/types';

// Typed invoke
const metrics = await invoke<SystemMetrics>('collect_metrics');
```

**Shared types in `types.ts`**:
```typescript
export interface SystemMetrics {
  cpu: CpuMetrics;
  gpu: GpuMetrics;
  memory: MemoryMetrics;
  storage: StorageMetrics;
  timestamp: number;
}
```

### Svelte Component Props

**Use `export let` for props with defaults**:

```svelte
<script lang="ts">
  export let title: string;
  export let value: string;
  export let subtitle: string = '';
  export let color: string = 'text-blue-500';
</script>

<div class="bg-white rounded-lg p-4">
  <h3>{title}</h3>
  <p class="text-2xl {color}">{value}</p>
  {#if subtitle}
    <p class="text-sm">{subtitle}</p>
  {/if}
</div>
```

**Pattern**: Always provide defaults for optional props.

### Lifecycle Hooks

**Use `onMount` for initialization and cleanup**:

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { metrics } from '@stores/metrics';

  let intervalId: number | null = null;

  onMount(() => {
    // Refresh immediately
    metrics.refresh();

    // Set up interval
    intervalId = setInterval(() => {
      metrics.refresh();
    }, 1000);

    // Cleanup function
    return () => {
      if (intervalId) {
        clearInterval(intervalId);
      }
    };
  });
</script>
```

---

## Platform-Specific Code

### Conditional Compilation

**Use `cfg!` macro for platform-specific code**:

```rust
pub fn get_config_path() -> Result<PathBuf> {
    let config_dir = if cfg!(target_os = "linux") {
        dirs::config_dir()
            .ok_or_else(|| AppError::ConfigError("Cannot find config directory".to_string()))?
            .join("divoom-pc-companion")
    } else if cfg!(target_os = "macos") {
        dirs::config_dir()
            .ok_or_else(|| AppError::ConfigError("Cannot find config directory".to_string()))?
            .join("DivoomPCCompanion")
    } else if cfg!(target_os = "windows") {
        dirs::config_dir()
            .ok_or_else(|| AppError::ConfigError("Cannot find config directory".to_string()))?
            .join("DivoomPCCompanion")
    } else {
        return Err(AppError::ConfigError("Unsupported platform".to_string()));
    };

    fs::create_dir_all(&config_dir)?;
    Ok(config_dir.join("config.json"))
}
```

### Platform-Specific Dependencies

**Declare platform-specific crates in `Cargo.toml`**:

```toml
[target.'cfg(target_os = "windows")'.dependencies]
wmi = "0.14"

[target.'cfg(target_os = "linux")'.dependencies]
zbus = "5"

[target.'cfg(target_os = "macos")'.dependencies]
core-foundation = "0.10"
```

---

## Configuration Management

### Default Config Pattern

**Provide sensible defaults**:

```rust
pub fn default_config() -> AppConfig {
    AppConfig {
        general: GeneralConfig {
            language: "auto".to_string(),
            start_minimized: false,
            autostart: false,
            update_interval: 1000,
        },
        metrics: MetricsConfig {
            enable_cpu: true,
            enable_gpu: true,
            enable_memory: true,
            enable_storage: true,
            storage_path: default_storage_path(),
        },
        // ...
    }
}
```

**Load with fallback**:
```rust
let config = load_config().unwrap_or_else(|_| default_config());
```

---

## Device Protocol

### HTTP Communication

**Pattern**: Use `reqwest` for HTTP, with timeout:

```rust
pub async fn discover_devices() -> Result<Vec<DivoomDevice>> {
    let client = Client::new();
    let response = client
        .get(DISCOVERY_URL)
        .timeout(std::time::Duration::from_secs(5))
        .send()
        .await
        .map_err(|_e| AppError::DeviceDiscoveryFailed)?;

    if !response.status().is_success() {
        return Err(AppError::DeviceDiscoveryFailed);
    }

    let json: Value = response.json().await?;
    // Parse response...
}
```

### Serde Field Renaming

**Match API field names with `serde(rename)`**:

```rust
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DeviceListResponse {
    #[serde(rename = "DeviceName")]
    pub device_name: String,

    #[serde(rename = "ProductName")]
    pub product_name: String,

    #[serde(rename = "Hardware")]
    pub hardware: String,
}
```

---

## UI Conventions

### Tailwind CSS Integration

**Custom colors in `tailwind.config.js`**:

```javascript
export default {
  theme: {
    extend: {
      colors: {
        'divoom-primary': '#FF6B35',
        'divoom-secondary': '#004E89',
        'divoom-dark': '#1A1A2E',
        'divoom-light': '#F7F7F7'
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}
```

**Dark mode support**:
```svelte
<div class="bg-white dark:bg-gray-900">
  <p class="text-gray-900 dark:text-gray-100">
    Auto dark mode based on system preference
  </p>
</div>
```

---

## Internationalization

### i18n Directory Structure

```
src/lib/i18n/locales/
├── en.json       # English (default)
├── ru.json       # Russian
├── zh-CN.json    # Chinese Simplified
└── ...
```

**Locale file format**:
```json
{
  "app": {
    "title": "Divoom PC Companion",
    "version": "Version {{version}}"
  },
  "metrics": {
    "cpu_usage": "CPU Usage",
    "cpu_temp": "CPU Temperature"
  }
}
```

---

## Testing Strategy

### Current State: No Tests

**Observation**: The project currently has 0 test files.

**Recommendation**: Add tests for:
1. **Unit tests**: Metric collection functions, config serialization
2. **Integration tests**: Device discovery, protocol parsing
3. **E2E tests**: Critical user flows with Playwright

---

## Build Configuration

### Tauri Config

**Key settings in `tauri.conf.json`**:

```json
{
  "app": {
    "windows": [{
      "title": "Divoom PC Companion",
      "width": 350,
      "height": 250,
      "decorations": false,
      "transparent": false,
      "alwaysOnTop": true,
      "resizable": false
    }],
    "security": {
      "csp": null  // TODO: Add CSP
    }
  },
  "bundle": {
    "active": true,
    "targets": "all",
    "icon": ["icons/32x32.png", "icons/128x128.png", "..."]
  }
}
```

**Pattern**: Frameless window for widget-style UI.

---

## Development Workflow

### Running the App

```bash
# Development mode
npm run tauri:dev

# Build for production
npm run tauri:build

# Type checking
npm run check
```

**Key**: Vite dev server runs on `http://localhost:1420` (not default 5173).

---

## Common Patterns

### Mutex for Shared State

**Use `Mutex<T>` for shared mutable state**:

```rust
pub struct AppState {
    pub collector: Mutex<MetricsCollector>,
    pub config: Mutex<AppConfig>,
}

// Access in commands
let config = state.config.lock().unwrap();
```

**TODO**: Replace `.unwrap()` with proper error handling.

### Default Trait for Data Structures

**Implement `Default` for complex types**:

```rust
impl Default for CpuMetrics {
    fn default() -> Self {
        Self {
            usage_percent: 0.0,
            temperature: None,
        }
    }
}
```

**Benefits**: Easy to create zero-values and fallback instances.

---

## Security Considerations

### Current Issues

1. **Missing CSP**: Content Security Policy is `null` in tauri.conf.json
2. **HTTP (not HTTPS)**: Device communication uses unencrypted HTTP
3. **No input validation**: Tauri commands don't validate IP addresses or IDs
4. **Unsafe mutex unwrap**: Multiple `.lock().unwrap()` calls

**Recommendations**:
```json
"security": {
  "csp": "default-src 'self'; connect-src 'self' http://app.divoom-gz.com http://localhost:*; script-src 'self';"
}
```

```rust
// Validate IP address
if ip.parse::<std::net::IpAddr>().is_err() {
    return Err("Invalid IP address format".to_string());
}
```

---

## Performance Patterns

### Refresh Intervals

**1-second refresh for metrics**:
```typescript
onMount(() => {
  metrics.refresh();
  intervalId = setInterval(() => metrics.refresh(), 1000);
});
```

**Make intervals configurable** via settings.

### Resource Limits

**Target**: < 50MB RAM, < 1% CPU idle

**Strategies**:
- Lazy load GPU monitoring
- Cache device list
- Reuse HTTP connections

---

## Documentation

### README Structure

The project includes:
1. **README.md**: User-facing documentation
2. **SPECIFICATION.md**: Complete technical specification
3. **READY_TO_RUN.md**: Build status and quick start

**Pattern**: Separate user docs from developer specs.

---

## Git History Patterns

### Commit Messages

**Most recent commits follow descriptive format**:
```
Improve CPU temperature monitoring and documentation - Add multiple methods...
Full code rework and Linux version
```

**Pattern**: Descriptive titles, detailed body explaining changes.

**No conventional commits**: Unlike the skill description, this repo doesn't use `feat:`, `fix:`, etc.

---

## What's Different Here

### Unlike Typical Web Projects

1. **Tauri, not webpack**: Desktop app framework
2. **Rust backend, not Node.js**: Native system access
3. **System metrics**: Cross-platform hardware monitoring
4. **Device protocol**: HTTP to local devices (not APIs)

### Unique to This Project

1. **Multi-language i18n**: 19 language support
2. **Platform-specific code**: Conditional compilation for Win/Linux/mac
3. **Hardware monitoring**: CPU/GPU/memory/storage metrics
4. **Widget UI**: Frameless, always-on-top window

---

## Key Takeaways

1. **Modular architecture** with clear domain boundaries
2. **Error handling** via `thiserror` for structured types
3. **Async I/O** with `tokio` for network operations
4. **Svelte stores** for state management
5. **TypeScript** for type safety across backend boundary
6. **Platform-specific code** via `cfg!` macro
7. **Serde** for all data/config serialization
8. **Tauri commands** for frontend-backend communication

---

**Repository**: DivoomPCMonitorTool
**Primary Language**: Rust (backend), TypeScript (frontend)
**Last Updated**: 2025-02-06
**Branch**: `rewrite`
