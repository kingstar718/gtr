package main

import (
	"testing"
)

func TestIsHTTPURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "HTTPS URL",
			input:    "https://api.example.com",
			expected: true,
		},
		{
			name:     "HTTP URL",
			input:    "http://example.com",
			expected: true,
		},
		{
			name:     "HTTPS with path",
			input:    "https://api.example.com/v1/users",
			expected: true,
		},
		{
			name:     "HTTP with query",
			input:    "http://example.com?foo=bar",
			expected: true,
		},
		{
			name:     "Non-HTTP URL",
			input:    "ftp://example.com",
			expected: false,
		},
		{
			name:     "Plain text",
			input:    "hello world",
			expected: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "File path",
			input:    "/path/to/file",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isHTTPURL(tt.input)
			if got != tt.expected {
				t.Errorf("isHTTPURL(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
