package analyzer

import (
	"fluent/code"
	"fluent/code/mod"
	"fluent/code/types"
	"fluent/code/wrapper"
	"fluent/logger"
	"fluent/stack"
	"fluent/token"
	"fluent/tokenUtil/splitter"
)

// AnalyzeFor analyzes the given for loop
func AnalyzeFor(
	declaration []token.Token,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.FluentMod,
	variables *stack.Stack,
) {
	if len(declaration) < 5 {
		logger.TokenError(
			declaration[0],
			"Invalid for loop",
			"A for loop must have the form 'for var in iterable'",
			"Check the for loop",
		)
	}

	// First token is the variable name
	variableNameToken := declaration[0]
	variableName := variableNameToken.GetValue()

	if variableNameToken.GetType() != token.Identifier {
		logger.TokenError(
			variableNameToken,
			"Invalid variable name",
			"A variable name must be an identifier",
			"Change the variable name",
		)
	}

	if _, ok := variables.Load(variableName); ok {
		logger.TokenError(
			variableNameToken,
			"Redefinition of value '"+variableName+"'",
			"Change the variable name",
		)
	}

	// Second token should always be "in"
	if declaration[1].GetType() != token.In {
		logger.TokenError(
			declaration[1],
			"Invalid for loop",
			"The second token in a for loop must be 'in'",
			"Check the for loop",
		)
	}

	iterationRaw := declaration[2:]
	iterations := splitter.SplitTokens(
		iterationRaw,
		token.Colon,
		token.Unknown,
		token.Unknown,
	)

	// The iterations should always be 2 of length
	if len(iterations) != 2 {
		logger.TokenError(
			declaration[2],
			"Invalid for loop",
			"The iterations must be of the form 'start:stop'",
			"Check the for loop",
		)
	}

	// Extract the statements
	// "1 : 5" => [1], [5]
	statement1 := iterations[0]
	statement2 := iterations[1]

	// Analyze the statements
	result1 := AnalyzeStatement(
		statement1,
		variables,
		functions,
		mods,
		dummyNothingType,
	)

	result2 := AnalyzeStatement(
		statement2,
		variables,
		functions,
		mods,
		dummyNothingType,
	)

	result1TypeWrapper := result1.GetType()
	result2TypeWrapper := result2.GetType()

	if !result1TypeWrapper.Compare(result2TypeWrapper) && (result1TypeWrapper.GetType() != types.IntType || result1TypeWrapper.GetType() != types.DecimalType) {
		logger.TokenError(
			declaration[2],
			"Invalid for loop",
			"Both statements must be of the same type",
		)
	}

	// Store the variable in the stack
	variables.CreateScope()
	variables.Append(variableName, wrapper.NewFluentObject(
		wrapper.ForceNewTypeWrapper(
			"num",
			make([]wrapper.TypeWrapper, 0),
			types.IntType,
		),
		0,
	), false)

	// Scope is automatically destroyed when finding a closing bracket
	// no need to destroy it here or keep track of the loop's body
}
