package cli

import (
	"fluent/ir"
	"fmt"
	"github.com/urfave/cli/v2"
)

// BuildCommand compiles the given Fluent project into an executable
func BuildCommand(context *cli.Context) {
	code, _ := CheckCommand(context)

	// Build the IR
	fluentIr := ir.EmitIR(*code)

	fmt.Println(fluentIr)
}
