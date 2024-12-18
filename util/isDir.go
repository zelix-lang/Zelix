package util

import "os"

// IsDir returns true if the path is a directory
func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}
