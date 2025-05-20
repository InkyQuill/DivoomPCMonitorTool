package metrics

import (
	"os/exec"
	"strconv"
	"strings"
)

// GetHDDTempFromSmartctl получает температуру HDD через smartctl
func GetHDDTempFromSmartctl(device string) float64 {
	out, err := exec.Command("smartctl", "-A", device).Output()
	if err != nil {
		return 0.0
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Temperature_Celsius") {
			fields := strings.Fields(line)
			if len(fields) >= 10 {
				temp, err := strconv.ParseFloat(fields[9], 64)
				if err == nil {
					return temp
				}
			}
		}
	}

	return 0.0
}
