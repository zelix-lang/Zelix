package util

import (
	"path/filepath"
	"strings"
)

// The system's path separator
var separator = string(filepath.Separator)

// DirName returns the directory name of the path
// Example: DirName("/path/to/file.txt") -> "to"
func DirName(path string) string {
	dir := filepath.Dir(path)

	// Handle cases with no directory or root directory
	if dir == "." || dir == separator {
		return ""
	}

	parts := strings.Split(dir, separator)
	if parts[0] == "" {
		parts = parts[1:]
	}

	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""

}
