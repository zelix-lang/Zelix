package analyzer

import (
	"fluent/code"
	"fluent/code/mod"
	"fluent/code/types"
	"fluent/logger"
	"fluent/stack"
	"fluent/token"
	"fluent/tokenUtil/splitter"
)

// boolOperators contains all the boolean operators
// use a map for O(1) lookup time
var boolOperators = map[token.Type]struct{}{
	token.And: {},
	token.Or:  {},
	token.Not: {},
}

// analyzeSingleBool analyzes a single boolean expression
func analyzeSingleBool(
	unit token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.FluentMod,
) {
	if unit.GetType() == token.Identifier {
		// Retrieve variables
		variable := AnalyzeStatement(
			[]token.Token{unit},
			variables,
			functions,
			mods,
			dummyNothingType,
		)

		typeWrapper := variable.GetType()

		if typeWrapper.GetType() != types.BooleanType {
			logger.TokenError(
				unit,
				"Invalid type for boolean expression",
				"This variable should be a boolean",
				"Check the variable type and its value",
			)
		}

		return
	}

	if unit.GetType() == token.BoolLiteral {
		logger.TokenError(
			unit,
			"Invalid boolean expression",
			"Expected a boolean expression",
			"Check the boolean expression",
		)
	}
}

// AnalyzeBool analyzes the given boolean expression
func AnalyzeBool(
	statement []token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.FluentMod,
	trace token.Token,
) {
	if len(statement) == 0 {
		logger.TokenError(
			trace,
			"Unexpected token",
			"Expected a boolean expression",
			"Build a boolean expression like 'true' or 'false'",
		)
	}

	if len(statement) == 1 {
		analyzeSingleBool(statement[0], variables, functions, mods)
		return
	}

	// Parse expressions with parentheses
	startAt := 0
	firstToken := statement[0]
	firstTokenType := firstToken.GetType()

	// Exclude negations
	for firstTokenType == token.Not {
		if startAt == len(statement)-1 {
			// All tokens were consumed
			logger.TokenError(
				firstToken,
				"Invalid boolean expression",
				"This expression is constructed purely of negations",
				"Build expressions like: '!true' or '!false'",
			)
		}

		startAt++
		firstToken = statement[startAt]
		firstTokenType = firstToken.GetType()
	}

	hasProcessedParens := false
	// The remaining statement after the negations
	remainingStatement := statement[startAt:]

	if firstTokenType == token.OpenParen {
		// Recursively analyze the statement
		statementBeforeParen, _ := splitter.ExtractTokensBefore(
			remainingStatement,
			token.CloseParen,
			true,
			token.OpenParen,
			token.CloseParen,
			true,
		)

		startAt += len(statementBeforeParen) + 1
		hasProcessedParens = true

		// Recursively analyze the statement inside the parentheses
		AnalyzeBool(
			statementBeforeParen[1:],
			variables,
			functions,
			mods,
			trace,
		)
	}

	// No more tokens to analyze
	if startAt == len(statement) {
		return
	}

	firstToken = statement[startAt]
	firstTokenType = firstToken.GetType()
}
