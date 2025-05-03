using System;
using System.Management;
using System.Runtime.InteropServices;
using System.Diagnostics;
using System.Linq;
using LibreHardwareMonitor.Hardware;

namespace DivoomPCMonitor
{
    public class SystemMonitor
    {
        private readonly Computer computer;
        private readonly UpdateVisitor updateVisitor;
        private readonly PerformanceCounter cpuCounter;
        private readonly PerformanceCounter ramCounter;
        private readonly bool isLibreHardwareMonitorAvailable;

        public SystemMonitor()
        {
            // Инициализация счетчиков производительности
            try
            {
                cpuCounter = new PerformanceCounter("Processor", "% Processor Time", "_Total");
                ramCounter = new PerformanceCounter("Memory", "Available MBytes");
                cpuCounter.NextValue(); // Первое значение всегда 0
            }
            catch (Exception)
            {
                // Если счетчики недоступны, продолжим без них
            }

            // Инициализация LibreHardwareMonitor
            try
            {
                computer = new Computer();
                computer.IsCpuEnabled = true;
                computer.IsGpuEnabled = true;
                computer.IsMemoryEnabled = true;
                computer.IsStorageEnabled = true;
                computer.Open();
                updateVisitor = new UpdateVisitor();
                isLibreHardwareMonitorAvailable = true;
            }
            catch (Exception)
            {
                isLibreHardwareMonitorAvailable = false;
            }
        }

        public void Update()
        {
            if (isLibreHardwareMonitorAvailable)
            {
                try
                {
                    computer.Accept(updateVisitor);
                }
                catch (Exception)
                {
                    // Игнорируем ошибки обновления
                }
            }
        }

        public string GetCpuTemperature()
        {
            if (isLibreHardwareMonitorAvailable)
            {
                try
                {
                    foreach (var hardware in computer.Hardware)
                    {
                        if (hardware.HardwareType == HardwareType.Cpu)
                        {
                            hardware.Update();
                            foreach (var sensor in hardware.Sensors)
                            {
                                if (sensor.SensorType == SensorType.Temperature)
                                {
                                    return sensor.Value?.ToString("0") + "°C" ?? "--";
                                }
                            }
                        }
                    }
                }
                catch (Exception)
                {
                    // Продолжаем с альтернативными методами
                }
            }

            // Альтернативный метод через WMI
            try
            {
                using (var searcher = new ManagementObjectSearcher(@"root\WMI", "SELECT * FROM MSAcpi_ThermalZoneTemperature"))
                {
                    foreach (ManagementObject obj in searcher.Get())
                    {
                        double temp = Convert.ToDouble(obj["CurrentTemperature"].ToString());
                        return ((temp - 2732) / 10.0).ToString("0") + "°C";
                    }
                }
            }
            catch (Exception)
            {
                // Продолжаем с другими методами
            }

            return "--";
        }

        public string GetCpuUsage()
        {
            try
            {
                if (cpuCounter != null)
                {
                    return cpuCounter.NextValue().ToString("0") + "%";
                }
            }
            catch (Exception)
            {
                // Продолжаем с альтернативными методами
            }

            if (isLibreHardwareMonitorAvailable)
            {
                try
                {
                    foreach (var hardware in computer.Hardware)
                    {
                        if (hardware.HardwareType == HardwareType.Cpu)
                        {
                            hardware.Update();
                            foreach (var sensor in hardware.Sensors)
                            {
                                if (sensor.SensorType == SensorType.Load)
                                {
                                    return sensor.Value?.ToString("0") + "%" ?? "--";
                                }
                            }
                        }
                    }
                }
                catch (Exception)
                {
                    // Продолжаем с другими методами
                }
            }

            return "--";
        }

        public string GetGpuTemperature()
        {
            if (isLibreHardwareMonitorAvailable)
            {
                try
                {
                    foreach (var hardware in computer.Hardware)
                    {
                        if (hardware.HardwareType == HardwareType.GpuNvidia || 
                            hardware.HardwareType == HardwareType.GpuAmd ||
                            hardware.HardwareType == HardwareType.GpuIntel)
                        {
                            hardware.Update();
                            foreach (var sensor in hardware.Sensors)
                            {
                                if (sensor.SensorType == SensorType.Temperature)
                                {
                                    return sensor.Value?.ToString("0") + "°C" ?? "--";
                                }
                            }
                        }
                    }
                }
                catch (Exception)
                {
                    // Продолжаем с другими методами
                }
            }

            return "--";
        }

