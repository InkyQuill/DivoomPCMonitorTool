mod core;
mod config;
mod devices;
mod system;
mod ui;
mod i18n;

use config::{load_config, save_config as save_config_to_file, AppConfig};
use core::SystemMetrics;
use devices::{discover_devices, send_metrics, DivoomDevice};
use system::MetricsCollector;
use std::sync::Mutex;
use tauri::State;

// Global state
pub struct AppState {
    pub collector: Mutex<MetricsCollector>,
    pub config: Mutex<AppConfig>,
}

#[tauri::command]
async fn discover_devices_command() -> Result<Vec<DivoomDevice>, String> {
    discover_devices()
        .await
        .map_err(|e| e.to_string())
}

#[tauri::command]
async fn send_metrics_command(
    ip: String,
    lcd_id: i32,
    metrics: Vec<String>,
) -> Result<(), String> {
    send_metrics(&ip, lcd_id, metrics)
        .await
        .map_err(|e| e.to_string())
}

#[tauri::command]
fn get_config(state: State<'_, AppState>) -> Result<AppConfig, String> {
    let config = state.config.lock().unwrap();
    Ok(config.clone())
}

#[tauri::command]
fn save_config_command(
    config: AppConfig,
    state: State<'_, AppState>,
) -> Result<(), String> {
    save_config_to_file(&config).map_err(|e| e.to_string())?;

    let mut state_config = state.config.lock().unwrap();
    *state_config = config;

    Ok(())
}

#[tauri::command]
fn collect_metrics(state: State<'_, AppState>) -> Result<SystemMetrics, String> {
    let mut collector = state.collector.lock().unwrap();
    collector.collect().map_err(|e| e.to_string())
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    let config = load_config().unwrap_or_else(|_| {
        crate::config::default_config()
    });

    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_fs::init())
        .manage(AppState {
            collector: Mutex::new(MetricsCollector::new()),
            config: Mutex::new(config),
        })
        .invoke_handler(tauri::generate_handler![
            discover_devices_command,
            send_metrics_command,
            get_config,
            save_config_command,
            collect_metrics
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
