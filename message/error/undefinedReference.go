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

// UndefinedReference prints an error message when a reference to a value is undefined.
//
// Parameters:
//
//	name: the name of the value.
func UndefinedReference(name string) {
	logger.Error("Undefined reference to value '" + name + "'")
	logger.Help(
		"This value is not defined in the current scope.",
		"Make sure you have defined the value before using it.",
	)
	logger.Info(
		"For more information, refer to:",
		"https://fluent-lang.github.io/book/codes/E0005",
		"Full details:",
	)
}
