/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package ansi

import (
	"github.com/muesli/termenv"
	"os"
)

// Create a new output for colorizing the terminal
var output = termenv.NewOutput(os.Stdout)

// Colorize applies a foreground color to the given message.
// Parameters:
//   - hex: A string representing the color in hexadecimal format.
//   - message: The message to be colorized.
//
// Returns:
//
//	A Style object with the applied color.
func Colorize(hex, message string) termenv.Style {
	coloredMessage := output.String(message)
	coloredMessage = coloredMessage.Foreground(output.Color(hex))
	return coloredMessage
}

// Underline applies an underline style to the given message.
// Parameters:
//   - message: The message to be underlined.
//
// Returns:
//
//	A Style object with the applied underline.
func Underline(message string) termenv.Style {
	return output.String(message).Underline()
}
