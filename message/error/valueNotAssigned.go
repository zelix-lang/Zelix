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

// ValueNotAssigned generates an error message indicating that some values are uninitialized.
// It provides a suggestion to define a constructor and initialize the values.
// The details parameter allows additional information to be included in the message.
//
// Returns:
//
//	A string containing the formatted error message.
func ValueNotAssigned() string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("This value is unitialized"))
	builder.WriteString(
		logger.BuildHelp(
			"Define a constructor and use 'this.variable = value'",
			"to initialize this value",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0019",
			"Full details:",
		),
	)

	return builder.String()
}
