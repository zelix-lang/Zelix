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

// InvalidDereference prints an error message indicating that a non-pointer value
// was dereferenced. It provides help and additional information on the correct
// usage of dereferencing in the Fluent programming language.
func InvalidDereference() {
	logger.Error("Dereference of a non-pointer value")
	logger.Help(
		"Dereferencing works for converting pointers to values.",
		"A value that is not a pointer cannot be dereferenced.",
	)
	logger.Info(
		"For more information, refer to:",
		"https://fluent-lang.github.io/book/codes/E0008",
		"Full details:",
	)
}
