package checker

import (
	"surf/ast"
	"surf/code"
	"surf/core/stack"
	"surf/logger"
	"surf/object"
)

var arithmeticOperators = map[code.TokenType]struct{}{
	code.Plus:     {},
	code.Minus:    {},
	code.Asterisk: {},
	code.Slash:    {},
}

// isOperator checks if the given token is an arithmetic operator
func isOperator(token code.Token) bool {
	_, found := arithmeticOperators[token.GetType()]
	return found
}

// checkOperator checks if the given token is an arithmetic operator
func checkOperator(token code.Token, statement []code.Token) {
	if !isOperator(token) || len(statement) == 1 {
		logger.TokenError(
			token,
			"Invalid arithmetic expression",
			"An arithmetic expression must start with a number or a variable",
			"Check the expression",
		)
	}
}

// checkStatement checks if the given statement is a valid arithmetic expression
func checkStatement(
	extractingIdentifier *bool,
	lastStatement *[]code.Token,
	currentToken code.Token,
	currentIndex int,
	variables *stack.StaticStack,
	functions *map[string]map[string]ast.Function,
) {
	*extractingIdentifier = false

	// Check the last statement
	if len(*lastStatement) == 0 {
		logger.TokenError(
			currentToken,
			"Unexpected token",
			"Expected an identifier or a number",
			"Check the arithmetic expression",
		)
	}

	// Analyze the last statement
	dummyType := object.NothingType
	dummyBool := false
	isArithmetic := false

	AnalyzeIdentifier(
		*lastStatement,
		variables,
		functions,
		&currentIndex,
		&dummyType,
		&isArithmetic,
		&dummyBool,
	)

	if !isArithmetic {
		logger.TokenError(
			currentToken,
			"Invalid arithmetic expression",
			"Expected an arithmetic expression",
			"Check the arithmetic expression",
		)
	}

	// Clear the slice
	*lastStatement = make([]code.Token, 0)
}

// AnalyzeArithmetic analyzes the given arithmetic expression
func AnalyzeArithmetic(
	statement []code.Token,
	variables *stack.StaticStack,
	functions *map[string]map[string]ast.Function,
) {
	statementLen := len(statement)
	if statementLen == 0 {
		return
	}

	// Check if there is an arithmetic operator
	checkOperator(statement[0], statement)

	// Used to check if the next token is an operator
	expectingOperator := true

	lastStatement := make([]code.Token, 0)
	extractingIdentifier := false

	// Used to skip indexes
	skipToIndex := 0

	for i, token := range statement {
		if i < skipToIndex {
			continue
		}

		if extractingIdentifier {
			if isOperator(token) {
				checkStatement(
					&extractingIdentifier,
					&lastStatement,
					token,
					i,
					variables,
					functions,
				)

				continue
			}

			lastStatement = append(lastStatement, token)
			continue
		}

		if expectingOperator {
			checkOperator(token, statement[i:])
			expectingOperator = false

			continue
		}

		tokenType := token.GetType()

		switch tokenType {
		case code.Identifier:
			extractingIdentifier = true
			lastStatement = append(lastStatement, token)

			break
		case code.DecimalLiteral:
			expectingOperator = true
			break
		case code.NumLiteral:
			expectingOperator = true
			break
		default:
			logger.TokenError(
				token,
				"Unexpected token",
				"Expected an identifier or a number",
				"Check the arithmetic expression",
			)
		}

	}

	// Check any remaining statement
	if len(lastStatement) > 0 {
		checkStatement(
			&extractingIdentifier,
			&lastStatement,
			statement[len(statement)-1],
			len(statement)-1,
			variables,
			functions,
		)
	}
}
