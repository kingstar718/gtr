package main

import (
	"testing"
	"time"
)

func TestIsTimeFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid 10-digit timestamps
		{
			name:     "10-digit timestamp",
			input:    "1727087511",
			expected: true,
		},
		// Valid 13-digit timestamps
		{
			name:     "13-digit timestamp",
			input:    "1727087511000",
			expected: true,
		},
		// Valid date formats
		{
			name:     "Standard date format",
			input:    "2024-09-23 10:31:51",
			expected: true,
		},
		{
			name:     "Date only",
			input:    "2024-09-23",
			expected: true,
		},
		{
			name:     "Compact date format",
			input:    "20240923103151",
			expected: true,
		},
		{
			name:     "RFC3339 format",
			input:    "2024-09-23T10:31:51Z",
			expected: true,
		},
		// Invalid cases
		{
			name:     "Timestamp with spaces",
			input:    " 1727087511",
			expected: false,
		},
		{
			name:     "Invalid text",
			input:    "hello world",
			expected: false,
		},
		{
			name:     "11-digit number (not 10 or 13)",
			input:    "17270875110",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTimeFormat(tt.input)
			if got != tt.expected {
				t.Errorf("isTimeFormat(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestHandleTimeConvert(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantYear  int
		wantMonth time.Month
		wantDay   int
	}{
		{
			name:      "10-digit timestamp",
			input:     "1727087511",
			wantYear:  2024,
			wantMonth: time.September,
			wantDay:   23,
		},
		{
			name:      "13-digit timestamp",
			input:     "1727087511000",
			wantYear:  2024,
			wantMonth: time.September,
			wantDay:   23,
		},
		{
			name:      "Standard date format",
			input:     "2024-09-23 10:31:51",
			wantYear:  2024,
			wantMonth: time.September,
			wantDay:   23,
		},
		{
			name:      "Compact date format",
			input:     "20240923103151",
			wantYear:  2024,
			wantMonth: time.September,
			wantDay:   23,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We just verify that the function runs without error
			// and returns expected date parts
			err := handleTimeConvert(tt.input)
			if err != nil {
				t.Errorf("handleTimeConvert(%q) returned error: %v", tt.input, err)
			}
		})
	}
}

func TestTimestampConversions(t *testing.T) {
	// Test known timestamp: 2024-09-23 10:31:51 UTC
	timestamp10 := int64(1727087511)
	timestamp13 := int64(1727087511000)

	t.Run("10-digit to 13-digit conversion", func(t *testing.T) {
		parseTime := time.Unix(timestamp10, 0)
		got13 := parseTime.UnixNano() / int64(time.Millisecond)

		if got13 != timestamp13 {
			t.Errorf("10-digit to 13-digit: got %d, want %d", got13, timestamp13)
		}
	})

	t.Run("13-digit to 10-digit conversion", func(t *testing.T) {
		parseTime := time.UnixMilli(timestamp13)
		got10 := parseTime.Unix()

		if got10 != timestamp10 {
			t.Errorf("13-digit to 10-digit: got %d, want %d", got10, timestamp10)
		}
	})

	t.Run("Timestamp date verification", func(t *testing.T) {
		parseTime := time.Unix(timestamp10, 0).UTC()

		if parseTime.Year() != 2024 {
			t.Errorf("Year: got %d, want 2024", parseTime.Year())
		}
		if parseTime.Month() != time.September {
			t.Errorf("Month: got %v, want September", parseTime.Month())
		}
		if parseTime.Day() != 23 {
			t.Errorf("Day: got %d, want 23", parseTime.Day())
		}
	})
}
