package fun

import (
	"zyro/code"
	"zyro/core/stack"
	"zyro/object"
	"zyro/token"
	"zyro/tokenUtil"
	"zyro/util"
)

// CallFun interprets a function and executes it
func CallFun(
	function *code.Function,
	runtime map[string]func(...object.ZyroObject),
	functions *map[string]map[string]*code.Function,
	args ...object.ZyroObject,
) {
	variables := stack.NewStack()
	// Create a new scope in the variables map
	variables.CreateScope()

	argKeys := util.MapKeys(function.GetParameters())

	// Put the arguments into variables
	for i, key := range argKeys {
		variables.Append(
			key,
			args[i],
			false,
		)
	}

	// Used to skip indexes
	skipToIndex := 0

	for i, unit := range function.GetBody() {
		if i < skipToIndex && skipToIndex > 0 {
			continue
		}

		tokenType := unit.GetType()

		if tokenType == token.Identifier {
			// Extract the statement
			statement := tokenUtil.ExtractTokensBefore(
				function.GetBody()[i:],
				token.Semicolon,
				// Don't handle nested statements here
				false,
				token.Unknown,
				token.Unknown,
				true,
			)

			CallStatement(statement, runtime, function.IsStd(), functions, variables)

			// Don't subtract 1 because the statement doesn't contain the semicolon
			skipToIndex = i + len(statement)
		}
	}
}
