use super::models::MetricsRequest;
use crate::core::Result;
use reqwest::Client;

pub async fn send_metrics(device_ip: &str, lcd_id: i32, metrics: Vec<String>) -> Result<()> {
    let client = Client::new();
    let url = format!("http://{}:80/post", device_ip);

    let request = MetricsRequest::new(lcd_id, metrics);

    client
        .post(&url)
        .timeout(std::time::Duration::from_secs(5))
        .json(&request)
        .send()
        .await
        .map_err(|e| crate::core::AppError::DataSendFailed(e.to_string()))?;

    Ok(())
}
