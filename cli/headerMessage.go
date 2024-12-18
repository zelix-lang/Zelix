package cli

import "surf/ansi"

// ShowHeaderMessage shows the header message of the Surf CLI
func ShowHeaderMessage() {
	println(ansi.Colorize("blue_bright", "Surf"))
	println(ansi.Colorize("black_bright", "A blazingly fast programming language"))
	println()
}
