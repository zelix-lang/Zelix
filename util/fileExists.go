package util

import "os"

// FileExists returns true if the file exists
func FileExists(file string) bool {
	_, err := os.Stat(file)

	return !os.IsNotExist(err)
}
