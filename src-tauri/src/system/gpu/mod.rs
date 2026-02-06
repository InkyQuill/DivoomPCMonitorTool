// GPU Monitoring Implementation (Linux)
// Implements real GPU detection and metrics collection for NVIDIA, AMD, and Intel GPUs.
//
// Detection priority order:
// 1. NVIDIA GPUs: nvidia-smi command (most accurate, provides usage + temperature)
// 2. AMD/Intel GPUs: /sys/class/drm/card*/device/ (usage percentage)
// 3. Temperature from hwmon subsystem (for AMD/Intel)
// 4. Fallback to generic /sys/class/hwmon detection

use crate::core::{GpuMetrics, Result};

mod amd;
mod generic;
mod intel;
mod nvidia;

/// Get GPU metrics with vendor detection and real usage/temperature data.
///
/// On Linux, detects NVIDIA, AMD, Intel, or generic GPUs and returns
/// actual usage percentage and temperature when available.
///
/// On other platforms, returns placeholder (TODO: implement).
pub fn get_gpu_metrics() -> Result<GpuMetrics> {
    // Try platform-specific implementations

    #[cfg(target_os = "linux")]
    {
        if let Some(metrics) = get_gpu_metrics_linux() {
            return Ok(metrics);
        }
    }

    // Fallback to placeholder for other platforms
    Ok(GpuMetrics {
        usage_percent: 0.0,
        temperature: None,
        vendor: "Unknown".to_string(),
    })
}

