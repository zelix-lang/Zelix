/*
   The Fluent Programming Language
   -----------------------------------------------------
   This code is released under the GNU GPL v3 license.
   For more information, please visit:
   https://www.gnu.org/licenses/gpl-3.0.html
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package error

import (
	"fluent/logger"
	"strings"
)

func MainFunctionIsPublic() string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("The main function cannot be public"))
	builder.WriteString(
		logger.BuildHelp(
			"If the main function is public, we risk",
			"exposing the program's internal structure.",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/Book/codes/E0028",
			"Full details:",
		),
	)

	return builder.String()
}
