package main

import (
	"testing"
)

func TestDetectInputType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected InputType
	}{
		// HTTP URLs
		{
			name:     "HTTPS URL",
			input:    "https://api.example.com",
			expected: TypeHTTP,
		},
		{
			name:     "HTTP URL",
			input:    "http://example.com",
			expected: TypeHTTP,
		},
		// Timestamps
		{
			name:     "10-digit timestamp",
			input:    "1727087511",
			expected: TypeTimestamp,
		},
		{
			name:     "13-digit timestamp",
			input:    "1727087511000",
			expected: TypeTimestamp,
		},
		// Coordinates
		{
			name:     "Coordinate with comma",
			input:    "113.901495,22.499501",
			expected: TypeCoordinate,
		},
		{
			name:     "Coordinate with pipe",
			input:    "113.901495|22.499501",
			expected: TypeCoordinate,
		},
		{
			name:     "Coordinate with space",
			input:    "113.901495 22.499501",
			expected: TypeCoordinate,
		},
		{
			name:     "Negative coordinates",
			input:    "-113.901495,-22.499501",
			expected: TypeCoordinate,
		},
		// Text (fallback)
		{
			name:     "Plain text",
			input:    "hello world",
			expected: TypeText,
		},
		{
			name:     "Base64 string",
			input:    "aGVsbG8gd29ybGQ=",
			expected: TypeText,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectInputType(tt.input)
			if got != tt.expected {
				t.Errorf("detectInputType(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
