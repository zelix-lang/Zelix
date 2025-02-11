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

// DataOutlivesStack prints an error message indicating that a value outlives the lifetime of its stack.
// It provides help and information on why this issue occurs.
func DataOutlivesStack() string {
	builder := strings.Builder{}

	builder.WriteString(logger.BuildError("This value outlives the lifetime of its stack"))
	builder.WriteString(
		logger.BuildHelp(
			"A pointer is used to pass references to values between functions",
			"When it goes out of scope, it gets deleted; hence, the value becomes invalid",
			"When you later try to use a reference that no longer exists, you will most likely",
			"cause a SIGSEGV or produce unpredictable results in the worst case scenario.",
		),
	)
	builder.WriteString(
		logger.BuildInfo(
			"For more information, refer to:",
			"https://fluent-lang.github.io/book/codes/E0003",
			"Full details:",
		),
	)

	return builder.String()
}
