use super::{schema::AppConfig, defaults::default_config};
use crate::core::{AppError, Result};
use std::fs;
use std::path::PathBuf;

pub fn get_config_path() -> Result<PathBuf> {
    let config_dir = if cfg!(target_os = "linux") {
        dirs::config_dir()
            .ok_or_else(|| AppError::ConfigError("Cannot find config directory".to_string()))?
            .join("divoom-pc-companion")
    } else if cfg!(target_os = "macos") {
        dirs::config_dir()
            .ok_or_else(|| AppError::ConfigError("Cannot find config directory".to_string()))?
            .join("DivoomPCCompanion")
    } else if cfg!(target_os = "windows") {
        dirs::config_dir()
            .ok_or_else(|| AppError::ConfigError("Cannot find config directory".to_string()))?
            .join("DivoomPCCompanion")
    } else {
        return Err(AppError::ConfigError("Unsupported platform".to_string()));
    };

    // Create directory if it doesn't exist
    fs::create_dir_all(&config_dir)?;

    Ok(config_dir.join("config.json"))
}

pub fn load_config() -> Result<AppConfig> {
    let config_path = get_config_path()?;

    if !config_path.exists() {
        let default_config = default_config();
        save_config(&default_config)?;
        return Ok(default_config);
    }

    let content = fs::read_to_string(&config_path)?;
    let config: AppConfig = serde_json::from_str(&content)?;
    Ok(config)
}

pub fn save_config(config: &AppConfig) -> Result<()> {
    let config_path = get_config_path()?;
    let content = serde_json::to_string_pretty(config)?;
    fs::write(&config_path, content)?;
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_get_config_path_returns_valid_path() {
        let path = get_config_path();
        assert!(path.is_ok());

        let path_buf = path.unwrap();
        assert!(path_buf.ends_with("config.json"));
    }

    #[test]
    fn test_get_config_path_creates_directory() {
        let path = get_config_path().unwrap();
        let parent = path.parent();

        assert!(parent.is_some());
        assert!(parent.unwrap().exists());
    }

    #[test]
    fn test_default_config_roundtrip() {
        let original = default_config();
        save_config(&original).unwrap();
        let loaded = load_config().unwrap();

        assert_eq!(original.general.language, loaded.general.language);
        assert_eq!(original.metrics.enable_cpu, loaded.metrics.enable_cpu);
        assert_eq!(original.device.device_type, loaded.device.device_type);
    }

    #[test]
    fn test_save_and_load_modified_config() {
        // Note: This test may fail when run in parallel with other tests
        // because all tests share the same config file. In production,
        // use a test-specific config path or run tests serially with
        // --test-threads=1

        let mut config = default_config();
        config.general.language = "test_lang".to_string(); // Use unique value
        config.general.update_interval = 9999;

        save_config(&config).unwrap();

        // Small delay to ensure file is written
        std::thread::sleep(std::time::Duration::from_millis(10));

        let loaded = load_config().unwrap();

        // Verify the values were saved
        assert_eq!(loaded.general.language, "test_lang");
        assert_eq!(loaded.general.update_interval, 9999);
    }

    #[test]
    fn test_load_config_creates_default_if_missing() {
        // Remove existing config if present
        if let Ok(path) = get_config_path() {
            let _ = fs::remove_file(&path);
        }

        let config = load_config().unwrap();
        assert_eq!(config.general.language, "en"); // Default value
    }

    #[test]
    fn test_config_json_is_valid() {
        let config = default_config();
        let json = serde_json::to_string_pretty(&config);
        assert!(json.is_ok());

        // Verify JSON is not empty
        let json_str = json.unwrap();
        assert!(!json_str.is_empty());
        assert!(json_str.contains("{"));
        assert!(json_str.contains("}"));
    }

    #[test]
    fn test_save_config_overwrites_existing() {
        let config1 = default_config();
        save_config(&config1).unwrap();

        let mut config2 = default_config();
        config2.general.language = "de".to_string();
        save_config(&config2).unwrap();

        let loaded = load_config().unwrap();
        assert_eq!(loaded.general.language, "de");
    }
}
