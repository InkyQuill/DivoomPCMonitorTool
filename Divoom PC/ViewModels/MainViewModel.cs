using System;
using System.Threading.Tasks;
using System.ComponentModel;
using System.Runtime.CompilerServices;
using DivoomPC.Services;
using DivoomPC.Models;
using DivoomPC.Views;
using System.Collections.ObjectModel;
using System.Windows.Input;
using Microsoft.UI.Xaml;
using System.Linq;

namespace DivoomPC.ViewModels
{
    public class MainViewModel : INotifyPropertyChanged
    {
        private readonly DivoomDeviceService _deviceService;
        private readonly SystemMonitorService _systemMonitor;
        private readonly SettingsService _settingsService;
        private readonly DispatcherTimer _updateTimer;
        private bool _isMinimizedToTray;

        public MainViewModel()
        {
            _deviceService = new DivoomDeviceService();
            _systemMonitor = new SystemMonitorService();
            _settingsService = new SettingsService();

            _updateTimer = new DispatcherTimer
            {
                Interval = TimeSpan.FromSeconds(1)
            };
            _updateTimer.Tick += (s, e) => UpdateTimer_Tick(s, e);

            Devices = new ObservableCollection<DivoomDeviceInfo>();
            MonitoringCards = new ObservableCollection<MonitoringCard>
            {
                new MonitoringCard("CPU Temperature"),
                new MonitoringCard("CPU Usage"),
                new MonitoringCard("GPU Temperature"),
                new MonitoringCard("GPU Usage"),
                new MonitoringCard("Memory Usage"),
                new MonitoringCard("HDD Temperature")
            };

            LCDOptions = new ObservableCollection<string> { "LCD 1", "LCD 2", "LCD 3", "LCD 4", "LCD 5" };
            SelectedDeviceIndex = -1;
            SelectedLCDIndex = -1;

            // Initialize commands
            RefreshCommand = new RelayCommand(async () => await UpdateDeviceListAsync());
            SettingsCommand = new RelayCommand(() => { });
            MinimizeCommand = new RelayCommand(() => { });
            MaximizeCommand = new RelayCommand(() => { });
            CloseCommand = new RelayCommand(() => { });
            ToggleAutoStartCommand = new RelayCommand(() => AutoStart = !AutoStart);
            ToggleStartMinimizedCommand = new RelayCommand(() => StartMinimized = !StartMinimized);
            SelectDeviceCommand = new RelayCommand<int>(index => SelectedDeviceIndex = index);
            SelectLCDCommand = new RelayCommand<int>(index => SelectedLCDIndex = index);

            LoadSettings();
            InitializeAsync();
        }

        private async void InitializeAsync()
        {
            await UpdateDeviceListAsync();
            if (_settingsService.StartMinimized)
            {
                MinimizeToTray();
            }
            _updateTimer.Start();
        }

        private void LoadSettings()
        {
            var savedDevice = _settingsService.SelectedDevice;
            var savedLCD = _settingsService.SelectedLCD;

            if (!string.IsNullOrEmpty(savedDevice))
            {
                var device = Devices.FirstOrDefault(d => d.DevicePrivateIP == savedDevice);
                if (device != null)
                {
                    var deviceIndex = Devices.IndexOf(device);
                    if (deviceIndex >= 0)
                    {
                        SelectedDeviceIndex = deviceIndex;
                        if (savedLCD >= 0)
                        {
                            SelectedLCDIndex = savedLCD;
                        }
                    }
                }
            }
        }

        private void UpdateTimer_Tick(object? sender, object e)
        {
            if (SelectedDeviceIndex >= 0 && SelectedLCDIndex >= 0 && Devices.Count > SelectedDeviceIndex)
            {
                _systemMonitor.Update();
                var cpuTemp = _systemMonitor.GetCpuTemperature();
                var cpuUse = _systemMonitor.GetCpuUsage();
                var gpuTemp = _systemMonitor.GetGpuTemperature();
                var gpuUse = _systemMonitor.GetGpuUsage();
                var memUse = _systemMonitor.GetMemoryUsage();
                var hddTemp = _systemMonitor.GetHardDiskTemperature();

                MonitoringCards[0].Value = $"{cpuTemp}°C";
                MonitoringCards[1].Value = $"{cpuUse}%";
                MonitoringCards[2].Value = $"{gpuTemp}°C";
                MonitoringCards[3].Value = $"{gpuUse}%";
                MonitoringCards[4].Value = $"{memUse}%";
                MonitoringCards[5].Value = $"{hddTemp}°C";

                _ = _deviceService.SendSystemInfoAsync(cpuTemp, cpuUse, gpuTemp, gpuUse, memUse, hddTemp);
            }
        }

