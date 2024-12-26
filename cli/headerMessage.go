package cli

import "fluent/ansi"

// ShowHeaderMessage shows the header message of the Fluent CLI
func ShowHeaderMessage() {
	println(ansi.Colorize("blue_bright", "Fluent"))
	println(ansi.Colorize("black_bright", "A blazingly fast programming language"))
	println()
}
