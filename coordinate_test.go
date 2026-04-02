package main

import (
	"math"
	"testing"
)

func TestCoordinateStringToFloat(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expectLng  float64
		expectLat  float64
		expectFail bool
	}{
		{
			name:      "Comma separator",
			input:     "113.901495,22.499501",
			expectLng: 113.901495,
			expectLat: 22.499501,
		},
		{
			name:      "Pipe separator",
			input:     "113.901495|22.499501",
			expectLng: 113.901495,
			expectLat: 22.499501,
		},
		{
			name:      "Negative coordinates",
			input:     "-113.901495,-22.499501",
			expectLng: -113.901495,
			expectLat: -22.499501,
		},
		{
			name:      "Integer coordinates",
			input:     "113,22",
			expectLng: 113,
			expectLat: 22,
		},
		{
			name:       "Invalid - invalid longitude",
			input:      "abc,22",
			expectFail: true,
		},
		{
			name:       "Invalid - no separator",
			input:      "113.901495",
			expectFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := coordinateStringToFloat(tt.input)
			if tt.expectFail {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error for input %q: %v", tt.input, err)
			}
			if len(result) != 2 {
				t.Errorf("Expected 2 values, got %d", len(result))
			}
			if math.Abs(result[0]-tt.expectLng) > 0.00001 {
				t.Errorf("Longitude: got %v, want %v", result[0], tt.expectLng)
			}
			if math.Abs(result[1]-tt.expectLat) > 0.00001 {
				t.Errorf("Latitude: got %v, want %v", result[1], tt.expectLat)
			}
		})
	}
}

func TestIsCoordinateFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid comma format",
			input:    "113.901495,22.499501",
			expected: true,
		},
		{
			name:     "Valid pipe format",
			input:    "113.901495|22.499501",
			expected: true,
		},
		{
			name:     "Valid space format",
			input:    "113.901495 22.499501",
			expected: true,
		},
		{
			name:     "Negative coordinates",
			input:    "-113.901495,-22.499501",
			expected: true,
		},
		{
			name:     "Invalid - missing longitude",
			input:    "22.499501",
			expected: false,
		},
		{
			name:     "Invalid - text",
			input:    "hello world",
			expected: false,
		},
		{
			name:     "Invalid - trailing space",
			input:    "113.901495, 22.499501 ",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCoordinateFormat(tt.input)
			if got != tt.expected {
				t.Errorf("isCoordinateFormat(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGPSUtilConversions(t *testing.T) {
	gps := &GPSUtil{}

	// Test data: WGS84 coordinates for Beijing
	wgs84Lat := 39.9042
	wgs84Lng := 116.4074

	t.Run("WGS84 to GCJ02 and back", func(t *testing.T) {
		gcj02 := gps.WGS84_To_Gcj02(wgs84Lat, wgs84Lng)
		backToWgs84 := gps.GCJ02_To_WGS84(gcj02[0], gcj02[1])

		// Check roundtrip accuracy (should be within 0.0001 degrees)
		if math.Abs(backToWgs84[0]-wgs84Lat) > 0.0001 {
			t.Errorf("WGS84 to GCJ02 roundtrip failed: lat %v → %v → %v",
				wgs84Lat, gcj02[0], backToWgs84[0])
		}
		if math.Abs(backToWgs84[1]-wgs84Lng) > 0.0001 {
			t.Errorf("WGS84 to GCJ02 roundtrip failed: lng %v → %v → %v",
				wgs84Lng, gcj02[1], backToWgs84[1])
		}
	})

	t.Run("GCJ02 to BD09 and back", func(t *testing.T) {
		gcj02Lat := 39.9042
		gcj02Lng := 116.4074
		bd09 := gps.gcj02_To_Bd09(gcj02Lat, gcj02Lng)
		backToGcj02 := gps.bd09_To_Gcj02(bd09[0], bd09[1])

		// Check roundtrip accuracy
		if math.Abs(backToGcj02[0]-gcj02Lat) > 0.0001 {
			t.Errorf("GCJ02 to BD09 roundtrip failed: lat %v → %v → %v",
				gcj02Lat, bd09[0], backToGcj02[0])
		}
		if math.Abs(backToGcj02[1]-gcj02Lng) > 0.0001 {
			t.Errorf("GCJ02 to BD09 roundtrip failed: lng %v → %v → %v",
				gcj02Lng, bd09[1], backToGcj02[1])
		}
	})

	t.Run("OutOfChina coordinates pass-through", func(t *testing.T) {
		// New York coordinates
		newYorkLat := 40.7128
		newYorkLng := -74.0060

		result := gps.WGS84_To_Gcj02(newYorkLat, newYorkLng)

		// Should return unchanged
		if math.Abs(result[0]-newYorkLat) > 0.00001 {
			t.Errorf("Out of China coords should not change: lat %v → %v", newYorkLat, result[0])
		}
		if math.Abs(result[1]-newYorkLng) > 0.00001 {
			t.Errorf("Out of China coords should not change: lng %v → %v", newYorkLng, result[1])
		}
	})

	t.Run("Retain6 decimal precision", func(t *testing.T) {
		num := 113.9014956789
		result := gps.retain6(num)
		expected := 113.901496

		if math.Abs(result-expected) > 0.0000001 {
			t.Errorf("retain6(%v) = %v, want %v", num, result, expected)
		}
	})
}
