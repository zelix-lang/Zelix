package analyzer

import (
	"zyro/code"
	"zyro/code/mod"
	"zyro/code/types"
	"zyro/code/wrapper"
	"zyro/logger"
	"zyro/stack"
	"zyro/token"
)

// AnalyzeType analyzes the given type and makes sure it matches the expected type
func AnalyzeType(
	statement []token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.ZyroMod,
	expected wrapper.ZyroObject,
) {
	expectedTypeWrapper := expected.GetType()
	isMod := expectedTypeWrapper.GetType() == types.ModType

	// Analyze the type
	value := AnalyzeStatement(
		statement,
		variables,
		functions,
		mods,
	)

	if isMod {
		module := expected.GetValue().(*mod.ZyroMod)
		valueTypeWrapper := value.GetType()

		if valueTypeWrapper.GetType() != types.ModType {
			logger.TokenError(
				statement[0],
				"Type mismatch",
				"This type does not match the value type",
				"Change the declaration or remove the assignment",
			)
		}

		gotMod := value.GetValue().(*mod.ZyroMod)

		if valueTypeWrapper.GetType() != types.ModType || gotMod.GetName() != module.GetName() {
			logger.TokenError(
				statement[0],
				"Type mismatch",
				"This type does not match the value type",
				"Change the declaration or remove the assignment",
			)
		}
	} else if !expectedTypeWrapper.Compare(value.GetType()) {
		logger.TokenError(
			statement[0],
			"Type mismatch",
			"This type does not match the value type",
			"Change the declaration or remove the assignment",
		)
	}
}
