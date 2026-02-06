// Protocol constants and helper functions for Divoom devices

pub const PROTOCOL_VERSION: &str = "1.0";

/// Check if device is TimeGate (hardware version 400)
pub fn is_timegate(hardware: &str) -> bool {
    hardware == "400"
}

/// Get screen count for device type
pub fn get_screen_count(hardware: &str) -> u32 {
    if is_timegate(hardware) {
        5
    } else {
        1
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_timegate_detection() {
        assert!(is_timegate("400"));
        assert!(!is_timegate("100"));
        assert!(!is_timegate("200"));
    }

    #[test]
    fn test_screen_count() {
        assert_eq!(get_screen_count("400"), 5);
        assert_eq!(get_screen_count("100"), 1);
    }
}
