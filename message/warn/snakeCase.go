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

package warn

import (
	"fluent/logger"
	"strings"
)

// SnakeCase generates a warning message indicating that the provided name should be in snake_case format.
// It also provides a link to more information about the warning.
//
// Parameters:
//   - name: The name that should be in snake_case format.
//
// Returns:
//
//	A string containing the warning message and additional information.
func SnakeCase(name string) string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildWarn("The value name '" + name + "' should be in snake_case."))
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0001",
			"Full details:",
		),
	)

	return builder.String()
}
