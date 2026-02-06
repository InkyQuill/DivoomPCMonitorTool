use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CpuMetrics {
    pub usage_percent: f32,
    pub temperature: Option<f32>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GpuMetrics {
    pub usage_percent: f32,
    pub temperature: Option<f32>,
    pub vendor: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MemoryMetrics {
    pub total_gb: f64,
    pub used_gb: f64,
    pub available_gb: f64,
    pub usage_percent: f32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct StorageMetrics {
    pub total_gb: f64,
    pub used_gb: f64,
    pub free_gb: f64,
    pub usage_percent: f32,
    pub temperature: Option<f32>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SystemMetrics {
    pub cpu: CpuMetrics,
    pub gpu: GpuMetrics,
    pub memory: MemoryMetrics,
    pub storage: StorageMetrics,
    pub timestamp: i64,
}

impl Default for CpuMetrics {
    fn default() -> Self {
        Self {
            usage_percent: 0.0,
            temperature: None,
        }
    }
}

impl Default for GpuMetrics {
    fn default() -> Self {
        Self {
            usage_percent: 0.0,
            temperature: None,
            vendor: "Unknown".to_string(),
        }
    }
}

impl Default for MemoryMetrics {
    fn default() -> Self {
        Self {
            total_gb: 0.0,
            used_gb: 0.0,
            available_gb: 0.0,
            usage_percent: 0.0,
        }
    }
}

impl Default for StorageMetrics {
    fn default() -> Self {
        Self {
            total_gb: 0.0,
            used_gb: 0.0,
            free_gb: 0.0,
            usage_percent: 0.0,
            temperature: None,
        }
    }
}

impl Default for SystemMetrics {
    fn default() -> Self {
        Self {
            cpu: CpuMetrics::default(),
            gpu: GpuMetrics::default(),
            memory: MemoryMetrics::default(),
            storage: StorageMetrics::default(),
            timestamp: chrono::Utc::now().timestamp(),
        }
    }
}