        public ObservableCollection<DivoomDeviceInfo> Devices { get; }
        public ObservableCollection<MonitoringCard> MonitoringCards { get; }
        public ObservableCollection<string> LCDOptions { get; }

        private int _selectedDeviceIndex;
        public int SelectedDeviceIndex
        {
            get => _selectedDeviceIndex;
            set
            {
                if (_selectedDeviceIndex != value)
                {
                    _selectedDeviceIndex = value;
                    OnPropertyChanged();
                    if (value >= 0)
                    {
                        SelectDeviceAsync(value);
                    }
                }
            }
        }

        private int _selectedLCDIndex;
        public int SelectedLCDIndex
        {
            get => _selectedLCDIndex;
            set
            {
                if (_selectedLCDIndex != value)
                {
                    _selectedLCDIndex = value;
                    OnPropertyChanged();
                    if (value >= 0)
                    {
                        SelectLCDAsync(value);
                    }
                }
            }
        }

        public bool ShowLCDSelector => SelectedDeviceIndex >= 0 && 
                                     Devices.Count > SelectedDeviceIndex && 
                                     Devices[SelectedDeviceIndex].Hardware == 400;

        public bool IsMinimizedToTray
        {
            get => _isMinimizedToTray;
            set
            {
                if (_isMinimizedToTray != value)
                {
                    _isMinimizedToTray = value;
                    OnPropertyChanged();
                }
            }
        }

        public bool AutoStart
        {
            get => _settingsService.AutoStart;
            set
            {
                if (_settingsService.AutoStart != value)
                {
                    _settingsService.AutoStart = value;
                    _ = _settingsService.SetAutoStartAsync(value);
                    OnPropertyChanged();
                }
            }
        }

        public bool StartMinimized
        {
            get => _settingsService.StartMinimized;
            set
            {
                if (_settingsService.StartMinimized != value)
                {
                    _settingsService.StartMinimized = value;
                    OnPropertyChanged();
                }
            }
        }

        public ICommand RefreshCommand { get; }
        public ICommand SettingsCommand { get; }
        public ICommand MinimizeCommand { get; }
        public ICommand MaximizeCommand { get; }
        public ICommand CloseCommand { get; }
        public ICommand ToggleAutoStartCommand { get; }
        public ICommand ToggleStartMinimizedCommand { get; }
        public ICommand SelectDeviceCommand { get; }
        public ICommand SelectLCDCommand { get; }

        private void MinimizeToTray()
        {
            IsMinimizedToTray = true;
        }

        private void RestoreFromTray()
        {
            IsMinimizedToTray = false;
        }

        private void Exit()
        {
            _updateTimer.Stop();
            _systemMonitor.Dispose();
            Application.Current.Exit();
        }

        private async Task UpdateDeviceListAsync()
        {
            await _deviceService.UpdateDeviceListAsync();
            Devices.Clear();
            foreach (var device in _deviceService.LocalList.DeviceList)
            {
                Devices.Add(device);
            }

            if (Devices.Count == 1)
            {
                SelectedDeviceIndex = 0;
            }
        }

        private async void SelectDeviceAsync(int index)
        {
            await _deviceService.SelectDeviceAsync(index);
            if (index >= 0)
            {
                _settingsService.SelectedDevice = Devices[index].DevicePrivateIP;
                if (_deviceService.LcdIndependence == 0)
                {
                    SelectedLCDIndex = 0;
                }
                OnPropertyChanged(nameof(ShowLCDSelector));
            }
        }

        private async void SelectLCDAsync(int index)
        {
            await _deviceService.SelectLCDAsync(index);
            if (index >= 0)
            {
                _settingsService.SelectedLCD = index;
            }
        }

        public event PropertyChangedEventHandler? PropertyChanged;
        protected virtual void OnPropertyChanged([CallerMemberName] string? propertyName = null)
        {
            PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(propertyName));
        }
    }

    public class RelayCommand<T> : ICommand
    {
        private readonly Action<T> _execute;
        private readonly Func<T, bool>? _canExecute;

        public RelayCommand(Action<T> execute, Func<T, bool>? canExecute = null)
        {
            _execute = execute ?? throw new ArgumentNullException(nameof(execute));
            _canExecute = canExecute;
        }

        public event EventHandler? CanExecuteChanged
        {
            add { }
            remove { }
        }

        public bool CanExecute(object? parameter) => 
            parameter is T typedParameter && (_canExecute?.Invoke(typedParameter) ?? true);

        public void Execute(object? parameter)
        {
            if (parameter is T typedParameter)
            {
                _execute(typedParameter);
            }
        }
    }
} 