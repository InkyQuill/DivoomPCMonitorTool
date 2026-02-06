use super::models::{DeviceListResponse, DivoomDevice};
use crate::core::{AppError, Result};
use reqwest::Client;
use serde_json::Value;

const DISCOVERY_URL: &str = "http://app.divoom-gz.com/Device/ReturnSameLANDevice";

pub async fn discover_devices() -> Result<Vec<DivoomDevice>> {
    let client = Client::new();
    let response = client
        .get(DISCOVERY_URL)
        .timeout(std::time::Duration::from_secs(5))
        .send()
        .await
        .map_err(|_e| AppError::DeviceDiscoveryFailed)?;

    if !response.status().is_success() {
        return Err(AppError::DeviceDiscoveryFailed);
    }

    let json: Value = response.json().await?;

    let device_list = json["DeviceList"]
        .as_array()
        .ok_or_else(|| AppError::DeviceDiscoveryFailed)?;

    let devices: Vec<DivoomDevice> = device_list
        .iter()
        .filter_map(|d| serde_json::from_value(d.clone()).ok())
        .map(|d: DeviceListResponse| DivoomDevice {
            ip: d.ip,
            product_name: d.product_name,
            device_name: d.device_name,
            hardware: d.hardware,
            mac: d.mac,
        })
        .collect();

    Ok(devices)
}
