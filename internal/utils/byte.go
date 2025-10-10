package utils

import "fmt"

func FormatBytes(n int64) string {
	if n < 0 {
		return "-"
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	v := float64(n)
	i := 0
	for v >= 1024 && i < len(units)-1 {
		v /= 1024
		i++
	}

	return fmt.Sprintf("%.2f%s", v, units[i])
}
