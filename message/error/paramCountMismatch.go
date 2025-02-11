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

// ParamCountMismatch generates an error message for a mismatched count of parameters.
//
// Parameters:
//   - count: The expected number of parameters as a string.
//
// Returns:
//
//	A formatted error message string.
func ParamCountMismatch(count string) string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("Mismatched count of parameters"))
	builder.WriteString(
		logger.BuildHelp(
			"This function expected "+count+" parameters.",
			"Make sure you have passed the correct number of parameters.",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0010",
			"Full details:",
		),
	)

	return builder.String()
}
