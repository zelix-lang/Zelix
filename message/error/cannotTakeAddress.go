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

// CannotTakeAddress returns a detailed error message explaining why the address of a value cannot be taken.
// It uses the logger package to build the error, help, and info messages.
func CannotTakeAddress() string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("Cannot take the address of this value"))
	builder.WriteString(
		logger.BuildHelp(
			"Pointers work for variables and function calls.",
			"They allow another function to modify local data that can later be used.",
			"This value does not hold useful information that can be re-utilized.",
			"You can define a local variable that hold this value and pass a pointer to it.",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0013",
			"Full details:",
		),
	)

	return builder.String()
}
