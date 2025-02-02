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

import "fmt"

// Colorize applies a foreground color to the given message.
// Parameters:
//   - ansi: A string representing the color in ANSI format.
//   - message: The message to be colorized.
//
// Returns:
//
//	A string with the applied color.
func Colorize(ansi, message string) string {
	return fmt.Sprintf("%s%s%s", ansi, message, Reset)
}
