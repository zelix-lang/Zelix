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

// CircularModuleDependency generates an error message for a circular module dependency.
// It takes a string representing the dependency chain and returns a formatted error message.
//
// Parameters:
//   - chain: A string representing the full dependency chain.
//
// Returns:
//   - A formatted error message string indicating the circular dependency.
func CircularModuleDependency(chain string) string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("This module depends on its own"))
	builder.WriteString(
		logger.BuildHelp(
			"This module has a property that is either the same module or",
			"another module that depends on this one.",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/Book/codes/E0020",
			"Full dependency chain:",
		),
	)

	builder.WriteString(chain)
	builder.WriteString(
		logger.BuildInfo(
			"Full details:",
		),
	)

	return builder.String()
}
