// Package input provides file and string input handling for comment detection.
package input

import (
	"os"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

// ReadFile reads file content with UTF-8 to Latin-1 fallback.
//
// Returns:
//   - File content as string on success
//   - Empty string if file not found (graceful handling)
//
// Notes:
//   - UTF-8 decode failure falls back to latin-1
//   - FileNotFoundError returns empty string (graceful handling)
func ReadFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		// FileNotFoundError or other errors -> return empty string
		return ""
	}

	// Check if valid UTF-8
	if utf8.Valid(data) {
		return string(data)
	}

	// Fallback to Latin-1 (ISO-8859-1) decoding
	decoded, err := charmap.ISO8859_1.NewDecoder().Bytes(data)
	if err != nil {
		// If Latin-1 decoding also fails, return raw bytes as string
		return string(data)
	}

	return string(decoded)
}

// ReadString passes through string content unchanged.
// This function exists for API consistency with ReadFile.
func ReadString(content string) string {
	return content
}
