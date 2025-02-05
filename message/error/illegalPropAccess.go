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

func IllegalPropAccess() string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("Illegal property access"))
	builder.WriteString(
		logger.BuildHelp(
			"The current scope does not have permission to access this property.",
			"Properties are private by default.",
			"Check if the property has any setter or getter methods.",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0015",
			"Full details:",
		),
	)

	return builder.String()
}
