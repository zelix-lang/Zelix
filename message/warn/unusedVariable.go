/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package warn

import "fluent/logger"

// UnusedVariable logs a warning about an unused variable and provides help and information links.
// Parameters:
// - name: The name of the unused variable.
func UnusedVariable(name *string) {
	logger.Warn("Unused value: '" + *name + "'.")
	logger.Help(
		"If this is intentional, add an underscore before the variable name.",
	)
	logger.Info(
		"For more information, refer to:",
		"https://fluent-lang.github.io/book/codes/E0006",
		"Inside of this function:",
	)
}
