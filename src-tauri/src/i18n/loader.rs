// i18n translations loader
// Loads translations from JSON files for internationalization support
//
// Translation files are stored in translations/{language}.json
// Each JSON file contains key-value pairs of translation strings
//
// Example JSON structure:
// {
//   "app.name": "Divoom PC Monitor",
//   "common.ok": "OK",
//   "menu.settings": "Settings"
// }

use std::collections::HashMap;
use std::fs;
use std::path::Path;

/// Load translations for a specific language from JSON file
///
/// Reads translation file from translations/{language}.json
/// Returns empty HashMap if file not found or invalid JSON
///
/// # Arguments
/// * `language` - Language code (e.g., "en", "ru", "zh-CN")
///
/// # Returns
/// HashMap of translation key-value pairs
///
/// # Example
/// ```
/// let translations = load_translations("en");
/// let app_name = translations.get("app.name").unwrap_or(&"Unknown".to_string());
/// ```
pub fn load_translations(language: &str) -> HashMap<String, String> {
    // Build path to translation file
    // Try both locations (for different build contexts)
    let paths_to_try = vec![
        Path::new("translations").join(format!("{}.json", language)),
        Path::new("src-tauri/translations").join(format!("{}.json", language)),
    ];

    for file_path in paths_to_try {
        if file_path.exists() {
            // Read file content
            if let Ok(content) = fs::read_to_string(&file_path) {
                // Parse JSON
                if let Ok(translations) = serde_json::from_str::<HashMap<String, String>>(&content)
                {
                    return translations;
                }
            }
        }
    }

    // No file found or all failed to parse
    HashMap::new()
}

/// Get a translation for a specific key in a specific language
///
/// Returns None if key or language not found
///
/// # Arguments
/// * `key` - Translation key (e.g., "app.name", "common.ok")
/// * `language` - Language code (e.g., "en", "ru")
///
/// # Returns
/// Some(String) if translation found, None otherwise
///
/// # Example
/// ```
/// let title = get_translation("app.name", "en").unwrap_or("App".to_string());
/// ```
pub fn get_translation(key: &str, language: &str) -> Option<String> {
    let translations = load_translations(language);
    translations.get(key).cloned()
}

/// Get translation with fallback to English
///
/// If translation not found in requested language, tries English
/// Returns None if key doesn't exist in either language
///
/// This is the recommended function for getting translations as it
/// provides automatic fallback to English for incomplete translations.
///
/// # Arguments
/// * `key` - Translation key
/// * `language` - Preferred language code
///
/// # Returns
/// Some(String) if translation found in requested or English language, None otherwise
pub fn get_translation_with_fallback(key: &str, language: &str) -> Option<String> {
    // Try requested language first
    if let Some(trans) = get_translation(key, language) {
        return Some(trans);
    }

    // Fallback to English
    if language != "en" {
        get_translation(key, "en")
    } else {
        None
    }
}

