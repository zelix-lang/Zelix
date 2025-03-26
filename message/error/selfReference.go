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

// SelfReference generates an error message indicating that a module cannot end
// without creating a copy of itself. It provides detailed information about
// why self-references are not allowed and includes a link to further documentation.
func SelfReference() string {
	builder := strings.Builder{}

	// Add an error message to the builder
	builder.WriteString(logger.BuildError("This module's block cannot end without creating a copy of itself"))

	// Add a help message to the builder explaining the issue with self-references
	builder.WriteString(
		logger.BuildHelp(
			"Things like: let my_mod: MyMod = new MyMod() inside a function",
			"that belongs to MyMod are not allowed because they create a self reference,",
			"which leads to an infinite loop.",
		),
	)

	// Add an informational message to the builder with a link to further documentation
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/Book/codes/E0021",
			"Full details:",
		),
	)

	return builder.String()
}
