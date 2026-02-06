// TypeScript types matching Rust structs

export interface CpuMetrics {
  usage_percent: number;
  temperature: number | null;
}

export interface GpuMetrics {
  usage_percent: number;
  temperature: number | null;
  vendor: string;
}

export interface MemoryMetrics {
  total_gb: number;
  used_gb: number;
  available_gb: number;
  usage_percent: number;
}

export interface StorageMetrics {
  total_gb: number;
  used_gb: number;
  free_gb: number;
  usage_percent: number;
  temperature: number | null;
}

export interface SystemMetrics {
  cpu: CpuMetrics;
  gpu: GpuMetrics;
  memory: MemoryMetrics;
  storage: StorageMetrics;
  timestamp: number;
}

export interface DivoomDevice {
  ip: string;
  product_name: string;
  device_name: string;
  hardware: string;
  mac: string;
}

export interface AppConfig {
  general: GeneralConfig;
  metrics: MetricsConfig;
  display: DisplayConfig;
  device: DeviceConfig;
  advanced: AdvancedConfig;
}

export interface GeneralConfig {
  language: string;
  start_minimized: boolean;
  autostart: boolean;
  update_interval: number;
}

export interface MetricsConfig {
  enable_cpu: boolean;
  enable_gpu: boolean;
  enable_memory: boolean;
  enable_storage: boolean;
  storage_path: string;
}

export interface DisplayConfig {
  screen_index: number;
  brightness: number;
  theme_mode: string;
}

export interface DeviceConfig {
  ip_address: string;
  device_type: string;
  connection_timeout: number;
  retry_attempts: number;
}

export interface AdvancedConfig {
  log_level: string;
  log_file: string;
}
