use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AppConfig {
    pub general: GeneralConfig,
    pub metrics: MetricsConfig,
    pub display: DisplayConfig,
    pub device: DeviceConfig,
    pub advanced: AdvancedConfig,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GeneralConfig {
    pub language: String,
    pub start_minimized: bool,
    pub autostart: bool,
    pub update_interval: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MetricsConfig {
    pub enable_cpu: bool,
    pub enable_gpu: bool,
    pub enable_memory: bool,
    pub enable_storage: bool,
    pub storage_path: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DisplayConfig {
    pub screen_index: u32,
    pub brightness: u32,
    pub theme_mode: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DeviceConfig {
    pub ip_address: String,
    pub device_type: String,
    pub connection_timeout: u64,
    pub retry_attempts: u32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AdvancedConfig {
    pub log_level: String,
    pub log_file: String,
}
