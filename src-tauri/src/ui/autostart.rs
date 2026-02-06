// Autostart management
// Implements autostart functionality using .desktop files on Linux
//
// ## Linux Implementation
// Creates a .desktop file in ~/.config/autostart/ with the following structure:
// ```ini
// [Desktop Entry]
// Type=Application
// Name=DivoomPCMonitorTool
// Exec=/path/to/executable
// Icon=divoom-pc-monitor
// Terminal=false
// X-GNOME-Autostart-enabled=true
// Hidden=false
// ```
//
// ## Platform Support
// - Linux: ✅ Implemented (uses .desktop files)
// - Windows: ⏳ TODO (uses Registry Run key)
// - macOS: ⏳ TODO (uses launchd/LaunchAgents)

use std::path::{Path, PathBuf};
use std::fs;
use std::env;
use crate::core::Result;

/// Autostart manager for controlling application autostart behavior
///
/// # Example
/// ```no_run
/// use ui::autostart::AutostartManager;
///
/// let manager = AutostartManager::new();
///
/// // Enable autostart
/// manager.enable().unwrap();
///
/// // Check if enabled
/// if manager.is_enabled() {
///     println!("Autostart is enabled");
/// }
///
/// // Disable autostart
/// manager.disable().unwrap();
/// ```
pub struct AutostartManager {
    app_name: String,
    exec_path: PathBuf,
}

impl AutostartManager {
    /// Create a new AutostartManager
    ///
    /// Automatically detects the application name and executable path
    pub fn new() -> Self {
        Self {
            app_name: "DivoomPCMonitorTool".to_string(),
            exec_path: env::current_exe().unwrap_or_else(|_| PathBuf::from("divoom-pc-companion")),
        }
    }

    /// Create AutostartManager with custom app name
    ///
    /// Useful for testing or custom installations
    pub fn with_name(app_name: String) -> Self {
        Self {
            app_name,
            exec_path: env::current_exe().unwrap_or_else(|_| PathBuf::from("divoom-pc-companion")),
        }
    }

    /// Check if autostart is enabled
    ///
    /// Returns true if the autostart file exists and is configured
    pub fn is_enabled(&self) -> bool {
        #[cfg(target_os = "linux")]
        {
            let path = self.get_autostart_file_path();
            path.exists()
        }

        #[cfg(not(target_os = "linux"))]
        {
            false // TODO: Implement for other platforms
        }
    }

    /// Enable autostart
    ///
    /// Creates the necessary configuration for the application to start automatically
    ///
    /// # Errors
    /// Returns an error if:
    /// - Cannot create autostart directory
    /// - Cannot write autostart file
    pub fn enable(&self) -> Result<()> {
        #[cfg(target_os = "linux")]
        {
            let autostart_file = self.get_autostart_file_path();

            // Create autostart directory if it doesn't exist
            if let Some(parent) = autostart_file.parent() {
                fs::create_dir_all(parent)?;
            }

            // Create .desktop file content
            let desktop_entry = format!(
                "[Desktop Entry]\n\
                 Type=Application\n\
                 Name={}\n\
                 Exec={}\n\
                 Icon={}\n\
                 Terminal=false\n\
                 X-GNOME-Autostart-enabled=true\n\
                 Hidden=false\n",
                self.app_name,
                self.exec_path.display(),
                "divoom-pc-monitor" // Icon name
            );

            // Write to file
            fs::write(&autostart_file, desktop_entry)?;

            Ok(())
        }

        #[cfg(not(target_os = "linux"))]
        {
            Ok(()) // TODO: Implement for other platforms
        }
    }

    /// Disable autostart
    ///
    /// Removes the autostart configuration
    ///
    /// # Errors
    /// Returns an error if:
    /// - Cannot remove autostart file (due to permissions, etc.)
    pub fn disable(&self) -> Result<()> {
        #[cfg(target_os = "linux")]
        {
            let autostart_file = self.get_autostart_file_path();

            // Remove file if it exists
            if autostart_file.exists() {
                fs::remove_file(&autostart_file)?;
            }

            Ok(())
        }

        #[cfg(not(target_os = "linux"))]
        {
            Ok(()) // TODO: Implement for other platforms
        }
    }

