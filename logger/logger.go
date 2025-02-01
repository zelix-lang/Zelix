/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package logger

import (
	"fluent/ansi"
	"fmt"
	"github.com/muesli/termenv"
)

// Prefixes for the different log levels
var errorPrefix = ansi.Colorize(ansi.BrightRed, "[ERROR] ").Bold()
var infoPrefix = ansi.Colorize(ansi.BrightBlue, "[INFO]  ").Bold()
var warnPrefix = ansi.Colorize(ansi.BrightYellow, "[WARN] ").Bold()
var helpPrefix = ansi.Colorize(ansi.BrightGreen, "[HELP]  ").Bold()

// printMessageImpl prints each message with the given prefix.
// Parameters:
//   - prefix: the style to be applied to the prefix.
//   - message: variadic string arguments to be printed.
//   - color: the color to be applied to the message.
//   - colorize: whether to colorize the message or not.
func printMessageImpl(prefix termenv.Style, color string, colorize bool, message ...string) {
	for _, m := range message {
		fmt.Print(prefix)
		if colorize {
			fmt.Print(ansi.Colorize(color, m))
		} else {
			fmt.Print(m)
		}

		// Print a new line
		fmt.Println()
	}
}

// Warn prints the provided messages with a warning prefix.
// Parameters:
//   - message: variadic string arguments to be logged.
func Warn(message ...string) {
	printMessageImpl(warnPrefix, "", false, message...)
}

// Info prints the provided messages with an info prefix.
// Parameters:
//   - message: variadic string arguments to be logged.
func Info(message ...string) {
	printMessageImpl(infoPrefix, ansi.BrightBlack, true, message...)
}

// Help prints the provided messages with a help prefix.
// Parameters:
//   - message: variadic string arguments to be logged.
func Help(message ...string) {
	printMessageImpl(helpPrefix, "", false, message...)
}

// Error prints the provided messages with an error prefix.
// Parameters:
//   - message: variadic string arguments to be logged.
func Error(message ...string) {
	printMessageImpl(errorPrefix, ansi.BrightRed, true, message...)
}
