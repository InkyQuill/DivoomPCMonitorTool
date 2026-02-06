pub mod detector;
pub mod loader;

pub use detector::{detect_system_language, get_supported_languages};
pub use loader::{get_translation, load_translations};
