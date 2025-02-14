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
	"fluent/logger"
	"strings"
)

// InvalidPointer generates an error message for invalid pointer usage.
// It provides detailed information and guidance on the correct usage of pointers.
func InvalidPointer() string {
	builder := strings.Builder{}

	// Add the main error message
	builder.WriteString(logger.BuildError("Invalid pointer"))

	// Add help message with guidance on pointer usage
	builder.WriteString(
		logger.BuildHelp(
			"Pointers can only be used with variables or parameters.",
			"Directly using pointers within a statement is not allowed.",
		),
	)

	// Add additional information and reference link
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0023",
			"Full details:",
		),
	)

	return builder.String()
}