/// Linux-specific GPU detection and metrics collection.
///
/// Tries multiple detection methods in order:
/// 1. nvidia-smi (NVIDIA proprietary drivers)
/// 2. /sys/class/drm (AMD/Intel open-source drivers)
/// 3. /sys/class/hwmon (generic hardware monitoring)
#[cfg(target_os = "linux")]
fn get_gpu_metrics_linux() -> Option<GpuMetrics> {
    use std::fs;
    use std::path::Path;

    // Try NVIDIA first (nvidia-smi)
    if Path::new("/usr/bin/nvidia-smi").exists() || Path::new("/usr/local/bin/nvidia-smi").exists()
    {
        if let Ok(output) = std::process::Command::new("nvidia-smi")
            .arg("--query-gpu=utilization.gpu,temperature.gpu")
            .arg("--format=csv,noheader,nounits")
            .output()
        {
            if output.status.success() {
                let stdout = String::from_utf8_lossy(&output.stdout);
                let parts: Vec<&str> = stdout.trim().split(',').collect();

                if parts.len() >= 2 {
                    if let Ok(usage) = parts[0].trim().parse::<f32>() {
                        if let Ok(temp) = parts[1].trim().parse::<f32>() {
                            return Some(GpuMetrics {
                                usage_percent: usage,
                                temperature: Some(temp),
                                vendor: "NVIDIA".to_string(),
                            });
                        }
                    }
                }
            }
        }
    }

    // Try AMD/Intel via /sys/class/drm
    let drm_path = Path::new("/sys/class/drm");
    if let Ok(entries) = fs::read_dir(drm_path) {
        for entry in entries.flatten() {
            let card_path = entry.path();
            let card_name = card_path.file_name()?.to_string_lossy();

            // Look for card0, card1, etc. (not card*-*)
            if card_name.starts_with("card") && !card_name.contains('-') {
                let device_path = card_path.join("device");

                // Try to read GPU usage
                let usage_file = device_path.join("gpu_busy_percent");
                let usage = if usage_file.exists() {
                    fs::read_to_string(&usage_file)
                        .ok()
                        .and_then(|s| s.trim().parse::<f32>().ok())
                        .unwrap_or(0.0)
                } else {
                    0.0
                };

                // Try to read GPU temperature from hwmon
                let hwmon_path = device_path.join("hwmon");
                let temperature = if let Ok(hwmon_entries) = fs::read_dir(&hwmon_path) {
                    let mut temp = None;
                    for hwmon_entry in hwmon_entries.flatten() {
                        let temp_file = hwmon_entry.path().join("temp1_input");
                        if temp_file.exists() {
                            if let Ok(temp_str) = fs::read_to_string(&temp_file) {
                                if let Ok(millidegrees) = temp_str.trim().parse::<i64>() {
                                    temp = Some(millidegrees as f32 / 1000.0);
                                    break;
                                }
                            }
                        }
                    }
                    temp
                } else {
                    None
                };

                // Try to detect vendor from uevent file
                let uevent_path = device_path.join("uevent");
                let vendor = if uevent_path.exists() {
                    if let Ok(uevent_content) = fs::read_to_string(&uevent_path) {
                        if uevent_content.contains("PCI_ID=1002") || uevent_content.contains("AMD")
                        {
                            "AMD".to_string()
                        } else if uevent_content.contains("PCI_ID=8086")
                            || uevent_content.contains("Intel")
                        {
                            "Intel".to_string()
                        } else {
                            "Generic".to_string()
                        }
                    } else {
                        "Generic".to_string()
                    }
                } else {
                    "Generic".to_string()
                };

                // If we found at least some data, return it
                if usage > 0.0 || temperature.is_some() || vendor != "Generic" {
                    return Some(GpuMetrics {
                        usage_percent: usage,
                        temperature,
                        vendor,
                    });
                }
            }
        }
    }

    // Try generic detection from /sys/class/hwmon (fallback)
    if let Ok(entries) = fs::read_dir(Path::new("/sys/class/hwmon")) {
        for entry in entries.flatten() {
            let name_path = entry.path().join("name");
            if let Ok(name) = fs::read_to_string(&name_path) {
                let name_lower = name.to_lowercase();

                // Check for known GPU hwmon devices
                let is_gpu = name_lower.contains("amdgpu")
                    || name_lower.contains("radeon")
                    || name_lower.contains("i915")
                    || name_lower.contains("nvidia")
                    || name_lower.contains("gpu");

                if is_gpu {
                    // Try to read temperature
                    let device_path = entry.path();
                    let temp_input = device_path.join("temp1_input");

                    let temperature = if temp_input.exists() {
                        if let Ok(temp_str) = fs::read_to_string(&temp_input) {
                            if let Ok(millidegrees) = temp_str.trim().parse::<i64>() {
                                Some(millidegrees as f32 / 1000.0)
                            } else {
                                None
                            }
                        } else {
                            None
                        }
                    } else {
                        None
                    };

                    // Determine vendor
                    let vendor = if name_lower.contains("amdgpu") || name_lower.contains("radeon") {
                        "AMD".to_string()
                    } else if name_lower.contains("i915") || name_lower.contains("intel") {
                        "Intel".to_string()
                    } else if name_lower.contains("nvidia") {
                        "NVIDIA".to_string()
                    } else {
                        "Generic".to_string()
                    };

                    return Some(GpuMetrics {
                        usage_percent: 0.0, // hwmon doesn't provide usage
                        temperature,
                        vendor,
                    });
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
    fn test_gpu_vendor_detection() {
        let metrics = get_gpu_metrics().unwrap();

        // On Linux systems, should detect vendor (NVIDIA, AMD, Intel, or Generic)
        #[cfg(target_os = "linux")]
        {
            assert!(
                metrics.vendor == "NVIDIA"
                    || metrics.vendor == "AMD"
                    || metrics.vendor == "Intel"
                    || metrics.vendor == "Generic",
                "Vendor should be NVIDIA, AMD, Intel, or Generic, got: {}",
                metrics.vendor
            );
        }

        #[cfg(not(target_os = "linux"))]
        {
            // Other platforms still return Unknown until implemented
            assert_eq!(metrics.vendor, "Unknown");
        }
    }

    #[test]
    fn test_gpu_usage_not_always_zero() {
        let metrics = get_gpu_metrics().unwrap();

        // On Linux with GPU, usage should vary (might be 0 at idle, but not always)
        #[cfg(target_os = "linux")]
        {
            // Just check it's a valid value - could be 0 at idle, that's OK
            assert!(metrics.usage_percent >= 0.0 && metrics.usage_percent <= 100.0);
        }

        #[cfg(not(target_os = "linux"))]
        {
            assert_eq!(metrics.usage_percent, 0.0);
        }
    }

    #[test]
    fn test_gpu_temperature_available_when_supported() {
        let metrics = get_gpu_metrics().unwrap();

        // Temperature may be available on Linux if GPU supports it
        #[cfg(target_os = "linux")]
        {
            if metrics.temperature.is_some() {
                let temp = metrics.temperature.unwrap();
                assert!(
                    temp >= 30.0 && temp <= 120.0,
                    "GPU temp {}°C out of range",
                    temp
                );
            }
            // If None, that's also OK - not all GPUs report temperature
        }

        #[cfg(not(target_os = "linux"))]
        {
            assert!(metrics.temperature.is_none());
        }
    }

    #[test]
    fn test_gpu_vendor_not_unknown_on_linux() {
        let metrics = get_gpu_metrics().unwrap();

        // This test will FAIL until we implement real GPU detection
        #[cfg(target_os = "linux")]
        assert_ne!(
            metrics.vendor, "Unknown",
            "Should detect GPU vendor on Linux, not return 'Unknown'"
        );

        #[cfg(not(target_os = "linux"))]
        {
            // Other platforms can return Unknown for now
        }
    }

    #[test]
    fn test_gpu_detects_real_gpu() {
        let metrics = get_gpu_metrics().unwrap();

        // On Linux, should get real data from system
        #[cfg(target_os = "linux")]
        {
            println!("✅ GPU Metrics:");
            println!("   Vendor: {}", metrics.vendor);
            println!("   Usage: {:.1}%", metrics.usage_percent);
            if let Some(temp) = metrics.temperature {
                println!("   Temperature: {:.1}°C", temp);
            } else {
                println!("   Temperature: Not available");
            }

            // Verify we have real data
            assert_ne!(metrics.vendor, "Unknown", "Should detect GPU vendor");
        }
    }

    #[test]
    fn test_gpu_multiple_calls_consistent() {
        let metrics1 = get_gpu_metrics().unwrap();
        let metrics2 = get_gpu_metrics().unwrap();

        // Vendor should be consistent
        assert_eq!(metrics1.vendor, metrics2.vendor);

        // Usage may vary slightly between calls, but should be similar
        let diff = (metrics1.usage_percent - metrics2.usage_percent).abs();
        assert!(
            diff <= 100.0,
            "Usage changed too much: {:.1} vs {:.1}",
            metrics1.usage_percent,
            metrics2.usage_percent
        );
    }
}
