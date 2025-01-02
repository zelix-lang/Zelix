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

// AnalyzeReturn analyzes a return statement
func AnalyzeReturn(
	statement []token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.FluentMod,
	returnVal wrapper.FluentObject,
) {
	// Check if the function returns something whilst it shouldn't
	returnType := returnVal.GetType()
	if len(statement) > 1 && returnType.GetType() == types.NothingType {
		logger.TokenError(
			statement[0],
			"Invalid return statement",
			"This function is not supposed to return anything",
			"Remove this return statement",
			"Or change the function definition",
		)
	}

	// Check if the function doesn't return anything whilst it should
	if len(statement) < 1 && returnType.GetType() != types.NothingType {
		logger.TokenError(
			statement[0],
			"Invalid return statement",
			"This function should return something",
			"Add a return statement",
			"Or change the function definition",
		)
	}

	// Drop empty return statements
	if len(statement) < 1 {
		return
	}

	// Check for type mismatch
	AnalyzeType(
		statement,
		variables,
		functions,
		mods,
		returnVal,
	)
}
