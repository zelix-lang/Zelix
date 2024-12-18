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
	variables *stack.StaticStack,
	functions *map[string]map[string]*ast.Function,
) object.SurfObjectType {
	// Used to know what to check for
	isArithmetic := false
	isFunCall := false

	// Used to check property access
	// i.e.: object.property
	lastValueType := object.NothingType
	startAt := 0

	firstToken := statement[0]
	firstTokenType := firstToken.GetType()

	switch firstTokenType {
	case code.Identifier:
		AnalyzeIdentifier(
			statement,
			variables,
			functions,
			&startAt,
			&lastValueType,
			&isArithmetic,
			&isFunCall,
		)
	default:
		lastValueType = tokenUtil.ToObjType(firstToken, variables)
		isArithmetic = lastValueType == object.IntType || lastValueType == object.DecimalType
		startAt = 1
	}

	// Analyze the rest of the statement
	remainingStatement := statement[startAt:]

	if len(remainingStatement) == 0 {
		return lastValueType
	}

	if lastValueType == object.NothingType {
		logger.TokenError(
			remainingStatement[0],
			"Illegal property access",
			"Cannot access properties of a non-object",
			"Check the object type",
		)
	}

	if isArithmetic {
		AnalyzeArithmetic(
			remainingStatement,
			variables,
			functions,
		)

		return lastValueType
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

	return lastValueType
}
