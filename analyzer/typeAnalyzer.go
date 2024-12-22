package analyzer

import (
	"zyro/code"
	"zyro/core/stack"
	"zyro/logger"
	"zyro/object"
	"zyro/token"
)

// AnalyzeType analyzes the given type and makes sure it matches the expected type
func AnalyzeType(
	statement []token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*code.ZyroMod,
	expected object.ZyroObject,
) {
	isMod := expected.GetType() == object.ModType

	// Analyze the type
	value := AnalyzeStatement(
		statement,
		variables,
		functions,
		mods,
	)

	if isMod {
		mod := expected.GetValue().(*code.ZyroMod)

		if value.GetType() != object.ModType {
			logger.TokenError(
				statement[0],
				"Type mismatch",
				"This type does not match the value type",
				"Change the declaration or remove the assignment",
			)
		}

		gotMod := value.GetValue().(*code.ZyroMod)

		if value.GetType() != object.ModType || gotMod.GetName() != mod.GetName() {
			logger.TokenError(
				statement[0],
				"Type mismatch",
				"This type does not match the value type",
				"Change the declaration or remove the assignment",
			)
		}
	} else if value.GetType() != expected.GetType() {
		logger.TokenError(
			statement[0],
			"Type mismatch",
			"This type does not match the value type",
			"Change the declaration or remove the assignment",
		)
	}
}
