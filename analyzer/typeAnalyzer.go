package analyzer

import (
	"fluent/code"
	"fluent/code/mod"
	"fluent/code/types"
	"fluent/code/wrapper"
	"fluent/logger"
	"fluent/stack"
	"fluent/token"
)

// AnalyzeType analyzes the given type and makes sure it matches the expected type
func AnalyzeType(
	statement []token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.FluentMod,
	expected wrapper.FluentObject,
) {
	expectedTypeWrapper := expected.GetType()
	AnalyzeGeneric(expectedTypeWrapper, mods, statement[0])

	isMod := expectedTypeWrapper.GetType() == types.ModType

	// Analyze the type
	value := AnalyzeStatement(
		statement,
		variables,
		functions,
		mods,
		expectedTypeWrapper,
	)

	valueTypeWrapper := value.GetType()

	if isMod {
		if valueTypeWrapper.GetType() != types.ModType {
			logger.TokenError(
				statement[0],
				"Type mismatch",
				"This type does not match the value type",
				"Change the declaration or remove the assignment",
				"Expected: "+expectedTypeWrapper.Marshal(),
				"Got: "+valueTypeWrapper.Marshal(),
			)
		}

		if !expectedTypeWrapper.Compare(valueTypeWrapper) {
			logger.TokenError(
				statement[0],
				"Type mismatch",
				"This type does not match the value type",
				"Change the declaration or remove the assignment",
				"Expected: "+expectedTypeWrapper.Marshal(),
				"Got: "+valueTypeWrapper.Marshal(),
			)
		}
	} else if !expectedTypeWrapper.Compare(valueTypeWrapper) {
		logger.TokenError(
			statement[0],
			"Type mismatch",
			"This type does not match the value type",
			"Change the declaration or remove the assignment",
			"Expected: "+expectedTypeWrapper.Marshal(),
			"Got: "+valueTypeWrapper.Marshal(),
		)
	}
}
