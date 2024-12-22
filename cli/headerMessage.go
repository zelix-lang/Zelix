package cli

import "zyro/ansi"

// ShowHeaderMessage shows the header message of the Zyro CLI
func ShowHeaderMessage() {
	println(ansi.Colorize("blue_bright", "Zyro"))
	println(ansi.Colorize("black_bright", "A blazingly fast programming language"))
	println()
}
