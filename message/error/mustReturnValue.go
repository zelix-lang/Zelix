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

import "fluent/logger"

// MustReturnValue prints an error message when a function must return a value.
func MustReturnValue() {
	logger.Error("This function must return a value")
	logger.Help(
		"This function has not returned any value, but it must.",
		"If you don't want to return a value, use the 'nothing' type.",
	)
	logger.Info(
		"For more information, refer to:",
		"https://fluent-lang.github.io/book/codes/E0004",
		"Full details:",
	)
}
