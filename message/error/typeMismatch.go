/*
   The Fluent Programming Language
   -----------------------------------------------------
   This code is released under the GNU GPL v3 license.
   For more information, please visit:
   https://www.gnu.org/licenses/gpl-3.0.html
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package error

import (
	"fluent/ansi"
	"fluent/logger"
	"fmt"
	"strings"
)

// TypeMismatch generates an error message for type mismatches.
// It takes the expected type and the actual type as arguments and returns a formatted error message.
//
// Parameters:
//   - expected: The expected type as a string.
//   - got: The actual type as a string.
//
// Returns:
//
//	A formatted error message string indicating the type mismatch.
func TypeMismatch(expected string, got string) string {
	// Use a builder to build the string
	var builder strings.Builder

	builder.WriteString(logger.BuildError("Type mismatch"))
	builder.WriteString(
		logger.BuildHelp(
			"The expected type does not match the gotten type.",
			"Modify the type to match the expected one, for example \"2\" instead of 2",
			fmt.Sprintf(
				"Expected %s%s%s, but got %s%s%s",
				ansi.BoldBrightGreen, expected, ansi.Reset,
				ansi.BoldBrightRed, got, ansi.Reset,
			),
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/Book/codes/E0009",
			"Full details:",
		),
	)

	return builder.String()
}
