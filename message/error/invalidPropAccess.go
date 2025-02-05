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

// InvalidPropAccess generates an error message for invalid property access.
// It constructs the message using the logger package to provide error, help, and info sections.
func InvalidPropAccess() string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("Invalid property access"))
	builder.WriteString(
		logger.BuildHelp(
			"The property you are trying to access does not exist in this object.",
			"Make sure you are accessing a module's property and that such exist.",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0014",
			"Full details:",
		),
	)

	return builder.String()
}
