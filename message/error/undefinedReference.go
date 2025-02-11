/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package error

import (
	"fluent/logger"
	"strings"
)

// UndefinedReference generates an error message for an undefined reference.
// It takes the name of the undefined value as a parameter and returns a formatted error message.
//
// Parameters:
// - name: The name of the undefined value.
//
// Returns:
// - A string containing the formatted error message.
func UndefinedReference(name string) string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("Undefined reference to value '" + name + "'"))
	builder.WriteString(
		logger.BuildHelp(
			"This value is not defined in the current scope.",
			"Make sure you have defined the value before using it.",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0005",
			"Full details:",
		),
	)

	return builder.String()
}