        public string GetGpuUsage()
        {
            if (isLibreHardwareMonitorAvailable)
            {
                try
                {
                    foreach (var hardware in computer.Hardware)
                    {
                        if (hardware.HardwareType == HardwareType.GpuNvidia || 
                            hardware.HardwareType == HardwareType.GpuAmd ||
                            hardware.HardwareType == HardwareType.GpuIntel)
                        {
                            hardware.Update();
                            foreach (var sensor in hardware.Sensors)
                            {
                                if (sensor.SensorType == SensorType.Load)
                                {
                                    return sensor.Value?.ToString("0") + "%" ?? "--";
                                }
                            }
                        }
                    }
                }
                catch (Exception)
                {
                    // Продолжаем с другими методами
                }
            }

            return "--";
        }

        public string GetMemoryUsage()
        {
            try
            {
                if (ramCounter != null)
                {
                    float availableMB = ramCounter.NextValue();
                    float totalMB = new Microsoft.VisualBasic.Devices.ComputerInfo().TotalPhysicalMemory / (1024 * 1024);
                    float usedPercent = ((totalMB - availableMB) / totalMB) * 100;
                    return usedPercent.ToString("0") + "%";
                }
            }
            catch (Exception)
            {
                // Продолжаем с альтернативными методами
            }

            if (isLibreHardwareMonitorAvailable)
            {
                try
                {
                    foreach (var hardware in computer.Hardware)
                    {
                        if (hardware.HardwareType == HardwareType.Memory)
                        {
                            hardware.Update();
                            foreach (var sensor in hardware.Sensors)
                            {
                                if (sensor.SensorType == SensorType.Load)
                                {
                                    return sensor.Value?.ToString("0") + "%" ?? "--";
                                }
                            }
                        }
                    }
                }
                catch (Exception)
                {
                    // Продолжаем с другими методами
                }
            }

            // Используем Win32 API как последний вариант
            try
            {
                MEMORYSTATUSEX memInfo = new MEMORYSTATUSEX();
                memInfo.dwLength = (uint)Marshal.SizeOf(typeof(MEMORYSTATUSEX));
                GlobalMemoryStatusEx(ref memInfo);
                return memInfo.dwMemoryLoad.ToString() + "%";
            }
            catch (Exception)
            {
                return "--";
            }
        }

        public string GetHardDiskTemperature()
        {
            if (isLibreHardwareMonitorAvailable)
            {
                try
                {
                    foreach (var hardware in computer.Hardware)
                    {
                        if (hardware.HardwareType == HardwareType.Storage)
                        {
                            hardware.Update();
                            foreach (var sensor in hardware.Sensors)
                            {
                                if (sensor.SensorType == SensorType.Temperature)
                                {
                                    return sensor.Value?.ToString("0") + "°C" ?? "--";
                                }
                            }
                        }
                    }
                }
                catch (Exception)
                {
                    // Продолжаем с другими методами
                }
            }

            return "--";
        }

        public void Close()
        {
            if (isLibreHardwareMonitorAvailable)
            {
                try
                {
                    computer.Close();
                }
                catch (Exception)
                {
                    // Игнорируем ошибки закрытия
                }
            }

            if (cpuCounter != null)
            {
                cpuCounter.Dispose();
            }
            if (ramCounter != null)
            {
                ramCounter.Dispose();
            }
        }

        [StructLayout(LayoutKind.Sequential)]
        public struct MEMORYSTATUSEX
        {
            public uint dwLength;
            public uint dwMemoryLoad;
            public ulong ullTotalPhys;
            public ulong ullAvailPhys;
            public ulong ullTotalPageFile;
            public ulong ullAvailPageFile;
            public ulong ullTotalVirtual;
            public ulong ullAvailVirtual;
            public ulong ullAvailExtendedVirtual;
        }

        [DllImport("kernel32.dll")]
        public static extern void GlobalMemoryStatusEx(ref MEMORYSTATUSEX stat);
    }
} 