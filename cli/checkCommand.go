package cli

import (
	"github.com/urfave/cli/v2"
	"os"
	"surf/ansi"
	"surf/logger"
	"surf/util"
)

// CheckCommand represents the check command of the Surf CLI
// it checks a surf file
func CheckCommand(context *cli.Context) (string, string) {
	// Get the file path
	filePath := context.Args().First()

	if len(filePath) == 0 {
		logger.Error("Empty file path")
		logger.Help(
			"Provide a file path after the run command",
			"For example: "+ansi.Colorize("green_bright_bold", "surf r file.surf"),
		)

		os.Exit(1)
	}

	// Check if the file exists
	if !util.FileExists(filePath) {
		logger.Error("File does not exist")
		logger.Help(
			"Make sure the file exists",
			"Check the file path",
		)

		os.Exit(1)
	}

	// Read the file
	input, err := os.ReadFile(filePath)

	if err != nil {
		logger.Error("Error reading file")
		logger.Help(
			"Make sure the file is readable",
			"Check the file path",
		)

		os.Exit(1)
	}

	return string(input), filePath
}
