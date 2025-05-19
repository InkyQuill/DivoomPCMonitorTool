using System;
using System.Threading.Tasks;
using LibreHardwareMonitor.Hardware;
using System.Linq;

namespace DivoomPC.Services
{
    public class SystemMonitorService : IDisposable
    {
        private readonly Computer _computer;
        private readonly UpdateVisitor _updateVisitor;

        public SystemMonitorService()
        {
            _computer = new Computer
            {
                IsCpuEnabled = true,
                IsGpuEnabled = true,
                IsMemoryEnabled = true,
                IsStorageEnabled = true
            };
            _computer.Open();
            _updateVisitor = new UpdateVisitor();
        }

        public void Update()
        {
            _computer.Accept(_updateVisitor);
        }

        public string GetCpuTemperature()
        {
            var cpu = _computer.Hardware.FirstOrDefault(h => h.HardwareType == HardwareType.Cpu);
            if (cpu == null) return "N/A";

            var tempSensor = cpu.Sensors.FirstOrDefault(s => s.SensorType == SensorType.Temperature);
            return tempSensor?.Value?.ToString("F1") ?? "N/A";
        }

        public string GetCpuUsage()
        {
            var cpu = _computer.Hardware.FirstOrDefault(h => h.HardwareType == HardwareType.Cpu);
            if (cpu == null) return "N/A";

            var loadSensor = cpu.Sensors.FirstOrDefault(s => s.SensorType == SensorType.Load);
            return loadSensor?.Value?.ToString("F1") ?? "N/A";
        }

        public string GetGpuTemperature()
        {
            var gpu = _computer.Hardware.FirstOrDefault(h => h.HardwareType == HardwareType.GpuNvidia || 
                                                           h.HardwareType == HardwareType.GpuAmd);
            if (gpu == null) return "N/A";

            var tempSensor = gpu.Sensors.FirstOrDefault(s => s.SensorType == SensorType.Temperature);
            return tempSensor?.Value?.ToString("F1") ?? "N/A";
        }

        public string GetGpuUsage()
        {
            var gpu = _computer.Hardware.FirstOrDefault(h => h.HardwareType == HardwareType.GpuNvidia || 
                                                           h.HardwareType == HardwareType.GpuAmd);
            if (gpu == null) return "N/A";

            var loadSensor = gpu.Sensors.FirstOrDefault(s => s.SensorType == SensorType.Load);
            return loadSensor?.Value?.ToString("F1") ?? "N/A";
        }

        public string GetMemoryUsage()
        {
            var memory = _computer.Hardware.FirstOrDefault(h => h.HardwareType == HardwareType.Memory);
            if (memory == null) return "N/A";

            var loadSensor = memory.Sensors.FirstOrDefault(s => s.SensorType == SensorType.Load);
            return loadSensor?.Value?.ToString("F1") ?? "N/A";
        }

        public string GetHardDiskTemperature()
        {
            var storage = _computer.Hardware.FirstOrDefault(h => h.HardwareType == HardwareType.Storage);
            if (storage == null) return "N/A";

            var tempSensor = storage.Sensors.FirstOrDefault(s => s.SensorType == SensorType.Temperature);
            return tempSensor?.Value?.ToString("F1") ?? "N/A";
        }

        public void Dispose()
        {
            _computer.Close();
        }
    }

    public class UpdateVisitor : IVisitor
    {
        public void VisitComputer(IComputer computer)
        {
            computer.Traverse(this);
        }

        public void VisitHardware(IHardware hardware)
        {
            hardware.Update();
            foreach (IHardware subHardware in hardware.SubHardware)
                subHardware.Accept(this);
        }

        public void VisitSensor(ISensor sensor) { }

        public void VisitParameter(IParameter parameter) { }
    }
} 