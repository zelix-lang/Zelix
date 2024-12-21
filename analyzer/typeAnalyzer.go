package analyzer

import (
	"surf/code"
	"surf/core/stack"
	"surf/logger"
	"surf/object"
	"surf/token"
)

// AnalyzeType analyzes the type of the given assignment
func AnalyzeType(
	assignment []token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]*code.SurfMod,
	expected object.SurfObject,
) {
	isMod := expected.GetType() == object.ModType

	// Analyze the type
	value := AnalyzeStatement(
		assignment,
		variables,
		functions,
		mods,
	)

	if isMod {
		mod := expected.GetValue().(*code.SurfMod)

		if value.GetType() != object.ModType {
			logger.TokenError(
				assignment[0],
				"Type mismatch",
				"This type does not match the value type",
				"Change the declaration or remove the assignment",
			)
		}

		gotMod := value.GetValue().(*code.SurfMod)

		if value.GetType() != object.ModType || gotMod.GetName() != mod.GetName() {
			logger.TokenError(
				assignment[0],
				"Type mismatch",
				"This type does not match the value type",
				"Change the declaration or remove the assignment",
			)
		}
	} else if value.GetType() != expected.GetType() {
		logger.TokenError(
			assignment[0],
			"Type mismatch",
			"This type does not match the value type",
			"Change the declaration or remove the assignment",
		)
	}
}
