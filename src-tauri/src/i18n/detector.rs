// System language detection
// TODO: Implement platform-specific language detection

pub fn detect_system_language() -> String {
    // Default to English
    // TODO: Detect system language on Linux, macOS, Windows
    "en".to_string()
}

pub fn get_supported_languages() -> Vec<String> {
    vec![
        "en".to_string(),
        "ru".to_string(),
        "zh-CN".to_string(),
        "de".to_string(),
        "es".to_string(),
        "fr".to_string(),
        "ja".to_string(),
        "ko".to_string(),
        "pt".to_string(),
        "tr".to_string(),
        "ar".to_string(),
        "it".to_string(),
        "pl".to_string(),
        "nl".to_string(),
        "vi".to_string(),
        "th".to_string(),
    ]
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_detect_system_language_returns_valid_string() {
        let lang = detect_system_language();
        assert!(!lang.is_empty());
        assert_eq!(lang, "en"); // Current implementation always returns "en"
    }

    #[test]
    fn test_get_supported_languages_not_empty() {
        let languages = get_supported_languages();
        assert!(!languages.is_empty());
    }

    #[test]
    fn test_get_supported_languages_count() {
        let languages = get_supported_languages();
        assert_eq!(languages.len(), 16);
    }

    #[test]
    fn test_get_supported_languages_contains_english() {
        let languages = get_supported_languages();
        assert!(languages.contains(&"en".to_string()));
    }

    #[test]
    fn test_get_supported_languages_contains_russian() {
        let languages = get_supported_languages();
        assert!(languages.contains(&"ru".to_string()));
    }

    #[test]
    fn test_get_supported_languages_contains_chinese() {
        let languages = get_supported_languages();
        assert!(languages.contains(&"zh-CN".to_string()));
    }

    #[test]
    fn test_get_supported_languages_all_valid_format() {
        let languages = get_supported_languages();

        for lang in languages {
            assert!(!lang.is_empty());
            assert!(lang.len() <= 6); // Language codes should be short
            assert!(lang.contains('-') || lang.len() == 2); // Either "xx" or "xx-XX" format
        }
    }

    #[test]
    fn test_get_supported_languages_no_duplicates() {
        let languages = get_supported_languages();
        let unique: std::collections::HashSet<_> = languages.iter().cloned().collect();
        assert_eq!(languages.len(), unique.len());
    }
}
