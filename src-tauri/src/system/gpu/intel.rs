// Intel GPU monitoring stub
// TODO: Implement Intel GPU monitoring

pub fn get_intel_metrics() -> crate::core::Result<crate::core::GpuMetrics> {
    Ok(crate::core::GpuMetrics {
        usage_percent: 0.0,
        temperature: None,
        vendor: "Intel".to_string(),
    })
}
