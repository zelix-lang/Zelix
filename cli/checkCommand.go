package cli

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"surf/analyzer"
	"surf/ansi"
	"surf/ast"
	"surf/lexer"
	"surf/logger"
	"surf/util"
	"time"
)

// CheckCommand represents the check command of the Surf CLI
// it checks a surf file
func CheckCommand(context *cli.Context) (*ast.FileCode, string) {
	// Get the file path
	filePath := context.Args().First()

	showTimer := context.Bool("timer")

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

	start := time.Now()

	// Lex the input
	tokens := lexer.Lex(string(input), filePath)

	if showTimer {
		fmt.Println(
			ansi.Colorize(
				"black_bright",
				"~ Lexed ("+time.Since(start).String()+")",
			),
		)
	}

	start = time.Now()
	// Parse the tokens into a FileCode
	fileCode := ast.Parse(tokens)

	if showTimer {
		fmt.Println(
			ansi.Colorize(
				"black_bright",
				"~ Parsed ("+time.Since(start).String()+")",
			),
		)
	}

	start = time.Now()
	// Analyze the file code
	analyzer.AnalyzeFileCode(fileCode, filePath)

	if showTimer {
		fmt.Println(
			ansi.Colorize(
				"black_bright",
				"~ Checked ("+time.Since(start).String()+")",
			),
		)
	}

	return fileCode, filePath
}
