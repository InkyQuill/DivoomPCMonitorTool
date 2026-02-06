export const SUPPORTED_LANGUAGES = [
  { code: 'en', name: 'English' },
  { code: 'ru', name: 'Русский' },
  { code: 'zh-CN', name: '简体中文' },
  { code: 'de', name: 'Deutsch' },
  { code: 'es', name: 'Español' },
  { code: 'fr', name: 'Français' },
  { code: 'ja', name: '日本語' },
  { code: 'ko', name: '한국어' },
  { code: 'pt', name: 'Português' },
  { code: 'tr', name: 'Türkçe' },
  { code: 'ar', name: 'العربية' },
  { code: 'it', name: 'Italiano' },
  { code: 'pl', name: 'Polski' },
  { code: 'nl', name: 'Nederlands' },
  { code: 'vi', name: 'Tiếng Việt' },
  { code: 'th', name: 'ไทย' },
];

export const DEVICE_TYPES = [
  { value: 'pixoo64', name: 'Pixoo 64' },
  { value: 'timegate', name: 'TimeGate' },
  { value: 'pixoo16', name: 'Pixoo 16' },
  { value: 'pixoo-max', name: 'Pixoo Max' },
];

export const UPDATE_INTERVALS = [
  { value: 500, name: '0.5s' },
  { value: 1000, name: '1s' },
  { value: 2000, name: '2s' },
  { value: 5000, name: '5s' },
];

export const DEFAULT_CONFIG = {
  general: {
    language: 'en',
    start_minimized: false,
    autostart: false,
    update_interval: 1000,
  },
  metrics: {
    enable_cpu: true,
    enable_gpu: true,
    enable_memory: true,
    enable_storage: true,
    storage_path: '/',
  },
  display: {
    screen_index: 0,
    brightness: 100,
    theme_mode: 'system',
  },
  device: {
    ip_address: '',
    device_type: 'pixoo64',
    connection_timeout: 5000,
    retry_attempts: 3,
  },
  advanced: {
    log_level: 'info',
    log_file: '',
  },
};
