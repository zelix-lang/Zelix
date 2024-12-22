package analyzer

import (
	"zyro/code"
	"zyro/core/stack"
	"zyro/logger"
	"zyro/object"
	"zyro/token"
	"zyro/tokenUtil"
)

// AnalyzeStatement analyzes the given statement
// and returns the object type that the statement
// returns
func AnalyzeStatement(
	statement []token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*code.ZyroMod,
) object.ZyroObject {
	// Used to know what to check for
	isArithmetic := false
	isFunCall := false

	// Used to check property access
	// i.e.: object.property
	lastValue := object.NewZyroObject(object.NothingType, nil)
	startAt := 0

	firstToken := statement[0]
	firstTokenType := firstToken.GetType()

	switch firstTokenType {
	case token.New:
		AnalyzeObjectCreation(
			tokenUtil.ExtractTokensBefore(
				statement,
				token.Dot,
				true,
				token.OpenParen,
				token.CloseParen,
				false,
			),
			variables,
			functions,
			mods,
			&startAt,
			&lastValue,
		)

		break
	case token.Identifier:
		AnalyzeIdentifier(
			statement,
			variables,
			functions,
			mods,
			&startAt,
			&lastValue,
			&isArithmetic,
			&isFunCall,
		)

		break
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
			mods,
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

	// Analyze reassignments
	if remainingStatement[0].GetType() == token.Assign {
		if isFunCall {
			logger.TokenError(
				remainingStatement[0],
				"Invalid operation",
				"Cannot assign to a method call",
				"Check the statement",
			)
		}

		variable, _ := variables.Load(firstToken.GetValue())
		// No need to check if it was found here, it was already checked

		if variable.IsConstant() {
			logger.TokenError(
				remainingStatement[0],
				"Cannot reassign constant",
				"Check the variable declaration",
			)
		}

		return lastValue
	}

	// The only valid operation after all that has been processed
	// is property access, therefore the fist token of the remaining
	// statement must be a dot
	if remainingStatement[0].GetType() != token.Dot {
		logger.TokenError(
			remainingStatement[0],
			"Invalid operation",
			"Invalid operation after identifier",
			"Check the statement",
		)
	}

	// Get tokens before an assignment
	// i.e.: object.property = value
	beforeAssignment := tokenUtil.ExtractTokensBefore(
		remainingStatement[1:],
		token.Assign,
		false,
		token.Unknown,
		token.Unknown,
		false,
	)

	// if beforeAssignment is empty, that means that
	// the statement ends in a dot: "object.property."
	// which is invalid
	if len(beforeAssignment) == 0 {
		logger.TokenError(
			remainingStatement[0],
			"Invalid operation",
			"Invalid operation after identifier",
			"Check the statement",
		)
	}

	// +1 for the dot
	// +1 for the assignment
	afterAssignment := remainingStatement[len(beforeAssignment)+2:]

	// Reset isFunCall to catch assignments to methods
	isFunCall = false
	props := tokenUtil.SplitTokens(
		beforeAssignment,
		token.Dot,
		token.OpenParen,
		token.CloseParen,
	)

	// Analyze all props
	for _, prop := range props {
		AnalyzePropAccess(
			prop,
			variables,
			functions,
			mods,
			&lastValue,
			&isFunCall,
			len(afterAssignment) > 0,
		)
	}

	if isFunCall && len(afterAssignment) > 0 {
		logger.TokenError(
			afterAssignment[0],
			"Invalid operation",
			"Cannot assign to a method call",
			"Check the statement",
		)
	}

	// Analyze assignment
	AnalyzeType(
		afterAssignment,
		variables,
		functions,
		mods,
		lastValue,
	)

	return lastValue
}
