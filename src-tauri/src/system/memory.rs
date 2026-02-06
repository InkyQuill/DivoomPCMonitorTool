use crate::core::{MemoryMetrics, Result};
use sysinfo::System;

pub fn get_memory_metrics(system: &System) -> Result<MemoryMetrics> {
    let total_memory = system.total_memory();
    let available_memory = system.available_memory();
    let used_memory = total_memory.saturating_sub(available_memory);

    let total_gb = total_memory as f64 / 1024.0 / 1024.0 / 1024.0;
    let used_gb = used_memory as f64 / 1024.0 / 1024.0 / 1024.0;
    let available_gb = available_memory as f64 / 1024.0 / 1024.0 / 1024.0;

    let usage_percent = if total_memory > 0 {
        (used_memory as f32 / total_memory as f32) * 100.0
    } else {
        0.0
    };

    Ok(MemoryMetrics {
        total_gb,
        used_gb,
        available_gb,
        usage_percent,
    })
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_get_memory_metrics_returns_valid_values() {
        let mut system = System::new_all();
        system.refresh_memory();

        let metrics = get_memory_metrics(&system);

        assert!(metrics.is_ok());
        let mem = metrics.unwrap();
        assert!(mem.total_gb > 0.0);
        assert!(mem.used_gb >= 0.0);
        assert!(mem.available_gb >= 0.0);
        assert!(mem.usage_percent >= 0.0);
        assert!(mem.usage_percent <= 100.0);
    }

    #[test]
    fn test_memory_values_consistency() {
        let mut system = System::new_all();
        system.refresh_memory();

        let metrics = get_memory_metrics(&system).unwrap();

        // Total should roughly equal used + available
        let sum = metrics.used_gb + metrics.available_gb;
        let diff = (metrics.total_gb - sum).abs();

        // Allow small rounding difference
        assert!(
            diff < 0.1,
            "Total GB ({}) should equal used + available ({})",
            metrics.total_gb,
            sum
        );
    }

    #[test]
    fn test_usage_percent_calculation() {
        let mut system = System::new_all();
        system.refresh_memory();

        let metrics = get_memory_metrics(&system).unwrap();

        // Verify percentage calculation matches actual ratio
        let expected_percent = (metrics.used_gb / metrics.total_gb) * 100.0;
        let diff = (metrics.usage_percent as f64 - expected_percent).abs();

        assert!(
            diff < 0.5,
            "Usage percent ({}) should match ratio ({})",
            metrics.usage_percent,
            expected_percent
        );
    }

    #[test]
    fn test_values_in_gigabytes() {
        let mut system = System::new_all();
        system.refresh_memory();

        let metrics = get_memory_metrics(&system).unwrap();

        // All values should be in GB (not KB/MB)
        // Total memory should be > 0 for any modern system
        assert!(metrics.total_gb > 0.0);
        // Most systems have at least 1GB
        assert!(metrics.total_gb >= 1.0);
    }
}
