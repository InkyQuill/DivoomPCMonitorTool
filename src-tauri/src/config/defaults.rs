use super::schema::*;

pub fn default_config() -> AppConfig {
    AppConfig {
        general: GeneralConfig {
            language: "en".to_string(),
            start_minimized: false,
            autostart: false,
            update_interval: 1000,
        },
        metrics: MetricsConfig {
            enable_cpu: true,
            enable_gpu: true,
            enable_memory: true,
            enable_storage: true,
            storage_path: "/".to_string(),
        },
        display: DisplayConfig {
            screen_index: 0,
            brightness: 100,
            theme_mode: "system".to_string(),
        },
        device: DeviceConfig {
            ip_address: String::new(),
            device_type: "pixoo64".to_string(),
            connection_timeout: 5000,
            retry_attempts: 3,
        },
        advanced: AdvancedConfig {
            log_level: "info".to_string(),
            log_file: "".to_string(),
        },
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_default_config_general() {
        let config = default_config();
        assert_eq!(config.general.language, "en");
        assert_eq!(config.general.start_minimized, false);
        assert_eq!(config.general.autostart, false);
        assert_eq!(config.general.update_interval, 1000);
    }

    #[test]
    fn test_default_config_metrics() {
        let config = default_config();
        assert_eq!(config.metrics.enable_cpu, true);
        assert_eq!(config.metrics.enable_gpu, true);
        assert_eq!(config.metrics.enable_memory, true);
        assert_eq!(config.metrics.enable_storage, true);
        assert_eq!(config.metrics.storage_path, "/");
    }

    #[test]
    fn test_default_config_display() {
        let config = default_config();
        assert_eq!(config.display.screen_index, 0);
        assert_eq!(config.display.brightness, 100);
        assert_eq!(config.display.theme_mode, "system");
    }

    #[test]
    fn test_default_config_device() {
        let config = default_config();
        assert_eq!(config.device.ip_address, "");
        assert_eq!(config.device.device_type, "pixoo64");
        assert_eq!(config.device.connection_timeout, 5000);
        assert_eq!(config.device.retry_attempts, 3);
    }

    #[test]
    fn test_default_config_advanced() {
        let config = default_config();
        assert_eq!(config.advanced.log_level, "info");
        assert_eq!(config.advanced.log_file, "");
    }

    #[test]
    fn test_default_config_serialization() {
        let config = default_config();
        let json = serde_json::to_string(&config);
        assert!(json.is_ok());

        // Can deserialize back
        let deserialized: AppConfig = serde_json::from_str(&json.unwrap()).unwrap();
        assert_eq!(deserialized.general.language, config.general.language);
    }
}
