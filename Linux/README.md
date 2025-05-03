# Divoom PC Monitor for Linux

A simple application for displaying system information on Divoom devices.

## Installation

1. Install Python 3.6 or higher
2. Install dependencies:
```bash
pip install -r requirements.txt
```
3. Make the script executable:
```bash
chmod +x divoom.py
```

## Usage

### Device Configuration
```bash
./divoom.py --configure
```

### Start Monitoring
```bash
./divoom.py
```

## Running as a Service

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

## Notes

- GPU monitoring requires additional setup depending on your graphics card:
  - NVIDIA: Install NVIDIA drivers and `nvidia-smi`
  - AMD: Install ROCm and `rocm-smi`
  - Intel: Install `intel-gpu-tools` package
- Configuration is saved in `~/.divoom_config.json`
- The script will automatically detect and use available GPU monitoring tools
- If no GPU is detected or monitoring tools are not available, GPU metrics will show as "--"
