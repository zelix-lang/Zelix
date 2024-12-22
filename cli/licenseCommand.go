package cli

import (
	"github.com/urfave/cli/v2"
	"zyro/ansi"
)

// printFullLicense prints the full license of the Zyro CLI
func printFullLicense() {
	println(
		"You may check the full license at:",
		ansi.Colorize("underline", "https://www.gnu.org/licenses/gpl-3.0.en.html"),
	)
}

// LicenseCommand represents the license command of the Zyro CLI
// it prints the license
func LicenseCommand(context *cli.Context) {
	ShowHeaderMessage()

	if context.Bool("full") {
		printFullLicense()
		return
	}

	println("Copyright (C) 2024 Rodrigo R. & All Zyro Contributors")
	println("This program comes with ABSOLUTELY NO WARRANTY; for details type `zyro/ license`.")
	println("This is free software, and you are welcome to redistribute it under certain conditions;")
	println("type `zyro/ license --full` for details.")
}
