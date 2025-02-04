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

// NothingParam generates an error message indicating that parameters cannot have type 'nothing'.
//
// Returns:
//
//	string: The constructed error message.
func NothingParam() string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("Parameters cannot have type 'nothing'"))
	builder.WriteString(
		logger.BuildHelp(
			"'nothing' is used for return types",
			"A parameter that holds 'nothing' is not useful in a function",
			"Change the type of the parameter to a valid type",
		),
	)

	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0002",
			"Full details:",
		),
	)

	// Return the error message
	return builder.String()
}
