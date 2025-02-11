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

// DoesNotHaveConstructor generates an error message indicating that a module does not have a constructor.
// It provides suggestions on how to resolve the issue and includes a link to further information.
//
// Returns:
//
//	A string containing the error message.
func DoesNotHaveConstructor() string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("This module does not have a constructor"))
	builder.WriteString(
		logger.BuildHelp(
			"There is no constructor that can receive these parameters.",
			"Remove the parameters or define a constructor inside the module.",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0017",
			"Full details:",
		),
	)

	return builder.String()
}
