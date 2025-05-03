# Divoom PC Monitor

A cross-platform tool for monitoring and displaying PC system information on Divoom devices (Pixoo64 and TimeGate).

## Project Structure

- `DivoomPCMonitorTool.zip` - Windows executable (.NET Framework 4.8, compatible with Windows 10/11)
- `Windows/` - .NET source code for Windows version
- `Linux/` - Python script and instructions for Linux version

## Windows Version

The Windows version is a .NET Framework 4.8 application that provides a graphical interface for monitoring and displaying system information on Divoom devices.

### Requirements
- Windows 10 or 11
- .NET Framework 4.8
- Divoom device (Pixoo64 or TimeGate) connected to the same network

### Features
- Real-time system monitoring
- Support for both single-screen and 5-screen devices
- Automatic device discovery
- System tray integration
- Auto-start option

## Linux Version

The Linux version is a Python script that provides similar functionality through the command line.

### Requirements
- Python 3.6 or higher
- Divoom device (Pixoo64 or TimeGate) connected to the same network
- Required system packages:
  - `lm-sensors` - for CPU temperature monitoring
  - `psutil` - for system monitoring
  - Optional GPU monitoring tools:
    - NVIDIA: `nvidia-smi` (installed with NVIDIA drivers)
    - AMD: `rocm-smi` (part of ROCm)
    - Intel: `intel-gpu-top` (part of intel-gpu-tools)

### Installation
1. Install Python 3.6 or higher
2. Install system packages:
```bash
# For Debian/Ubuntu:
sudo apt-get install lm-sensors
sudo sensors-detect  # Configure sensors

# For Fedora:
sudo dnf install lm_sensors
sudo sensors-detect  # Configure sensors

# For Arch Linux:
sudo pacman -S lm-sensors
sudo sensors-detect  # Configure sensors
```
3. Install Python dependencies:
```bash
pip install -r Linux/requirements.txt
```
4. Make the script executable:
```bash
chmod +x Linux/divoom.py
```

### Usage
- Configure device:
```bash
./Linux/divoom.py --configure
```
- Start monitoring:
```bash
./Linux/divoom.py
```

### Running as a Service
1. Create service file `/etc/systemd/system/divoom.service`:
```ini
[Unit]
Description=Divoom PC Monitor
After=network.target

[Service]
Type=simple
User=YOUR_USERNAME
ExecStart=/usr/bin/python3 /path/to/divoom.py
Restart=always

[Install]
WantedBy=multi-user.target
```

2. Start the service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable divoom
sudo systemctl start divoom
```

## Supported Devices
- Pixoo64
- TimeGate

## Features
- CPU temperature and usage monitoring
- GPU temperature and usage monitoring (if available)
- Memory usage monitoring
- Disk usage monitoring
- Automatic device discovery
- Configuration persistence
- Cross-platform support

## Notes
- Configuration is saved in `~/.divoom_config.json` on Linux
- GPU monitoring requires appropriate drivers and tools
- The application must be on the same network as the Divoom device
