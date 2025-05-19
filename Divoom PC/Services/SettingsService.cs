using System;
using System.Threading.Tasks;
using Windows.Storage;
using Microsoft.Windows.AppLifecycle;
using Windows.ApplicationModel;

namespace DivoomPC.Services
{
    public class SettingsService
    {
        private readonly ApplicationDataContainer _localSettings;
        private const string AutoStartKey = "AutoStart";
        private const string StartMinimizedKey = "StartMinimized";
        private const string SelectedDeviceKey = "SelectedDevice";
        private const string SelectedLCDKey = "SelectedLCD";

        public SettingsService()
        {
            _localSettings = ApplicationData.Current.LocalSettings;
        }

        public bool AutoStart
        {
            get => GetSetting<bool>(AutoStartKey, false);
            set => SetSetting(AutoStartKey, value);
        }

        public bool StartMinimized
        {
            get => GetSetting<bool>(StartMinimizedKey, false);
            set => SetSetting(StartMinimizedKey, value);
        }

        public string SelectedDevice
        {
            get => GetSetting<string>(SelectedDeviceKey, string.Empty);
            set => SetSetting(SelectedDeviceKey, value);
        }

        public int SelectedLCD
        {
            get => GetSetting<int>(SelectedLCDKey, -1);
            set => SetSetting(SelectedLCDKey, value);
        }

        private T GetSetting<T>(string key, T defaultValue)
        {
            if (_localSettings.Values.TryGetValue(key, out object? value) && value != null)
            {
                return (T)value;
            }
            return defaultValue;
        }

        private void SetSetting<T>(string key, T value)
        {
            _localSettings.Values[key] = value;
        }

        public async Task SetAutoStartAsync(bool enable)
        {
            try
            {
                var startupTask = await StartupTask.GetAsync("DivoomPCStartupTask");
                
                if (enable)
                {
                    var state = await startupTask.RequestEnableAsync();
                    if (state != StartupTaskState.Enabled)
                    {
                        throw new Exception("Failed to enable startup task");
                    }
                }
                else
                {
                    startupTask.Disable();
                }
            }
            catch (Exception ex)
            {
                throw new Exception($"Failed to configure autostart: {ex.Message}");
            }
        }
    }
} 