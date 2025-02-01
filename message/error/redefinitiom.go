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

// Redefinition prints the error message for redefining a value.
//
// Parameters:
//
//	name - The name of the value that is being redefined.
func Redefinition(name string) {
	logger.Error("Redefinition of value '" + name + "'")
	logger.Help(
		"This value is already defined in the current scope.",
		"Please, change the name of the value or remove the redefinition.",
	)
	logger.Info(
		"For more information, refer to:",
		"https://fluent-lang.github.io/book/codes/E0007",
		"Full details:",
	)
}
