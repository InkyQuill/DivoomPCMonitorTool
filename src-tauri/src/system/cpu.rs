use crate::core::{CpuMetrics, Result};
use sysinfo::System;

#[cfg(target_os = "linux")]
mod linux;

pub fn get_cpu_metrics(system: &System) -> Result<CpuMetrics> {
    let usage_percent = system.global_cpu_usage();

    #[cfg(target_os = "linux")]
    let temperature = linux::get_cpu_temperature();

    #[cfg(not(target_os = "linux"))]
    let temperature = None; // TODO: Implement for other platforms

    Ok(CpuMetrics {
        usage_percent,
        temperature,
    })
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_get_cpu_metrics_returns_valid_metrics() {
        let mut system = System::new_all();
        system.refresh_all();

        let metrics = get_cpu_metrics(&system);

        assert!(metrics.is_ok());
        let cpu = metrics.unwrap();
        assert!(cpu.usage_percent >= 0.0);
        assert!(cpu.usage_percent <= 100.0);

        // Temperature is now implemented on Linux
        #[cfg(target_os = "linux")]
        {
            if cpu.temperature.is_some() {
                let temp = cpu.temperature.unwrap();
                assert!(temp >= 20.0 && temp <= 110.0);
            }
            // May be None if no thermal sensor available
        }

        #[cfg(not(target_os = "linux"))]
        assert!(cpu.temperature.is_none()); // Not yet implemented on other platforms
    }

    #[test]
    fn test_cpu_usage_range() {
        let mut system = System::new_all();
        system.refresh_cpu_all();

        let metrics = get_cpu_metrics(&system).unwrap();
        assert!(metrics.usage_percent >= 0.0);
        assert!(metrics.usage_percent <= 100.0);
    }

    #[test]
    fn test_temperature_none_initially() {
        let system = System::new_all();
        let metrics = get_cpu_metrics(&system).unwrap();

        // On non-Linux platforms, temperature is None
        // On Linux, it may be Some if thermal sensor is available
        #[cfg(not(target_os = "linux"))]
        assert!(metrics.temperature.is_none());

        #[cfg(target_os = "linux")]
        {
            // On Linux, temperature may be available
            if metrics.temperature.is_some() {
                let temp = metrics.temperature.unwrap();
                assert!(temp >= 0.0 && temp <= 150.0);
            }
        }
    }
}
