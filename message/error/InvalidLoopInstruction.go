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

// InvalidLoopInstruction generates an error message for invalid loop instructions.
// It returns a string containing the error message, help information, and a reference link.
func InvalidLoopInstruction() string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("Invalid instruction"))
	builder.WriteString(
		logger.BuildHelp(
			"This instruction is not inside a loop block.",
			"Use break/continue instructions only inside loop blocks.",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/Book/codes/E0022",
			"Full details:",
		),
	)

	return builder.String()
}
