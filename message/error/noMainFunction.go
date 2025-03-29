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

// NoMainFunction generates an error message indicating that the main function is not defined.
// It provides guidance on how to define the main function and where to find more information.
func NoMainFunction() string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("The main function is not defined"))
	builder.WriteString(
		logger.BuildHelp(
			"Define a main function with the signature 'fun main()'",
			"to start the program",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/Book/codes/E0023",
			"Full details:",
		),
	)

	return builder.String()
}
