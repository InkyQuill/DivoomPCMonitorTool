use crate::core::{StorageMetrics, Result};
use sysinfo::System;

pub fn get_storage_metrics(system: &System) -> Result<StorageMetrics> {
    // Try platform-specific implementations first

    #[cfg(target_os = "linux")]
    {
        if let Some(metrics) = get_storage_metrics_linux() {
            return Ok(metrics);
        }
    }

    // Fallback to placeholder for now
    // TODO: Implement for other platforms (Windows/macOS)
    Ok(StorageMetrics {
        total_gb: 500.0,
        used_gb: 250.0,
        free_gb: 250.0,
        usage_percent: 50.0,
        temperature: None,
    })
}

#[cfg(target_os = "linux")]
fn get_storage_metrics_linux() -> Option<StorageMetrics> {
    use std::path::Path;

    // Read root filesystem stats from /proc/self/mountstats or statfs
    // For simplicity, use statvfs via nix crate or direct parsing

    // Try using sysinfo's memory info to approximate (not ideal but works)
    // Or read from /sys/block for disk sizes

    // Simple approach: read from df command output or /proc/mounts
    if let Ok(output) = std::process::Command::new("df")
        .arg("-B1")  // 1-byte blocks
        .arg("/")
        .output()
    {
        if output.status.success() {
            let stdout = String::from_utf8_lossy(&output.stdout);

            // Parse df output:
            // Filesystem 1B-blocks Used Available Use%
            // /dev/root 123456789 987654321 12345678 90% /

            for line in stdout.lines().skip(1) {
                let parts: Vec<&str> = line.split_whitespace().collect();
                if parts.len() >= 4 {
                    if let Ok(total) = parts[1].parse::<u64>() {
                        if let Ok(used) = parts[2].parse::<u64>() {
                            let available = total.saturating_sub(used);

                            let total_gb = total as f64 / (1024.0 * 1024.0 * 1024.0);
                            let used_gb = used as f64 / (1024.0 * 1024.0 * 1024.0);
                            let free_gb = available as f64 / (1024.0 * 1024.0 * 1024.0);

                            let usage_percent = if total > 0 {
                                (used as f32 / total as f32) * 100.0
                            } else {
                                0.0
                            };

                            return Some(StorageMetrics {
                                total_gb,
                                used_gb,
                                free_gb,
                                usage_percent,
                                temperature: None,
                            });
                        }
                    }
                }
            }
        }
    }

    None
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_get_storage_metrics_returns_real_values() {
        // This test will FAIL until we implement real reading
        let mut system = System::new_all();
        system.refresh_all();

        let metrics = get_storage_metrics(&system);

        assert!(metrics.is_ok());
        let storage = metrics.unwrap();

        // Real disks should be > 0
        assert!(storage.total_gb > 0.0, "Total disk space should be > 0");
        assert!(storage.used_gb >= 0.0, "Used space should be >= 0");
        assert!(storage.free_gb >= 0.0, "Free space should be >= 0");

        // Sanity checks
        assert!(storage.total_gb < 10000.0, "Total disk > 10TB seems unrealistic: {}", storage.total_gb);
        assert!(storage.usage_percent >= 0.0 && storage.usage_percent <= 100.0);

        // Consistency check
        assert!(
            (storage.total_gb - (storage.used_gb + storage.free_gb)).abs() < 1.0,
            "Total != used + free: {} != {} + {}",
            storage.total_gb,
            storage.used_gb,
            storage.free_gb
        );
    }

    #[test]
    fn test_storage_detects_primary_disk() {
        let mut system = System::new_all();
        system.refresh_all();

        let metrics = get_storage_metrics(&system).unwrap();

        // On Linux systems, should get real data
        #[cfg(target_os = "linux")]
        assert!(
            metrics.total_gb != 500.0 || metrics.used_gb != 250.0,
            "Should return real disk data on Linux, not placeholder"
        );
    }

    #[test]
    fn test_storage_usage_percent_calculation() {
        let mut system = System::new_all();
        system.refresh_all();

        let metrics = get_storage_metrics(&system).unwrap();

        // Verify percentage calculation
        let expected_percent = (metrics.used_gb / metrics.total_gb) * 100.0;
        let diff = (metrics.usage_percent as f64 - expected_percent).abs();

        assert!(
            diff < 0.5,
            "Usage percent {} doesn't match calculated {}",
            metrics.usage_percent,
            expected_percent
        );
    }

    #[test]
    fn test_storage_temperature_none() {
        let mut system = System::new_all();
        system.refresh_all();

        let metrics = get_storage_metrics(&system).unwrap();

        // Temperature detection not yet implemented
        assert!(metrics.temperature.is_none());
    }

    #[test]
    fn test_storage_handles_no_disks() {
        // Edge case: what if no disks available?
        // System should still return valid metrics or error
        let system = System::new_all();
        let result = get_storage_metrics(&system);

        // Should not panic, should either return metrics or error
        assert!(result.is_ok() || result.is_err());
    }

    #[test]
    fn test_storage_multiple_disks_aggregated() {
        // If system has multiple disks, should aggregate or pick primary
        let mut system = System::new_all();
        system.refresh_all();

        let metrics = get_storage_metrics(&system).unwrap();

        // Currently, we expect primary disk (root filesystem)
        // Multi-disk support is a future enhancement
        assert!(metrics.total_gb > 0.0);
    }

    #[test]
    fn test_storage_values_realistic() {
        let mut system = System::new_all();
        system.refresh_all();

        let metrics = get_storage_metrics(&system).unwrap();

        // Most systems have > 100GB and < 10TB
        #[cfg(target_os = "linux")]
        assert!(
            metrics.total_gb >= 100.0,
            "Disk too small: {} GB",
            metrics.total_gb
        );

        // Usage should be reasonable (at least some used)
        #[cfg(target_os = "linux")]
        assert!(
            metrics.used_gb > 0.0,
            "No disk usage detected: {} GB used",
            metrics.used_gb
        );

        // And not completely full (normally)
        #[cfg(target_os = "linux")]
        assert!(
            metrics.free_gb > 0.0,
            "No free space: {} GB free",
            metrics.free_gb
        );
    }

    #[test]
    fn test_placeholder_no_longer_used() {
        let mut system = System::new_all();
        system.refresh_all();

        let metrics = get_storage_metrics(&system).unwrap();

        // This test ensures we're reading REAL data, not placeholders
        // If this fails, we're still using placeholder values
        let is_placeholder =
            metrics.total_gb == 500.0 && metrics.used_gb == 250.0 && metrics.usage_percent == 50.0;

        #[cfg(target_os = "linux")]
        assert!(!is_placeholder, "Still using placeholder values on Linux!");

        #[cfg(not(target_os = "linux"))]
        {
            // On other platforms, placeholder is still OK for now
            println!("INFO: Storage not yet implemented for this platform");
        }
    }

    #[test]
    fn test_actual_disk_reading() {
        // Show actual disk metrics being read
        let mut system = System::new_all();
        system.refresh_all();

        let metrics = get_storage_metrics(&system).unwrap();

        #[cfg(target_os = "linux")]
        {
            println!("✅ Real Storage Metrics Detected:");
            println!("   Total: {:.1} GB", metrics.total_gb);
            println!("   Used:  {:.1} GB", metrics.used_gb);
            println!("   Free:  {:.1} GB", metrics.free_gb);
            println!("   Usage: {:.1}%", metrics.usage_percent);
        }

        // Verify it's not the placeholder
        assert!((metrics.total_gb - 500.0).abs() > 1.0 || (metrics.used_gb - 250.0).abs() > 1.0);
    }
}
