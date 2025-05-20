package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var logFile *os.File

// InitLogger инициализирует логгер
func InitLogger() error {
	// Создаем директорию для логов
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	logDir := filepath.Join(homeDir, ".divoom")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Открываем файл лога
	logPath := filepath.Join(logDir, "app.log")
	logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	return nil
}

// Log записывает сообщение в лог
func Log(format string, args ...interface{}) {
	if logFile == nil {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)

	if _, err := logFile.WriteString(logEntry); err != nil {
		fmt.Printf("Failed to write to log file: %v\n", err)
	}
}

// Close закрывает файл лога
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
