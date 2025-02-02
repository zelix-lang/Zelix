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

import (
	"fluent/logger"
)

func TypeMismatch() {
	logger.Error("Type mismatch")
	logger.Help(
		"The expected type does not match the gotten type.",
		"Modify the type to match the expected one, for example \"2\" instead of 2",
	)
	logger.Info(
		"For more information, refer to:",
		"https://fluent-lang.github.io/book/codes/E0009",
		"Full details:",
	)
}
