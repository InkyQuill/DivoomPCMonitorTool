package autostart

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	AppName = "System Monitor"
)

// IsEnabled проверяет, включен ли автозапуск
func IsEnabled() bool {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("schtasks", "/query", "/tn", AppName)
		return cmd.Run() == nil
	case "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return false
		}
		desktopFile := filepath.Join(homeDir, ".config", "autostart", AppName+".desktop")
		_, err = os.Stat(desktopFile)
		return err == nil
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return false
		}
		plistFile := filepath.Join(homeDir, "Library", "LaunchAgents", AppName+".plist")
		_, err = os.Stat(plistFile)
		return err == nil
	default:
		return false
	}
}

// Enable включает автозапуск приложения
func Enable() error {
	switch runtime.GOOS {
	case "windows":
		return enableWindows()
	case "linux":
		return enableLinux()
	case "darwin":
		return enableDarwin()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// Disable отключает автозапуск приложения
func Disable() error {
	switch runtime.GOOS {
	case "windows":
		return disableWindows()
	case "linux":
		return disableLinux()
	case "darwin":
		return disableDarwin()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// --- Windows ---
func enableWindows() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// Создаем XML для Task Scheduler
	xmlContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
    <Description>%s Autostart</Description>
  </RegistrationInfo>
  <Triggers>
    <LogonTrigger>
      <Enabled>true</Enabled>
    </LogonTrigger>
  </Triggers>
  <Principals>
    <Principal id="Author">
      <LogonType>InteractiveToken</LogonType>
      <RunLevel>HighestAvailable</RunLevel>
    </Principal>
  </Principals>
  <Settings>
    <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
    <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>false</StopIfGoingOnBatteries>
    <AllowHardTerminate>true</AllowHardTerminate>
    <StartWhenAvailable>true</StartWhenAvailable>
    <RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
    <IdleSettings>
      <StopOnIdleEnd>false</StopOnIdleEnd>
      <RestartOnIdle>false</RestartOnIdle>
    </IdleSettings>
    <AllowStartOnDemand>true</AllowStartOnDemand>
    <Enabled>true</Enabled>
    <Hidden>false</Hidden>
    <RunOnlyIfIdle>false</RunOnlyIfIdle>
    <WakeToRun>false</WakeToRun>
    <ExecutionTimeLimit>PT0S</ExecutionTimeLimit>
    <Priority>7</Priority>
  </Settings>
  <Actions Context="Author">
    <Exec>
      <Command>%s</Command>
    </Exec>
  </Actions>
</Task>`, AppName, exePath)

	// Создаем временный XML файл
	tmpFile := filepath.Join(os.TempDir(), "autostart.xml")
	if err := os.WriteFile(tmpFile, []byte(xmlContent), 0644); err != nil {
		return fmt.Errorf("failed to create task XML: %v", err)
	}
	defer os.Remove(tmpFile)

	// Создаем задачу в Task Scheduler
	cmd := exec.Command("schtasks", "/create", "/tn", AppName, "/xml", tmpFile, "/f")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create scheduled task: %v", err)
	}

	return nil
}

func disableWindows() error {
	cmd := exec.Command("schtasks", "/delete", "/tn", AppName, "/f")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete scheduled task: %v", err)
	}
	return nil
}

// --- Linux ---
func enableLinux() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// Создаем .desktop файл
	desktopContent := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Exec="%s"
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
`, AppName, exePath)

	autostartDir := filepath.Join(homeDir, ".config", "autostart")
	if err := os.MkdirAll(autostartDir, 0755); err != nil {
		return fmt.Errorf("failed to create autostart directory: %v", err)
	}

	desktopFile := filepath.Join(autostartDir, AppName+".desktop")
	if err := os.WriteFile(desktopFile, []byte(desktopContent), 0644); err != nil {
		return fmt.Errorf("failed to create desktop file: %v", err)
	}

	return nil
}

func disableLinux() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	desktopFile := filepath.Join(homeDir, ".config", "autostart", AppName+".desktop")
	if err := os.Remove(desktopFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove desktop file: %v", err)
	}

	return nil
}

// --- MacOS ---
func enableDarwin() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// Создаем .plist файл
	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>%s</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>`, AppName, exePath)

	launchAgentsDir := filepath.Join(homeDir, "Library", "LaunchAgents")
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		return fmt.Errorf("failed to create LaunchAgents directory: %v", err)
	}

	plistFile := filepath.Join(launchAgentsDir, AppName+".plist")
	if err := os.WriteFile(plistFile, []byte(plistContent), 0644); err != nil {
		return fmt.Errorf("failed to create plist file: %v", err)
	}

	return nil
}

func disableDarwin() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	plistFile := filepath.Join(homeDir, "Library", "LaunchAgents", AppName+".plist")
	if err := os.Remove(plistFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove plist file: %v", err)
	}

	return nil
}
