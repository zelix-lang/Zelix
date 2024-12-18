package cli

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"surf/ansi"
	"surf/core"
	"time"
)

// RunCommand represents the run command of the Surf CLI
// it runs a surf file
func RunCommand(context *cli.Context) {
	showTimer := context.Bool("timer")
	fileCode, filePath := CheckCommand(context)

	start := time.Now()
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
