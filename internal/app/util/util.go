package util

import (
	"fmt"
)

func GetHumanReadableFileSize(fileSizeInBytes int64) string {
	if fileSizeInBytes < 1000 {
		return fmt.Sprintf("%d B", fileSizeInBytes)
	}

	sizeUnits := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	unitIndex := 0
	size := float64(fileSizeInBytes)

	for size >= 1000 && unitIndex < len(sizeUnits)-1 {
		size /= 1000
		unitIndex++
	}

	return fmt.Sprintf("%.2f %s", size, sizeUnits[unitIndex])
}