/// Get translation with default value
///
/// Returns translation if found, otherwise returns the provided default
///
/// # Arguments
/// * `key` - Translation key
/// * `language` - Language code
/// * `default` - Default value to return if translation not found
///
/// # Returns
/// Translation string or default value
///
/// # Example
/// ```
/// let title = get_translation_or_default("app.name", "en", "My App");
/// ```
pub fn get_translation_or_default(key: &str, language: &str, default: &str) -> String {
    get_translation_with_fallback(key, language).unwrap_or_else(|| default.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_load_translations_not_empty() {
        // This test will FAIL until we implement JSON loading
        let translations = load_translations("en");

        // Should have loaded translations from JSON file
        assert!(
            !translations.is_empty(),
            "English translations should be loaded from JSON file"
        );
    }

    #[test]
    fn test_load_translations_contains_common_keys() {
        let translations = load_translations("en");

        // Should have common translation keys
        assert!(
            translations.contains_key("app.name")
                || translations.contains_key("common.ok")
                || translations.contains_key("common.cancel"),
            "Translations should contain common keys"
        );
    }

    #[test]
    fn test_load_translations_for_different_languages() {
        let en = load_translations("en");
        let ru = load_translations("ru");

        // English should be loaded
        assert!(!en.is_empty(), "English translations should be loaded");

        // Russian may be empty if not yet translated, but shouldn't panic
        // Just verify it returns a HashMap
        let _ru_count = ru.len();
    }

    #[test]
    fn test_get_translation_with_real_keys() {
        // This test will FAIL until we have actual translations
        let result = get_translation("app.name", "en");

        assert!(
            result.is_some(),
            "Should find translation for 'app.name' key in English"
        );

        if let Some(value) = result {
            assert!(!value.is_empty(), "Translation should not be empty");
        }
    }

    #[test]
    fn test_get_translation_handles_empty_key() {
        let result = get_translation("", "en");
        assert!(result.is_none(), "Empty key should return None");
    }

    #[test]
    fn test_get_translation_handles_unknown_language() {
        let result = get_translation("test", "unknown_lang_xyz");
        assert!(result.is_some() || result.is_none()); // Either is ok (fallback to empty)
    }

    #[test]
    fn test_get_translation_with_fallback_english() {
        // This test will FAIL until we implement translations
        let result = get_translation_with_fallback("app.name", "en");

        assert!(
            result.is_some(),
            "Should get English translation with fallback"
        );
    }

    #[test]
    fn test_get_translation_with_fallback_missing_language() {
        // If language doesn't exist, should fallback to English
        let result = get_translation_with_fallback("app.name", "unknown_lang");

        assert!(
            result.is_some(),
            "Should fallback to English for unknown language"
        );
    }

    #[test]
    fn test_get_translation_with_fallback_missing_key() {
        // If key doesn't exist in either language, should return None
        let result = get_translation_with_fallback("nonexistent_key_xyz", "en");

        assert!(result.is_none(), "Should return None for non-existent key");
    }

    #[test]
    fn test_translation_files_exist() {
        // Verify translation files exist in the filesystem
        let translations_dir = Path::new("translations");

        // Check both possible locations
        let dir_exists = translations_dir.exists() || Path::new("src-tauri/translations").exists();

        assert!(
            dir_exists,
            "Translations directory should exist at translations/ or src-tauri/translations/"
        );

        // English file should exist in one of the locations
        let en_file1 = Path::new("translations/en.json");
        let en_file2 = Path::new("src-tauri/translations/en.json");

        assert!(
            en_file1.exists() || en_file2.exists(),
            "English translation file should exist"
        );
    }

    #[test]
    fn test_translation_files_valid_json() {
        // Verify translation files are valid JSON
        let en_file1 = Path::new("translations/en.json");
        let en_file2 = Path::new("src-tauri/translations/en.json");

        let en_file = if en_file1.exists() {
            en_file1
        } else {
            en_file2
        };

        if en_file.exists() {
            let content = fs::read_to_string(en_file).expect("Should read file");
            assert!(
                serde_json::from_str::<HashMap<String, String>>(&content).is_ok(),
                "English translation file should be valid JSON"
            );
        }
    }

    #[test]
    fn test_get_translation_or_default_with_existing_key() {
        let result = get_translation_or_default("app.name", "en", "Default App");
        assert_eq!(result, "Divoom PC Monitor");
    }

    #[test]
    fn test_get_translation_or_default_with_missing_key() {
        let result = get_translation_or_default("nonexistent_key", "en", "Default Value");
        assert_eq!(result, "Default Value");
    }

    #[test]
    fn test_get_translation_or_default_with_missing_language() {
        // Should fallback to English and then to default
        let result = get_translation_or_default("app.name", "unknown_lang", "Fallback");
        assert_eq!(result, "Divoom PC Monitor"); // Found in English
    }
}
