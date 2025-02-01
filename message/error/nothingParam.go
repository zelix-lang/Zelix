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

// NothingParam prints error and help messages related to the misuse of the 'nothing' type as a parameter.
func NothingParam() {
	logger.Error("Parameters cannot have type 'nothing'")
	logger.Help(
		"'nothing' is used for return types",
		"A parameter that holds 'nothing' is not useful in a function",
		"Change the type of the parameter to a valid type",
	)
	logger.Info(
		"For more information, refer to:",
		"https://fluent-lang.github.io/book/codes/E0002",
		"Full details:",
	)
}
