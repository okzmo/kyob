package utils

import (
	"fmt"
	"math"
)

func BytesToHuman(bytes int64) string {
	const unit = 1024
	units := []string{"B", "KB", "MB", "GB"}

	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	exp := int(math.Log(float64(bytes)) / math.Log(unit))
	value := float64(bytes) / math.Pow(unit, float64(exp))

	if value >= 100 {
		return fmt.Sprintf("%.0f %s", value, units[exp])
	} else if value >= 10 {
		return fmt.Sprintf("%.1f %s", value, units[exp])
	} else {
		return fmt.Sprintf("%.2f %s", value, units[exp])
	}
}
