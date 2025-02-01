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
	"fluent/ansi"
	"fmt"
)

// ShowHeaderMessage shows the header message of the Fluent CLI
func ShowHeaderMessage() {
	fmt.Println(ansi.Colorize(ansi.BrightBlue, "Fluent"))
	fmt.Println(ansi.Colorize(ansi.BrightBlack, "A blazingly fast programming language"))
	fmt.Println()
}
