package cli

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"surf/ansi"
	"surf/ast"
	"surf/checker"
	"surf/core"
	"surf/lexer"
	"time"
)

// RunCommand represents the run command of the Surf CLI
// it runs a surf file
func RunCommand(context *cli.Context) {
	showTimer := context.Bool("timer")
	input, filePath := CheckCommand(context)

	start := time.Now()

	// Lex the input
	tokens := lexer.Lex(input, filePath)

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
	checker.AnalyzeFileCode(*fileCode, filePath)

	if showTimer {
		fmt.Println(
			ansi.Colorize(
				"black_bright",
				"~ Checked ("+time.Since(start).String()+")",
			),
		)
	}

	start = time.Now()
	// Interpret the file code
	core.Interpret(fileCode, filePath)

	if showTimer {
		fmt.Println(
			ansi.Colorize(
				"black_bright",
				"~ Interpreted ("+time.Since(start).String()+")",
			),
		)
	}
}
