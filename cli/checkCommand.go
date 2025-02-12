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

package cli

import (
	"fluent/analyzer"
	"fluent/filecode"
	"fluent/filecode/converter"
	"fluent/logger"
	"github.com/urfave/cli/v3"
	"os"
	"path/filepath"
)

// CheckCommand is a CLI command that checks the provided path for code files,
// converts them to file codes, and analyzes the project's codebase.
// It returns a map of file codes.
// Parameters:
// - context: the CLI context containing the command arguments.
func CheckCommand(context *cli.Command) []filecode.FileCode {
	ShowHeaderMessage()

	// Retrieve the path from the context
	path := context.Args().First()

	// Check if the path exists
	if path == "" {
		logger.Error("No path provided")
		logger.Info("Usage: fluent <check|c> <path>")
		os.Exit(1)
	}

	// Convert the path to an absolute path
	path, absError := filepath.Abs(path)

	if absError != nil {
		logger.Error("Could not convert the path to an absolute path")
		logger.Help("Validate the path and try again")
		os.Exit(1)
	}

	// Convert the code to file codes
	fileCodes := converter.ConvertToFileCode(path, false)

	// Analyze the project's codebase
	sortedFileCodes := analyzer.AnalyzeCode(fileCodes, path, false)

	// The build command depends on the check command
	// hence, it also needs the file codes
	return sortedFileCodes
}
