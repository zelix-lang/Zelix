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
	"strings"
)

func CannotInferType() string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("Cannot infer this object's type"))
	builder.WriteString(
		logger.BuildHelp(
			"Fluent cannot infer the type of this object.",
			"Type inferring works for assignments, parameters and declarations",
			"Statements like [1, 2, 3] cannot have their type inferred because",
			"at such point, the Fluent compiler does not expect a type to infer or match",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0011",
			"Full details:",
		),
	)

	return builder.String()
}
