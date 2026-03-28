package cfg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"empty string", "", 0},
		{"invalid format", "invalid", 0},
		{"bytes", "100", 100},
		{"bytes with b suffix", "100b", 100},
		{"bytes with B suffix", "100B", 100},
		{"kilobytes lowercase", "1kb", 1024},
		{"kilobytes uppercase", "1KB", 1024},
		{"megabytes lowercase", "1mb", 1024 * 1024},
		{"megabytes uppercase", "1MB", 1024 * 1024},
		{"gigabytes lowercase", "1gb", 1024 * 1024 * 1024},
		{"gigabytes uppercase", "1GB", 1024 * 1024 * 1024},
		{"64kb from config", "64kb", 65536},
		{"1mb from config", "1mb", 1048576},
		{"spaces", " 100 kb ", 102400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseSize(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}
