package cli

import (
	"github.com/urfave/cli/v2"
)

// RunCommand represents the run command of the Fluent CLI
// it runs a fluent file
func RunCommand(context *cli.Context) {
	//context.Bool("timer")
	CheckCommand(context)

	// TODO!
}
