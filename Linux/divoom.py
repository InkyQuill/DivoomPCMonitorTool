#!/usr/bin/env python3
import argparse
import json
import os
import sys
import time
import requests
import psutil
import subprocess
import logging
from typing import Optional, Dict, List

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

class DivoomDeviceManager:
    def __init__(self):
        self.device_ip = None
        self.device_type = None
        self.selected_lcd = 0
        self.config_file = os.path.expanduser("~/.divoom_config.json")
        self.load_config()

    def load_config(self):
        try:
            if os.path.exists(self.config_file):
                with open(self.config_file, 'r') as f:
                    config = json.load(f)
                    self.device_ip = config.get('device_ip')
                    self.device_type = config.get('device_type')
                    self.selected_lcd = config.get('selected_lcd', 0)
        except json.JSONDecodeError:
            logger.error("Invalid configuration file format")
        except Exception as e:
            logger.error(f"Error loading configuration: {e}")

    def save_config(self):
        try:
            config = {
                'device_ip': self.device_ip,
                'device_type': self.device_type,
                'selected_lcd': self.selected_lcd
            }
            with open(self.config_file, 'w') as f:
                json.dump(config, f)
        except Exception as e:
            logger.error(f"Error saving configuration: {e}")

    def discover_devices(self) -> List[Dict]:
        try:
            response = requests.get("http://app.divoom-gz.com/Device/ReturnSameLANDevice", timeout=5)
            if response.status_code == 200:
                return response.json().get('DeviceList', [])
            else:
                logger.error(f"Failed to discover devices. Status code: {response.status_code}")
        except requests.exceptions.Timeout:
            logger.error("Device discovery timed out")
        except requests.exceptions.RequestException as e:
            logger.error(f"Network error during device discovery: {e}")
        except json.JSONDecodeError:
            logger.error("Invalid response format from device discovery")
        except Exception as e:
            logger.error(f"Unexpected error during device discovery: {e}")
        return []

    def configure(self):
        devices = self.discover_devices()
        if not devices:
            print("No Divoom devices found in the network")
            return

        print("\nFound devices:")
        for i, device in enumerate(devices, 1):
            print(f"{i}. {device['DeviceName']} ({device['DevicePrivateIP']})")

        while True:
            try:
                choice = int(input("\nSelect device (number): "))
                if 1 <= choice <= len(devices):
                    selected = devices[choice - 1]
                    self.device_ip = selected['DevicePrivateIP']
                    self.device_type = selected['Hardware']
                    break
                print("Invalid choice")
            except ValueError:
                print("Please enter a number")

        if self.device_type == 400:  # 5-screen device
            while True:
                try:
                    lcd = int(input("\nSelect display (1-5): "))
                    if 1 <= lcd <= 5:
                        self.selected_lcd = lcd - 1
                        break
                    print("Invalid choice")
                except ValueError:
                    print("Please enter a number")

        self.save_config()
        print("\nConfiguration saved")

    def send_system_info(self):
        if not self.device_ip:
            logger.warning("Device not configured. Run with --configure parameter")
            return

        try:
            cpu_temp = self.get_cpu_temperature()
            cpu_usage = self.get_cpu_usage()
            memory_usage = self.get_memory_usage()
            gpu_temp = self.get_gpu_temperature()
            gpu_usage = self.get_gpu_usage()
            disk_usage = self.get_disk_usage()

            data = {
                "Command": "Device/UpdatePCParaInfo",
                "ScreenList": [{
                    "LcdId": self.selected_lcd,
                    "DispData": [
                        cpu_usage,  # CPU Usage
                        gpu_usage,  # GPU Usage
                        cpu_temp,   # CPU Temp
                        gpu_temp,   # GPU Temp
                        memory_usage,  # Memory Usage
                        disk_usage     # Disk Usage
                    ]
                }]
            }

            response = requests.post(
                f"http://{self.device_ip}:80/post",
                json=data,
                timeout=2
            )
            if response.status_code != 200:
                logger.error(f"Failed to send data to device. Status code: {response.status_code}")
        except requests.exceptions.Timeout:
            logger.error("Timeout while sending data to device")
        except requests.exceptions.RequestException as e:
            logger.error(f"Network error while sending data to device: {e}")
        except Exception as e:
            logger.error(f"Unexpected error while sending system info: {e}")

    def get_cpu_temperature(self) -> str:
        # Method 1: Using psutil (primary method)
        try:
            temps = psutil.sensors_temperatures()
            # Check various possible sensor names
            for sensor_name in ['coretemp', 'k10temp', 'zenpower', 'acpitz']:
                if sensor_name in temps:
                    # Get maximum temperature from all cores
                    max_temp = max(temp.current for temp in temps[sensor_name])
                    return f"{max_temp:.0f}°C"
        except Exception as e:
            logger.debug(f"Error getting temperature via psutil: {e}")

        # Method 2: Using sysfs (for Linux)
        try:
            with open('/sys/class/thermal/thermal_zone0/temp', 'r') as f:
                temp = float(f.read().strip()) / 1000.0
                return f"{temp:.0f}°C"
        except Exception as e:
            logger.debug(f"Error getting temperature via sysfs: {e}")

        # Method 3: Using lm-sensors
        try:
            result = subprocess.run(['sensors', '-j'], capture_output=True, text=True)
            if result.returncode == 0:
                data = json.loads(result.stdout)
                # Check various possible sensor names
                for chip in data.values():
                    for key, value in chip.items():
                        if 'temp' in key.lower() and 'input' in key:
                            temp = float(value)
                            return f"{temp:.0f}°C"
        except Exception as e:
            logger.debug(f"Error getting temperature via lm-sensors: {e}")

        return "--"

    def get_cpu_usage(self) -> str:
        try:
            return f"{psutil.cpu_percent(interval=1):.0f}%"
        except Exception as e:
            logger.warning(f"Error getting CPU usage: {e}")
            return "--"

    def get_memory_usage(self) -> str:
        try:
            return f"{psutil.virtual_memory().percent:.0f}%"
        except Exception as e:
            logger.warning(f"Error getting memory usage: {e}")
            return "--"

    def get_gpu_temperature(self) -> str:
        # Try NVIDIA first
        try:
            result = subprocess.run(['nvidia-smi', '--query-gpu=temperature.gpu', '--format=csv,noheader'],
                                  capture_output=True, text=True)
            if result.returncode == 0:
                temp = result.stdout.strip()
                if temp.isdigit():
                    return f"{temp}°C"
        except Exception as e:
            logger.debug(f"NVIDIA GPU temperature check failed: {e}")

        # Try AMD
        try:
            result = subprocess.run(['rocm-smi', '--showtemp', '--json'],
                                  capture_output=True, text=True)
            if result.returncode == 0:
                data = json.loads(result.stdout)
                if 'card0' in data and 'Temperature' in data['card0']:
                    temp = data['card0']['Temperature']
                    return f"{temp}°C"
        except Exception as e:
            logger.debug(f"AMD GPU temperature check failed: {e}")

        # Try Intel
        try:
            result = subprocess.run(['intel_gpu_top', '-J'],
                                  capture_output=True, text=True)
            if result.returncode == 0:
                data = json.loads(result.stdout)
                if 'engines' in data and 'render' in data['engines']:
                    temp = data['engines']['render'].get('temperature', 0)
                    return f"{temp}°C"
        except Exception as e:
            logger.debug(f"Intel GPU temperature check failed: {e}")

        return "--"

    def get_gpu_usage(self) -> str:
        # Try NVIDIA first
        try:
            result = subprocess.run(['nvidia-smi', '--query-gpu=utilization.gpu', '--format=csv,noheader'],
                                  capture_output=True, text=True)
            if result.returncode == 0:
                usage = result.stdout.strip().replace(' %', '')
                if usage.isdigit():
                    return f"{usage}%"
        except Exception as e:
            logger.debug(f"NVIDIA GPU usage check failed: {e}")

        # Try AMD
        try:
            result = subprocess.run(['rocm-smi', '--showuse', '--json'],
                                  capture_output=True, text=True)
            if result.returncode == 0:
                data = json.loads(result.stdout)
                if 'card0' in data and 'GPU use (%)' in data['card0']:
                    usage = data['card0']['GPU use (%)']
                    return f"{usage}%"
        except Exception as e:
            logger.debug(f"AMD GPU usage check failed: {e}")

        # Try Intel
        try:
            result = subprocess.run(['intel_gpu_top', '-J'],
                                  capture_output=True, text=True)
            if result.returncode == 0:
                data = json.loads(result.stdout)
                if 'engines' in data and 'render' in data['engines']:
                    usage = data['engines']['render'].get('busy', 0)
                    return f"{usage}%"
        except Exception as e:
            logger.debug(f"Intel GPU usage check failed: {e}")

        return "--"

    def get_disk_usage(self) -> str:
        try:
            return f"{psutil.disk_usage('/').percent:.0f}%"
        except Exception as e:
            logger.warning(f"Error getting disk usage: {e}")
            return "--"

def main():
    parser = argparse.ArgumentParser(description='Divoom PC Monitor for Linux')
    parser.add_argument('--configure', action='store_true', help='Configure device')
    args = parser.parse_args()

    try:
        manager = DivoomDeviceManager()

        if args.configure:
            manager.configure()
            return

        while True:
            try:
                manager.send_system_info()
                time.sleep(1)
            except Exception as e:
                logger.error(f"Error in main loop: {e}")
                time.sleep(5)  # Wait before retrying
    except KeyboardInterrupt:
        logger.info("Program stopped by user")
    except Exception as e:
        logger.error(f"Fatal error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main() 