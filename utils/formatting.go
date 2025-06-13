package utils

import (
	"fmt"
	"math"
)

// FormatBytes converts a byte count into a human-readable string with appropriate suffixes.
func FormatBytes(bytes int64) string {
	if bytes == 0 {
		return "0B"
	}

	if bytes < 1024 {
		return fmt.Sprintf("%dB", bytes)
	}

	suffixes := []string{"K", "M", "G", "T", "P"} // Add more if needed: E, Z, Y
	val := float64(bytes)
	i := 0

	for val >= 1024 && i < len(suffixes) {
		val /= 1024
		i++
	}

	// Round to one decimal place
	// math.Pow10(1) is 10. Multiplying by 10, rounding, then dividing by 10 achieves one decimal place.
	roundedVal := math.Round(val*10) / 10

	return fmt.Sprintf("%.1f%s", roundedVal, suffixes[i-1])
}
