// AMD GPU monitoring stub
// TODO: Implement AMD GPU monitoring

pub fn get_amd_metrics() -> crate::core::Result<crate::core::GpuMetrics> {
    Ok(crate::core::GpuMetrics {
        usage_percent: 0.0,
        temperature: None,
        vendor: "AMD".to_string(),
    })
}
