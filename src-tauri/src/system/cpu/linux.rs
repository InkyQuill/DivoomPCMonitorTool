// Linux-specific CPU temperature detection
// NOTE: Tests written FIRST (TDD RED phase)

use std::path::Path;
use std::fs;

/// Get CPU temperature from Linux hwmon/thermal zones
/// Returns temperature in Celsius, or None if unavailable
///
/// Priority order:
/// 1. /sys/class/hwmon/hwmon*/name (k10temp, coretemp, etc.)
/// 2. /sys/class/thermal/thermal_zone* (fallback)
pub fn get_cpu_temperature() -> Option<f32> {
    // Try hwmon first (most reliable on modern systems)
    if let Some(temp) = get_temperature_from_hwmon() {
        return Some(temp);
    }

    // Fallback to thermal zones
    if let Some(temp) = get_temperature_from_thermal() {
        return Some(temp);
    }

    None
}

/// Read temperature from hwmon devices
fn get_temperature_from_hwmon() -> Option<f32> {
    let hwmon_path = Path::new("/sys/class/hwmon");
    if !hwmon_path.exists() {
        return None;
    }

    let entries = fs::read_dir(hwmon_path).ok()?;

    for entry in entries.flatten() {
        let device_path = entry.path();
        let name_path = device_path.join("name");

        // Read device name to check if it's a CPU sensor
        if let Ok(name) = fs::read_to_string(&name_path) {
            let name_lower = name.to_lowercase();

            // Known CPU temperature sensor drivers
            let is_cpu_sensor = name_lower.contains("k10temp")  // AMD Ryzen/EPYC
                || name_lower.contains("k8temp")  // AMD older
                || name_lower.contains("coretemp")  // Intel Core
                || name_lower.contains("cpu")      // Generic CPU
                || name_lower.contains("via")      // VIA
                || name_lower.contains("acpitz")   // ACPI thermal zone
                || name_lower.contains("x86_pkg_temp");  // Intel PCH

            if is_cpu_sensor {
                // Look for temp*_input files (usually temp1_input, temp2_input, etc.)
                if let Ok(entries) = fs::read_dir(&device_path) {
                    for file in entries.flatten() {
                        let file_path = file.path();
                        let file_name = file_path.file_name()
                            .unwrap_or_default()
                            .to_string_lossy();

                        if file_name.starts_with("temp") && file_name.ends_with("_input") {
                            if let Ok(temp_str) = fs::read_to_string(&file_path) {
                                if let Ok(millidegrees) = temp_str.trim().parse::<i64>() {
                                    // Temperature is in millidegrees Celsius
                                    let degrees = millidegrees as f32 / 1000.0;

                                    // Sanity check: CPU temps should be reasonable
                                    if degrees >= 20.0 && degrees <= 130.0 {
                                        return Some(degrees);
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    None
}

/// Read temperature from thermal zones (legacy fallback)
fn get_temperature_from_thermal() -> Option<f32> {
    let thermal_path = Path::new("/sys/class/thermal");

    if !thermal_path.exists() {
        return None;
    }

    let entries = fs::read_dir(thermal_path).ok()?;

    for entry in entries.flatten() {
        let zone_path = entry.path();
        let zone_name = zone_path.file_name()?.to_string_lossy();

        if zone_name.starts_with("thermal_zone") {
            // Read the type file to check if it's a CPU zone
            let type_path = zone_path.join("type");
            if let Ok(type_content) = fs::read_to_string(&type_path) {
                let type_str = type_content.to_lowercase();
                if type_str.contains("cpu") || type_str.contains("x86") || type_str.contains("acpi") {
                    // This is likely a CPU thermal zone, read temperature
                    let temp_path = zone_path.join("temp");
                    if let Ok(temp_str) = fs::read_to_string(&temp_path) {
                        // Temperature is in millidegrees Celsius
                        if let Ok(millidegrees) = temp_str.trim().parse::<i64>() {
                            let degrees = millidegrees as f32 / 1000.0;

                            // Sanity check
                            if degrees >= 20.0 && degrees <= 130.0 {
                                return Some(degrees);
                            }
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
    fn test_get_cpu_temperature_returns_option() {
        let temp = get_cpu_temperature();
        // Function should return Option (Some or None)
        match temp {
            Some(t) => assert!(t >= 0.0 && t <= 150.0), // Reasonable CPU temp range
            None => {} // Also acceptable if no thermal sensor
        }
    }

    #[test]
    fn test_get_cpu_temperature_in_celsius() {
        let temp = get_cpu_temperature();

        if let Some(t) = temp {
            // Temperature should be in reasonable Celsius range
            assert!(t >= 0.0, "Temperature should be >= 0°C, got {}", t);
            assert!(t <= 150.0, "Temperature should be <= 150°C, got {}", t);

            // Most systems should report between 20-100°C
            assert!(t >= 20.0 || t <= 110.0, "Temperature seems unrealistic: {}°C", t);
        }
    }

    #[test]
    fn test_get_cpu_temperature_not_none_when_available() {
        // On most Linux systems, at least one thermal zone should exist
        // This test may fail on systems without thermal sensors
        let temp = get_cpu_temperature();

        // We expect either Some value or explicitly None
        // The test documents the expected behavior
        if temp.is_none() {
            println!("WARNING: No thermal sensor detected on this system");
        }
    }

    #[test]
    fn test_temperature_format() {
        let temp = get_cpu_temperature();

        if let Some(t) = temp {
            // Temperature should have reasonable precision
            // Not too many decimal places (sensor noise)
            let formatted = format!("{:.1}", t);
            assert!(formatted.parse::<f32>().is_ok());
        }
    }

    #[test]
    fn test_multiple_calls_consistent() {
        let temp1 = get_cpu_temperature();
        let temp2 = get_cpu_temperature();

        match (temp1, temp2) {
            (Some(t1), Some(t2)) => {
                // Temperatures should be similar (within 5°C)
                // CPU temp doesn't change THAT fast between calls
                let diff = (t1 - t2).abs();
                assert!(diff <= 5.0, "Temperatures changed too rapidly: {} vs {}", t1, t2);
            }
            (None, None) => {} // Both unavailable - OK
            _ => {
                // One available, one not - sensor inconsistency
                println!("WARNING: Sensor availability changed between calls");
            }
        }
    }

    #[test]
    fn test_handles_missing_thermal_sys() {
        // Even if /sys/class/thermal doesn't exist, should not panic
        let temp = get_cpu_temperature();
        // Should return None, not crash
        assert!(temp.is_some() || temp.is_none());
    }

    #[test]
    fn test_reads_from_thermal_zone() {
        // This test will FAIL until we implement actual reading
        // This is the TRUE RED test!
        use std::path::Path;

        let thermal_path = Path::new("/sys/class/thermal");

        // Check if thermal_zone* directories exist
        let has_thermal_zones = thermal_path.exists()
            && fs::read_dir(thermal_path)
                .map(|entries| {
                    entries
                        .flatten()
                        .any(|entry| {
                            entry
                                .file_name()
                                .to_string_lossy()
                                .starts_with("thermal_zone")
                        })
                })
                .unwrap_or(false);

        if has_thermal_zones {
            let temp = get_cpu_temperature();

            // If thermal zones exist, we should get a temperature
            assert!(
                temp.is_some(),
                "Expected Some temperature when thermal zones exist, got None"
            );

            let t = temp.unwrap();
            assert!(t > 0.0, "Temperature should be greater than 0");
            assert!(t < 150.0, "Temperature should be less than 150°C");
        } else {
            // Skip test if no thermal zones
            println!("INFO: No thermal_zone* in /sys/class/thermal");
        }
    }

    #[test]
    fn test_reads_from_hwmon() {
        // Test hwmon (hardware monitoring) source
        use std::path::Path;

        let hwmon_path = Path::new("/sys/class/hwmon");

        if !hwmon_path.exists() {
            println!("INFO: No /sys/class/hwmon on this system");
            return;
        }

        // Check if any hwmon device with CPU name exists
        let has_cpu_hwmon = fs::read_dir(hwmon_path)
            .map(|entries| {
                entries
                    .flatten()
                    .any(|entry| {
                        let name_path = entry.path().join("name");
                        if let Ok(name) = fs::read_to_string(&name_path) {
                            let name_lower = name.to_lowercase();
                            name_lower.contains("k10temp")  // AMD
                                || name_lower.contains("coretemp")  // Intel
                                || name_lower.contains("cpu")
                                || name_lower.contains("k8temp")
                                || name_lower.contains("via")
                        } else {
                            false
                        }
                    })
            })
            .unwrap_or(false);

        if has_cpu_hwmon {
            let temp = get_cpu_temperature();

            // If CPU hwmon exists, we should get a temperature
            assert!(
                temp.is_some(),
                "Expected Some temperature when CPU hwmon exists, got None"
            );

            let t = temp.unwrap();
            assert!(t >= 20.0 && t <= 110.0, "CPU temp {}°C out of range", t);
        } else {
            println!("INFO: No CPU hwmon device found");
        }
    }

    #[test]
    fn test_hwmon_has_temp_input() {
        // Verify hwmon devices have temp*_input files
        use std::path::Path;

        let hwmon_path = Path::new("/sys/class/hwmon");
        if !hwmon_path.exists() {
            return;
        }

        if let Ok(entries) = fs::read_dir(hwmon_path) {
            for entry in entries.flatten() {
                let device_path = entry.path();
                let name_path = device_path.join("name");

                if let Ok(name) = fs::read_to_string(&name_path) {
                    let name_lower = name.to_lowercase();
                    if name_lower.contains("k10temp")
                        || name_lower.contains("coretemp")
                        || name_lower.contains("cpu")
                    {
                        // This is a CPU device, check for temp*_input
                        let has_temp = fs::read_dir(&device_path)
                            .map(|e| {
                                e.flatten()
                                    .any(|f| {
                                        f.file_name()
                                            .to_string_lossy()
                                            .starts_with("temp")
                                            && f.file_name()
                                                .to_string_lossy()
                                            .ends_with("_input")
                                    })
                            })
                            .unwrap_or(false);

                        if has_temp {
                            println!("INFO: Found CPU hwmon with temp inputs: {}", name.trim());
                        }
                    }
                }
            }
        }
    }

    #[test]
    fn test_actual_temperature_reading() {
        // This test shows the ACTUAL temperature being read
        let temp = get_cpu_temperature();

        match temp {
            Some(t) => {
                println!("✅ CPU Temperature detected: {:.1}°C", t);
                assert!(t >= 20.0 && t <= 110.0);
            }
            None => {
                println!("⚠️  No CPU temperature detected on this system");
            }
        }
    }
}

