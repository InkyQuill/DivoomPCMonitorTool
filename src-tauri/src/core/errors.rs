use thiserror::Error;

#[derive(Error, Debug)]
pub enum AppError {
    #[error("Device discovery failed")]
    DeviceDiscoveryFailed,

    #[error("Device connection failed: {0}")]
    DeviceConnectionFailed(String),

    #[error("Data send failed: {0}")]
    DataSendFailed(String),

    #[error("Metric collection failed: {0}")]
    MetricCollectionFailed(String),

    #[error("Configuration error: {0}")]
    ConfigError(String),

    #[error("IO error: {0}")]
    IoError(#[from] std::io::Error),

    #[error("HTTP error: {0}")]
    HttpError(#[from] reqwest::Error),

    #[error("JSON error: {0}")]
    JsonError(#[from] serde_json::Error),

    #[error("GPU not available")]
    GpuNotAvailable,

    #[error("Unsupported GPU vendor: {0}")]
    UnsupportedGpuVendor(String),
}

pub type Result<T> = std::result::Result<T, AppError>;
