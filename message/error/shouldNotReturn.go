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

// ShouldNotReturn generates an error message indicating that a function should not return a value.
// It provides detailed information and guidance on how to resolve the issue.
// Returns a string containing the error message.
func ShouldNotReturn() string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("This function should not return a value"))
	builder.WriteString(
		logger.BuildHelp(
			"Functions with the 'nothing' return type should not return a value.",
			"Remove the return statement from the function or modify the return type.",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0012",
			"Full details:",
		),
	)

	return builder.String()
}
