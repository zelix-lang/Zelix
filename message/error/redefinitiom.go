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

// Redefinition generates an error message for redefined values.
//
// Parameters:
//   - name: The name of the value that is being redefined.
//
// Returns:
//
//	A string containing the error message.
func Redefinition(name string) string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("Redefinition of value '" + name + "'"))
	builder.WriteString(
		logger.BuildHelp(
			"This value is already defined in the current scope.",
			"Please, change the name of the value or remove the redefinition.",
		),
	)

	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/Book/codes/E0007",
			"Full details:",
		),
	)

	return builder.String()
}
