package analyzer

import (
	"surf/ast"
	"surf/code"
	"surf/core/stack"
	"surf/logger"
	"surf/object"
	"surf/tokenUtil"
)

// AnalyzeStatement analyzes the given statement
// and returns the object type that the statement
// returns
func AnalyzeStatement(
	statement []code.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*ast.Function,
) object.SurfObject {
	// Used to know what to check for
	isArithmetic := false
	isFunCall := false

	// Used to check property access
	// i.e.: object.property
	lastValue := object.NewSurfObject(object.NothingType, nil)
	startAt := 0

	firstToken := statement[0]
	firstTokenType := firstToken.GetType()

	switch firstTokenType {
	case code.New:
		// TODO! Parse object creation
	case code.Identifier:
		AnalyzeIdentifier(
			statement,
			variables,
			functions,
			&startAt,
			&lastValue,
			&isArithmetic,
			&isFunCall,
		)
	default:
		lastValue = tokenUtil.ToObj(firstToken, variables)
		isArithmetic = lastValue.GetType() == object.IntType || lastValue.GetType() == object.DecimalType
		startAt = 1
	}

	// Analyze the rest of the statement
	remainingStatement := statement[startAt:]

	if len(remainingStatement) == 0 {
		return lastValue
	}

	if isArithmetic {
		AnalyzeArithmetic(
			remainingStatement,
			variables,
			functions,
		)

		return lastValue
	}

	if lastValue.GetType() != object.ModType {
		logger.TokenError(
			remainingStatement[0],
			"Illegal property access",
			"Cannot access properties of a non-object",
			"Check the object type",
		)
	}

	// The only valid operation after all that has been processed
	// is property access, therefore the fist token of the remaining
	// statement must be a dot
	if remainingStatement[0].GetType() != code.Dot {
		logger.TokenError(
			remainingStatement[0],
			"Invalid operation",
			"Invalid operation after identifier",
			"Check the statement",
		)
	}

	// TODO! Parse property access

	return lastValue
}
