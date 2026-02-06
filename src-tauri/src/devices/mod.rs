pub mod discovery;
pub mod models;
pub mod protocol;
pub mod sender;

pub use discovery::discover_devices;
pub use models::*;
pub use protocol::{get_screen_count, is_timegate};
pub use sender::send_metrics;