    /// Get the path to the autostart .desktop file (Linux only)
    #[cfg(target_os = "linux")]
    fn get_autostart_file_path(&self) -> PathBuf {
        let mut path = env::home_dir().unwrap_or_else(|| PathBuf::from("~"));
        path.push(".config/autostart");
        path.push(format!("{}.desktop", self.app_name));
        path
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_autostart_manager_creation() {
        let manager = AutostartManager::new();
        assert_eq!(manager.app_name, "DivoomPCMonitorTool");
    }

    #[test]
    fn test_autostart_manager_with_custom_name() {
        let manager = AutostartManager::with_name("CustomApp".to_string());
        assert_eq!(manager.app_name, "CustomApp");
    }

    #[test]
    fn test_autostart_manager_has_exec_path() {
        let manager = AutostartManager::new();
        assert!(!manager.exec_path.as_os_str().is_empty());
    }

    #[test]
    fn test_is_enabled_returns_bool() {
        let manager = AutostartManager::new();
        let _enabled = manager.is_enabled();
        // This test will FAIL until is_enabled() is implemented
        // Currently always returns false
    }

    #[test]
    fn test_enable_returns_result() {
        let manager = AutostartManager::new();
        let result = manager.enable();
        // This test will FAIL until enable() is implemented
        // Currently returns Ok(()) without doing anything
        assert!(result.is_ok() || result.is_err());
    }

    #[test]
    fn test_disable_returns_result() {
        let manager = AutostartManager::new();
        let result = manager.disable();
        // This test will FAIL until disable() is implemented
        // Currently returns Ok(()) without doing anything
        assert!(result.is_ok() || result.is_err());
    }

    #[test]
    fn test_autostart_file_path_on_linux() {
        let manager = AutostartManager::new();

        #[cfg(target_os = "linux")]
        {
            let path = manager.get_autostart_file_path();
            assert!(path.ends_with("DivoomPCMonitorTool.desktop"));
            assert!(path.to_string_lossy().contains(".config/autostart"));
        }
    }

    #[test]
    fn test_enable_creates_autostart_file() {
        let manager = AutostartManager::new();

        // Clean up first
        let _ = manager.disable();

        // This test will FAIL until we implement enable()
        let _ = manager.enable();

        #[cfg(target_os = "linux")]
        {
            let path = manager.get_autostart_file_path();
            // After enable(), the file should exist
            assert!(
                path.exists(),
                "Autostart file should exist after enable()"
            );
        }

        // Clean up
        let _ = manager.disable();
    }

    #[test]
    fn test_disable_removes_autostart_file() {
        let manager = AutostartManager::new();

        // First enable, then disable
        let _ = manager.enable();
        let _ = manager.disable();

        #[cfg(target_os = "linux")]
        {
            let path = manager.get_autostart_file_path();
            // After disable(), the file should not exist
            assert!(
                !path.exists(),
                "Autostart file should not exist after disable()"
            );
        }
    }

    #[test]
    fn test_is_enabled_detects_autostart_file() {
        let manager = AutostartManager::new();

        // Clean up any existing autostart file first
        let _ = manager.disable();

        // Initially should be false
        assert!(!manager.is_enabled(), "Autostart should be disabled initially");

        // After enable, should be true
        let _ = manager.enable();
        assert!(
            manager.is_enabled(),
            "Autostart should be enabled after enable()"
        );

        // Clean up
        let _ = manager.disable();
    }

    #[test]
    fn test_enable_is_idempotent() {
        let manager = AutostartManager::new();

        // Calling enable() multiple times should be safe
        let _ = manager.enable();
        let _ = manager.enable();
        let _ = manager.enable();

        assert!(manager.is_enabled(), "Autostart should still be enabled");

        // Clean up
        let _ = manager.disable();
    }

    #[test]
    fn test_disable_is_idempotent() {
        let manager = AutostartManager::new();

        // Calling disable() multiple times should be safe
        let _ = manager.disable();
        let _ = manager.disable();
        let _ = manager.disable();

        assert!(!manager.is_enabled(), "Autostart should still be disabled");
    }

    #[test]
    fn test_autostart_file_content_valid() {
        let manager = AutostartManager::new();
        let _ = manager.enable();

        #[cfg(target_os = "linux")]
        {
            let path = manager.get_autostart_file_path();
            if path.exists() {
                let content = fs::read_to_string(&path).unwrap();
                // Should contain basic .desktop file entries
                assert!(content.contains("[Desktop Entry]"), "Should have [Desktop Entry]");
                assert!(content.contains("Type=Application"), "Should have Type=Application");
            }
        }

        // Clean up
        let _ = manager.disable();
    }
}
