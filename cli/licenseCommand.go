package cli

import (
	"fluent/ansi"
	"github.com/urfave/cli/v2"
)

// printFullLicense prints the full license of the Fluent CLI
func printFullLicense() {
	println(
		"You may check the full license at:",
		ansi.Colorize("underline", "https://www.gnu.org/licenses/gpl-3.0.en.html"),
	)
}

// LicenseCommand represents the license command of the Fluent CLI
// it prints the license
func LicenseCommand(context *cli.Context) {
	ShowHeaderMessage()

	if context.Bool("full") {
		printFullLicense()
		return
	}

	println("Copyright (C) 2024 Rodrigo R. & All Fluent Contributors")
	println("This program comes with ABSOLUTELY NO WARRANTY; for details type `fluent/ license`.")
	println("This is free software, and you are welcome to redistribute it under certain conditions;")
	println("type `fluent/ license --full` for details.")
}
