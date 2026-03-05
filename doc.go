// Package main implements gtr - a Go terminal tools for quick conversions and transformations.
//
// gtr provides a command-line interface for performing various data transformations:
//
//   - Coordinate conversion: Convert between WGS84, GCJ02 (Chinese coordinate system), and BD09 (Baidu coordinate system)
//   - HTTP requests: Send GET/POST requests and view detailed response information
//   - Time conversion: Convert between Unix timestamps (10-digit and 13-digit) and human-readable date formats
//   - Text encoding/decoding: Base64, URL encoding/decoding, and MD5 hashing
//
// Usage:
//
//	gtr coordinate <longitude> <latitude>     # Coordinate system conversion
//	gtr http <URL> [<JSON_DATA>]              # HTTP GET/POST requests
//	gtr time <timestamp|date>                 # Time format conversion
//	gtr text <operation> <subop> <text>       # Text encoding/decoding
//
// For more information, visit: https://github.com/kingstar718/gtr
package main
