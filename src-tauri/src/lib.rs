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
use tauri::{Manager, State};
use ui::tray::{handle_tray_event, TrayIcon};

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
        .setup(|app| {
            // Setup system tray
            #[cfg(all(not(target_os = "android"), not(target_os = "ios")))]
            {
                use tauri::menu::{Menu, MenuItem};
                use tauri::tray::{TrayIconBuilder, TrayIconId};

                let tray_id = TrayIconId::new("main-tray");

                // Create menu items
                let show_item = MenuItem::with_id(app, "show", "Show", true, None::<String>)?;
                let hide_item = MenuItem::with_id(app, "hide", "Hide", true, None::<String>)?;
                let settings_item = MenuItem::with_id(app, "settings", "Settings", true, None::<String>)?;
                let quit_item = MenuItem::with_id(app, "quit", "Quit", true, None::<String>)?;

                // Build menu
                let menu = Menu::with_items(app, &[&show_item, &hide_item, &settings_item, &quit_item])?;

                // Build tray icon with menu
                let _tray = TrayIconBuilder::with_id(tray_id)
                    .menu(&menu)
                    .tooltip("Divoom PC Monitor")
                    .on_menu_event(|app, event| {
                        handle_tray_event(app, event.id().as_ref());
                    })
                    .build(app)?;
            }

            Ok(())
        })
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
