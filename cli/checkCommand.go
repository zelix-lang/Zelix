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
	"fluent/pkg"
	"fluent/util"
	"github.com/urfave/cli/v3"
	"os"
	"path"
	"path/filepath"
)

// CheckCommand performs a series of checks on the provided path.
// It retrieves the path from the context, validates it, converts it to file codes,
// and analyzes the project's codebase. It returns the sorted file codes and the original file codes.
//
// Parameters:
//
//	context (*cli.Command): The CLI command context containing the arguments.
//
// Returns:
//   - A slice of sorted file codes
//   - A map of original file codes.
//   - The original path where the project is located.
func CheckCommand(context *cli.Command) ([]filecode.FileCode, map[string]filecode.FileCode, string) {
	ShowHeaderMessage()

	// Retrieve the path from the context
	target := context.Args().First()

	// Check if the path exists
	if target == "" {
		logger.Error("No path provided")
		logger.Info("Usage: fluent <check|c> <path>")
		os.Exit(1)
	}

	// Convert the path to an absolute path
	target, absError := filepath.Abs(target)

	if absError != nil {
		logger.Error("Could not convert the path to an absolute path")
		logger.Help("Validate the path and try again")
		os.Exit(1)
	}

	// Make sure the path exists and is a directory
	if !util.DirExists(target) {
		logger.Error("The provided path does not exist or is not a directory")
		logger.Help("Validate the path and try again")
		os.Exit(1)
	}

	// Get the package.fluent file
	packagePath := filepath.Join(target, "package.fluent")
	if !util.FileExists(packagePath) {
		logger.Error("The package.fluent file does not exist in the provided path")
		logger.Help("Make sure the file exists and try again")
		os.Exit(1)
	}

	// Parse the package file
	metadata := pkg.ParsePackage(packagePath)

	// Join the path with the entry file
	target = path.Join(target, metadata.Entry)

	// Convert the code to file codes
	fileCodes := converter.ConvertToFileCode(target, false)

	// Analyze the project's codebase
	sortedFileCodes := analyzer.AnalyzeCode(fileCodes, target, false)

	// The build command depends on the check command
	// hence, it also needs the file codes
	return sortedFileCodes, fileCodes, target
}
