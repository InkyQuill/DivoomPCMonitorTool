use crate::core::{GpuMetrics, Result};

mod nvidia;
mod amd;
mod intel;
mod generic;

pub fn get_gpu_metrics() -> Result<GpuMetrics> {
    // Try to detect GPU vendor and get metrics
    // For now, return default metrics
    // TODO: Implement vendor-specific detection

    Ok(GpuMetrics {
        usage_percent: 0.0,
        temperature: None,
        vendor: "Unknown".to_string(),
    })
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_get_gpu_metrics_returns_ok() {
        let result = get_gpu_metrics();
        assert!(result.is_ok());
    }

    #[test]
    fn test_get_gpu_metrics_returns_valid_metrics() {
        let metrics = get_gpu_metrics().unwrap();

        assert!(metrics.usage_percent >= 0.0);
        assert!(metrics.usage_percent <= 100.0);
        assert!(!metrics.vendor.is_empty());
    }

    #[test]
    fn test_get_gpu_metrics_temperature_none() {
        let metrics = get_gpu_metrics().unwrap();

        // Temperature detection not yet implemented
        assert!(metrics.temperature.is_none());
    }

    #[test]
    fn test_get_gpu_metrics_vendor_unknown() {
        let metrics = get_gpu_metrics().unwrap();

        // Currently returns "Unknown" until vendor detection is implemented
        assert_eq!(metrics.vendor, "Unknown");
    }

    #[test]
    fn test_get_gpu_metrics_usage_zero() {
        let metrics = get_gpu_metrics().unwrap();

        // Currently returns 0.0 until GPU monitoring is implemented
        assert_eq!(metrics.usage_percent, 0.0);
    }

    #[test]
    fn test_get_gpu_metrics_multiple_calls() {
        let metrics1 = get_gpu_metrics().unwrap();
        let metrics2 = get_gpu_metrics().unwrap();

        // Should return consistent results
        assert_eq!(metrics1.usage_percent, metrics2.usage_percent);
        assert_eq!(metrics1.vendor, metrics2.vendor);
    }
}
