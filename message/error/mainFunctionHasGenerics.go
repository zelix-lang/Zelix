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

// MainFunctionHasGenerics generates an error message indicating that the main function cannot have generics.
// It builds a detailed error message using the logger package and returns it as a string.
func MainFunctionHasGenerics() string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("The main function cannot have generics"))
	builder.WriteString(
		logger.BuildHelp(
			"The main function cannot be built dynamically",
			"with generic types.",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/Book/codes/E0026",
			"Full details:",
		),
	)

	return builder.String()
}
