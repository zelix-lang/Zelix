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

// ShouldNotHaveGenerics generates an error message indicating that the function should not have generics.
// It builds a detailed error message using the logger package and returns it as a string.
func ShouldNotHaveGenerics() string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("This function should not have generics"))
	builder.WriteString(
		logger.BuildHelp(
			"Functions within a module should not have generics.",
			"Define the generics in the module signature if needed.",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0018",
			"Full details:",
		),
	)

	return builder.String()
}
