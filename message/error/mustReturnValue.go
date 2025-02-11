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

// MustReturnValue constructs an error message indicating that a function must return a value.
// It uses the logger package to build the error, help, and info messages, and concatenates them into a single string.
func MustReturnValue() string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("This function must return a value"))
	builder.WriteString(
		logger.BuildHelp(
			"This function has not returned any value, but it must.",
			"If you don't want to return a value, use the 'nothing' type.",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0004",
			"Full details:",
		),
	)

	return builder.String()
}
