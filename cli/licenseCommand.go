/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package cli

import (
	"fmt"
	"github.com/urfave/cli/v3"
)

// printFullLicense prints the full license of the Fluent CLI
func printFullLicense() {
	fmt.Println(
		"You may check the full license at:",
		"https://www.gnu.org/licenses/gpl-3.0.en.html",
	)
}

// LicenseCommand represents the license command of the Fluent CLI
// it prints the license
func LicenseCommand(context *cli.Command) {
	ShowHeaderMessage()

	if context.Bool("full") {
		printFullLicense()
		return
	}

	println("Copyright (C) 2024 Rodrigo R. & All Fluent Contributors")
	println("This program comes with ABSOLUTELY NO WARRANTY; for details type `fluent license`.")
	println("This is free software, and you are welcome to redistribute it under certain conditions;")
	println("type `fluent license --full` for details.")
}
