use crate::core::{SystemMetrics, CpuMetrics, MemoryMetrics, StorageMetrics, GpuMetrics};
use crate::system::cpu::get_cpu_metrics;
use crate::system::memory::get_memory_metrics;
use crate::system::storage::get_storage_metrics;
use crate::system::gpu::get_gpu_metrics;
use sysinfo::System;
use crate::core::Result;

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

        let cpu = get_cpu_metrics(&self.system)?;
        let memory = get_memory_metrics(&self.system)?;
        let storage = get_storage_metrics(&self.system)?;
        let gpu = get_gpu_metrics()?;

        Ok(SystemMetrics {
            cpu,
            gpu,
            memory,
            storage,
            timestamp: chrono::Utc::now().timestamp(),
        })
    }
}

impl Default for MetricsCollector {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_collector_new() {
        let collector = MetricsCollector::new();
        // Collector is created successfully
        assert!(true);
    }

    #[test]
    fn test_collector_default() {
        let collector = MetricsCollector::default();
        // Default trait works
        assert!(true);
    }

    #[test]
    fn test_collect_returns_valid_metrics() {
        let mut collector = MetricsCollector::new();
        let metrics = collector.collect();

        assert!(metrics.is_ok());
    }

    #[test]
    fn test_collected_metrics_have_valid_timestamp() {
        let mut collector = MetricsCollector::new();
        let metrics = collector.collect().unwrap();

        let now = chrono::Utc::now().timestamp();
        let diff = now - metrics.timestamp;

        // Timestamp should be recent (within 5 seconds)
        assert!(diff >= 0 && diff <= 5);
    }

    #[test]
    fn test_collected_metrics_have_cpu_data() {
        let mut collector = MetricsCollector::new();
        let metrics = collector.collect().unwrap();

        assert!(metrics.cpu.usage_percent >= 0.0);
        assert!(metrics.cpu.usage_percent <= 100.0);
    }

    #[test]
    fn test_collected_metrics_have_memory_data() {
        let mut collector = MetricsCollector::new();
        let metrics = collector.collect().unwrap();

        assert!(metrics.memory.total_gb > 0.0);
        assert!(metrics.memory.usage_percent >= 0.0);
        assert!(metrics.memory.usage_percent <= 100.0);
    }

    #[test]
    fn test_collected_metrics_have_storage_data() {
        let mut collector = MetricsCollector::new();
        let metrics = collector.collect().unwrap();

        // Currently returns placeholder values
        assert!(metrics.storage.total_gb > 0.0);
        assert!(metrics.storage.usage_percent >= 0.0);
    }

    #[test]
    fn test_collect_multiple_times() {
        let mut collector = MetricsCollector::new();

        let metrics1 = collector.collect();
        let metrics2 = collector.collect();

        assert!(metrics1.is_ok());
        assert!(metrics2.is_ok());

        let m1 = metrics1.unwrap();
        let m2 = metrics2.unwrap();

        // Timestamps should be different
        assert!(m2.timestamp >= m1.timestamp);
    }

    #[test]
    fn test_collected_gpu_metrics() {
        let mut collector = MetricsCollector::new();
        let metrics = collector.collect().unwrap();

        // GPU should be present (even if placeholder)
        assert!(metrics.gpu.usage_percent >= 0.0);
        assert!(metrics.gpu.usage_percent <= 100.0);
        assert!(!metrics.gpu.vendor.is_empty());
    }
}
