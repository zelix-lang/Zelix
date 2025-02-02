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

// ParamCountMismatch logs an error message indicating a mismatch in the count of parameters.
// It provides a help message with the expected count and a link to more information.
// Parameters:
//   - count: A string representing the expected number of parameters.
func ParamCountMismatch(count string) {
	logger.Error("Mismatched count of parameters")
	logger.Help(
		"This function expected "+count+" parameters.",
		"Make sure you have passed the correct number of parameters.",
	)
	logger.Info(
		"For more information, refer to:",
		"https://fluent-lang.github.io/book/codes/E0010",
		"Full details:",
	)
}
