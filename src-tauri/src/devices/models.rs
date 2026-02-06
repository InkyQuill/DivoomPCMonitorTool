use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DivoomDevice {
    pub ip: String,
    pub product_name: String,
    pub device_name: String,
    pub hardware: String,
    pub mac: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DeviceListResponse {
    #[serde(rename = "DeviceName")]
    pub device_name: String,
    #[serde(rename = "ProductName")]
    pub product_name: String,
    #[serde(rename = "Hardware")]
    pub hardware: String,
    #[serde(rename = "Id")]
    pub id: String,
    #[serde(rename = "Ip")]
    pub ip: String,
    #[serde(rename = "Mac")]
    pub mac: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TimeGateScreenInfo {
    pub total_screens: u32,
    pub current_screen: u32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MetricsRequest {
    #[serde(rename = "Command")]
    pub command: String,
    #[serde(rename = "TextSec")]
    pub text_sec: i32,
    #[serde(rename = "LcdId")]
    pub lcd_id: i32,
    #[serde(rename = "TextContent")]
    pub text_content: Vec<String>,
}

impl MetricsRequest {
    pub fn new(lcd_id: i32, metrics: Vec<String>) -> Self {
        Self {
            command: "Device/UpdatePCParaInfo".to_string(),
            text_sec: 20,
            lcd_id,
            text_content: metrics,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_divoom_device_creation() {
        let device = DivoomDevice {
            ip: "192.168.1.100".to_string(),
            product_name: "Pixoo64".to_string(),
            device_name: "My Device".to_string(),
            hardware: "100".to_string(),
            mac: "AA:BB:CC:DD:EE:FF".to_string(),
        };

        assert_eq!(device.ip, "192.168.1.100");
        assert_eq!(device.product_name, "Pixoo64");
        assert_eq!(device.device_name, "My Device");
        assert_eq!(device.hardware, "100");
        assert_eq!(device.mac, "AA:BB:CC:DD:EE:FF");
    }

    #[test]
    fn test_divoom_device_serialization() {
        let device = DivoomDevice {
            ip: "192.168.1.100".to_string(),
            product_name: "Pixoo64".to_string(),
            device_name: "My Device".to_string(),
            hardware: "100".to_string(),
            mac: "AA:BB:CC:DD:EE:FF".to_string(),
        };

        let json = serde_json::to_string(&device);
        assert!(json.is_ok());

        let deserialized: DivoomDevice = serde_json::from_str(&json.unwrap()).unwrap();
        assert_eq!(deserialized.ip, device.ip);
        assert_eq!(deserialized.product_name, device.product_name);
    }

    #[test]
    fn test_device_list_response_deserialization() {
        let json = r#"{
            "DeviceName": "My Device",
            "ProductName": "Pixoo64",
            "Hardware": "100",
            "Id": "12345",
            "Ip": "192.168.1.100",
            "Mac": "AA:BB:CC:DD:EE:FF"
        }"#;

        let response: DeviceListResponse = serde_json::from_str(json).unwrap();
        assert_eq!(response.device_name, "My Device");
        assert_eq!(response.product_name, "Pixoo64");
        assert_eq!(response.hardware, "100");
        assert_eq!(response.ip, "192.168.1.100");
    }

    #[test]
    fn test_timegate_screen_info() {
        let info = TimeGateScreenInfo {
            total_screens: 5,
            current_screen: 2,
        };

        assert_eq!(info.total_screens, 5);
        assert_eq!(info.current_screen, 2);
    }

    #[test]
    fn test_metrics_request_new() {
        let metrics = vec!["CPU: 50%".to_string(), "RAM: 60%".to_string()];
        let request = MetricsRequest::new(0, metrics);

        assert_eq!(request.command, "Device/UpdatePCParaInfo");
        assert_eq!(request.text_sec, 20);
        assert_eq!(request.lcd_id, 0);
        assert_eq!(request.text_content.len(), 2);
        assert_eq!(request.text_content[0], "CPU: 50%");
    }

    #[test]
    fn test_metrics_request_serialization() {
        let metrics = vec!["CPU: 50%".to_string()];
        let request = MetricsRequest::new(0, metrics);

        let json = serde_json::to_string(&request);
        assert!(json.is_ok());

        let json_str = json.unwrap();
        assert!(json_str.contains("Device/UpdatePCParaInfo"));
        assert!(json_str.contains("TextSec"));
        assert!(json_str.contains("LcdId"));
        assert!(json_str.contains("TextContent"));
    }

    #[test]
    fn test_metrics_request_with_different_lcd_id() {
        let request1 = MetricsRequest::new(0, vec![]);
        let request2 = MetricsRequest::new(3, vec![]);

        assert_eq!(request1.lcd_id, 0);
        assert_eq!(request2.lcd_id, 3);
    }

    #[test]
    fn test_metrics_request_with_empty_metrics() {
        let request = MetricsRequest::new(0, vec![]);
        assert_eq!(request.text_content.len(), 0);
    }

    #[test]
    fn test_metrics_request_with_multiple_metrics() {
        let metrics = vec![
            "CPU: 50%".to_string(),
            "GPU: 60%".to_string(),
            "RAM: 70%".to_string(),
            "Storage: 80%".to_string(),
        ];
        let request = MetricsRequest::new(0, metrics);

        assert_eq!(request.text_content.len(), 4);
        assert_eq!(request.text_content[3], "Storage: 80%");
    }
}
