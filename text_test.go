package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/url"
	"testing"
)

func TestBase64Operations(t *testing.T) {
	tests := []struct {
		name            string
		plaintext       string
		expectedEncoded string
	}{
		{
			name:            "Basic string",
			plaintext:       "hello world",
			expectedEncoded: "aGVsbG8gd29ybGQ=",
		},
		{
			name:            "Empty string",
			plaintext:       "",
			expectedEncoded: "",
		},
		{
			name:            "Special characters",
			plaintext:       "hello@world!",
			expectedEncoded: "aGVsbG9Ad29ybGQh",
		},
		{
			name:            "Numbers",
			plaintext:       "12345",
			expectedEncoded: "MTIzNDU=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+" encode", func(t *testing.T) {
			encoded := base64.StdEncoding.EncodeToString([]byte(tt.plaintext))
			if encoded != tt.expectedEncoded {
				t.Errorf("Base64 encode(%q) = %q, want %q", tt.plaintext, encoded, tt.expectedEncoded)
			}
		})

		t.Run(tt.name+" decode", func(t *testing.T) {
			decoded, err := base64.StdEncoding.DecodeString(tt.expectedEncoded)
			if err != nil {
				t.Errorf("Base64 decode(%q) returned error: %v", tt.expectedEncoded, err)
			}
			if string(decoded) != tt.plaintext {
				t.Errorf("Base64 decode(%q) = %q, want %q", tt.expectedEncoded, string(decoded), tt.plaintext)
			}
		})

		t.Run(tt.name+" roundtrip", func(t *testing.T) {
			encoded := base64.StdEncoding.EncodeToString([]byte(tt.plaintext))
			decoded, _ := base64.StdEncoding.DecodeString(encoded)
			if string(decoded) != tt.plaintext {
				t.Errorf("Base64 roundtrip failed: %q → %q → %q", tt.plaintext, encoded, string(decoded))
			}
		})
	}

	t.Run("Invalid base64", func(t *testing.T) {
		_, err := base64.StdEncoding.DecodeString("!!!invalid!!!")
		if err == nil {
			t.Errorf("Base64 decode should fail for invalid data")
		}
	})
}

func TestURLOperations(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedEnc string
		expectedDec string
	}{
		{
			name:        "Plain text",
			input:       "hello world",
			expectedEnc: "hello+world",
			expectedDec: "hello world",
		},
		{
			name:        "Special characters",
			input:       "hello@example.com",
			expectedEnc: "hello%40example.com",
			expectedDec: "hello@example.com",
		},
		{
			name:        "URL path",
			input:       "path/to/file",
			expectedEnc: "path%2Fto%2Ffile",
			expectedDec: "path/to/file",
		},
		{
			name:        "Empty string",
			input:       "",
			expectedEnc: "",
			expectedDec: "",
		},
		{
			name:        "Already encoded",
			input:       "hello%20world",
			expectedEnc: "hello%2520world",
			expectedDec: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+" encode", func(t *testing.T) {
			encoded := url.QueryEscape(tt.input)
			if encoded != tt.expectedEnc {
				t.Errorf("URL encode(%q) = %q, want %q", tt.input, encoded, tt.expectedEnc)
			}
		})

		t.Run(tt.name+" decode", func(t *testing.T) {
			decoded, err := url.QueryUnescape(tt.input)
			if err != nil {
				t.Errorf("URL decode(%q) returned error: %v", tt.input, err)
			}
			if decoded != tt.expectedDec {
				t.Errorf("URL decode(%q) = %q, want %q", tt.input, decoded, tt.expectedDec)
			}
		})
	}
}

func TestMD5Hash(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "hello",
			input:    "hello",
			expected: "5d41402abc4b2a76b9719d911017c592",
		},
		{
			name:     "hello world",
			input:    "hello world",
			expected: "5eb63bbbe01eeed093cb22bb8f5acdc3",
		},
		{
			name:     "password",
			input:    "password",
			expected: "5f4dcc3b5aa765d61d8327deb882cf99",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name:     "numbers",
			input:    "12345",
			expected: "827ccb0eea8a706c4c34a16891f84e7b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := md5.Sum([]byte(tt.input))
			got := fmt.Sprintf("%x", hash)
			if got != tt.expected {
				t.Errorf("MD5(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
