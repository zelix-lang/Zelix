package analyzer

import (
	"fluent/code"
	"fluent/code/mod"
	"fluent/logger"
	"fluent/stack"
	"fluent/token"
	"fluent/tokenUtil/splitter"
)

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
		if statement[0].GetType() != token.BoolLiteral {
			logger.TokenError(
				statement[0],
				"Invalid boolean expression",
				"Expected a boolean expression",
				"Check the boolean expression",
			)
		}

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

	if firstTokenType == token.OpenParen {
		// Recursively analyze the statement
		statementBeforeParen, _ := splitter.ExtractTokensBefore(
			statement[startAt:],
			token.CloseParen,
			true,
			token.OpenParen,
			token.CloseParen,
			true,
		)

		startAt = len(statementBeforeParen) + 1

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
