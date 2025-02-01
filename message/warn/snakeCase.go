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

// SnakeCase logs a warning saying that the provided name is not in snake_case format.
// It also provides a link to more information about the warning.
//
// Parameters:
//   - name: The name to be logged.
func SnakeCase(name string) {
	logger.Warn("The value name '" + name + "' should be in snake_case.")
	logger.Info(
		"For more information, refer to:",
		"https://fluent-lang.github.io/book/codes/E0001",
		"Full details:",
	)
}
