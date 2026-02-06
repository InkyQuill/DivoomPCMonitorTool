// Generic GPU monitoring stub

pub fn get_generic_metrics() -> crate::core::Result<crate::core::GpuMetrics> {
    Ok(crate::core::GpuMetrics {
        usage_percent: 0.0,
        temperature: None,
        vendor: "Unknown".to_string(),
    })
}
