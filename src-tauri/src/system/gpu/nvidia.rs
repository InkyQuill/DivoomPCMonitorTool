// NVIDIA GPU monitoring stub
// TODO: Implement NVIDIA GPU monitoring using NVML

pub fn get_nvidia_metrics() -> crate::core::Result<crate::core::GpuMetrics> {
    Ok(crate::core::GpuMetrics {
        usage_percent: 0.0,
        temperature: None,
        vendor: "NVIDIA".to_string(),
    })
}
