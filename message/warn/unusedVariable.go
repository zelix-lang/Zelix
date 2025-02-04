/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package warn

import (
	"fluent/logger"
	"strings"
)

// UnusedVariable generates a warning message for an unused variable.
// It provides suggestions and additional information for the user.
//
// Parameters:
// - name: The name of the unused variable.
//
// Returns:
// - A formatted warning message as a string.
func UnusedVariable(name string) string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildWarn("Unused value: '" + name + "'."))
	builder.WriteString(
		logger.BuildHelp(
			"If this is intentional, add an underscore before the variable name.",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0006",
			"Inside of this function:",
		),
	)

	return builder.String()
}
