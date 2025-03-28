/*
   The Fluent Programming Language
   -----------------------------------------------------
   This code is released under the GNU GPL v3 license.
   For more information, please visit:
   https://www.gnu.org/licenses/gpl-3.0.html
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package util

import (
	"os"
	"strings"
)

var cwd, _ = os.Getwd()

// DiscardCwd removes the current working directory (cwd) prefix from the given path.
// If the cwd is not a prefix of the path, the original path is returned.
// If the path starts with a slash after removing the cwd, the slash is also removed.
// Parameters:
// - path: The file path from which to discard the cwd.
// Returns: The modified path with the cwd and leading slash removed, if applicable.
func DiscardCwd(path string) string {
	// Check if the cwd is a prefix of the path
	if strings.HasPrefix(path, cwd) {
		// Remove the cwd from the path
		path = (path)[len(cwd):]
	} else {
		// Return the path as is
		return path
	}

	// Check if the path starts with a slash
	if strings.HasPrefix(path, "/") {
		// Remove the slash
		path = (path)[1:]
	}

	return path
}

// DirExists checks if a directory exists at the given path.
// Parameters:
// - path: The file path to check.
// Returns: True if the directory exists, false otherwise.
func DirExists(path string) bool {
	info, err := os.Stat(path)

	if os.IsNotExist(err) || info == nil {
		return false
	}

	return info.IsDir()
}

// FileExists checks if a file exists at the given path.
// Parameters:
// - path: The file path to check.
// Returns: True if the file exists and is not a directory, false otherwise.
func FileExists(path string) bool {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

// GetDir returns the directory part of the given file path.
// Parameters:
// - path: The file path from which to extract the directory.
// Returns: The directory part of the path, or an empty string if no directory is found.
func GetDir(path string) string {
	// Get the last index of the separator
	lastIndex := strings.LastIndex(path, string(os.PathSeparator))

	// Check if the separator was found
	if lastIndex == -1 {
		return ""
	}

	return path[:lastIndex]
}

// FileName returns the file name from the given file path.
// Parameters:
// - path: The file path from which to extract the file name.
// Returns: The file name part of the path, or the entire path if no separator is found.
func FileName(path *string) string {
	// Get the last index of the separator
	lastIndex := strings.LastIndex(*path, string(os.PathSeparator))

	// Check if the separator was found
	if lastIndex == -1 {
		return *path
	}

	return (*path)[lastIndex+1:]
}

// ReadDir reads the contents of the directory specified by the given path.
// Parameters:
// - path: The file path of the directory to read.
// Returns: A slice of os.FileInfo containing the directory's contents, or an empty slice if the directory does not exist or an error occurs.
func ReadDir(path string) []os.FileInfo {
	if !DirExists(path) {
		return make([]os.FileInfo, 0)
	}

	// Read the directory
	dir, err := os.Open(path)
	if err != nil {
		return make([]os.FileInfo, 0)
	}

	// Close the directory
	defer func(dir *os.File) {
		err := dir.Close()
		if err != nil {
			panic(err)
		}
	}(dir)

	// Read the directory
	files, err := dir.Readdir(0)

	if err != nil {
		return make([]os.FileInfo, 0)
	}

	return files
}
