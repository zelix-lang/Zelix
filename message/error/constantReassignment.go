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

// ConstantReassignment generates an error message for invalid reassignment to a constant value.
// It provides a detailed error message, help suggestion, and additional information link.
func ConstantReassignment() string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("Invalid reassignment to constant value"))
	builder.WriteString(
		logger.BuildHelp(
			"This value is constant, hence it cannot be reassigned upon being declared.",
			"Try changing the declaration type from 'const' to 'let' when this value is declared.",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0016",
			"Full details:",
		),
	)

	return builder.String()
}
