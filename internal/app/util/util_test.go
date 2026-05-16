package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHumanReadableFileSize(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{"Bytes", 500, "500 B"},
		{"Kilobytes", 1024, "1.02 KB"},
		{"Megabytes", 1048576, "1.05 MB"},
		{"Gigabytes", 1073741824, "1.07 GB"},
		{"Terabytes", 1099511627776, "1.10 TB"},
		{"Zero", 0, "0 B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, GetHumanReadableFileSize(tt.input))
		})
	}
}
