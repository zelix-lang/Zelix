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

func ShouldNotReturn() {
	logger.Error("This function should not return a value")
	logger.Help(
		"Functions with the 'nothing' return type should not return a value.",
		"Remove the return statement from the function or modify the return type.",
	)
	logger.Info(
		"For more information, refer to:",
		"https://fluent-lang.github.io/book/codes/E0012",
		"Full details:",
	)
}
