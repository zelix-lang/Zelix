/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package util

import (
	"fluent/logger"
	"os"
)

// ReadFile reads the contents of a file at the given path.
// Parameters:
//   - path: The path to the file to be read.
//
// Returns:
//   - The contents of the file as a string.
func ReadFile(path string) string {
	// Read the file
	file, err := os.ReadFile(path)

	if err != nil {
		logger.Error("Could not read the file at:", path)
		logger.Help("Make sure the provided path points to a file")
		logger.Help("Make sure the provided path exists")
		os.Exit(1)
	}

	// Convert the contents to a string
	contents := string(file)
	return contents
}
