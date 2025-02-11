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

// InvalidDereference constructs an error message for invalid dereference attempts.
// It returns a string containing the error message, help information, and a reference link.
func InvalidDereference() string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("Dereference of a non-pointer value"))
	builder.WriteString(
		logger.BuildHelp(
			"Dereferencing works for converting pointers to values.",
			"A value that is not a pointer cannot be dereferenced.",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0008",
			"Full details:",
		),
	)

	return builder.String()
}
