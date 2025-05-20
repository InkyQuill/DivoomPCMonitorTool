package metrics

import "runtime"

// NewCollector создает новый коллектор метрик в зависимости от операционной системы
func NewCollector() Collector {
	switch runtime.GOOS {
	case "linux":
		return NewLinuxCollector()
	case "windows":
		return NewWindowsCollector()
	case "darwin":
		return NewDarwinCollector()
	default:
		return NewLinuxCollector() // По умолчанию используем Linux коллектор
	}
}
