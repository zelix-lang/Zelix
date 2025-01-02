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

var arithmeticOperators = map[token.Type]struct{}{
	token.Plus:     {},
	token.Minus:    {},
	token.Asterisk: {},
	token.Slash:    {},
	token.Percent:  {},
}

// isOperator checks if the given token is an arithmetic operator
func isArithmeticOperator(token token.Token) bool {
	_, found := arithmeticOperators[token.GetType()]
	return found
}

// checkOperator checks if the given token is an arithmetic operator
func checkArithmeticOperator(token token.Token, statement []token.Token) {
	if !isArithmeticOperator(token) || len(statement) == 1 {
		logger.TokenError(
			token,
			"Invalid arithmetic expression",
			"An arithmetic expression must start with a number or a variable",
			"Got: "+token.GetValue(),
			"Check the expression",
		)
	}
}

// checkStatement checks if the given statement is a valid arithmetic expression
func checkStatement(
	extractingIdentifier *bool,
	lastStatement *[]token.Token,
	currentToken token.Token,
	currentIndex int,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.FluentMod,
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
	dummyObj := wrapper.NewFluentObject(
		dummyNothingType,
		nil,
	)

	dummyBool := false
	isArithmetic := false
	isBool := false

	AnalyzeIdentifier(
		*lastStatement,
		variables,
		functions,
		mods,
		&currentIndex,
		&dummyObj,
		&isArithmetic,
		&dummyBool,
		&isBool,
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
	*lastStatement = make([]token.Token, 0)
}

// AnalyzeArithmetic analyzes the given arithmetic expression
func AnalyzeArithmetic(
	statement []token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.FluentMod,
	lastValue *wrapper.FluentObject,
) {
	statementLen := len(statement)
	if statementLen == 0 {
		return
	}

	// There is an edge case where this can be a boolean expression
	// for example: 1.25 + 2.5 == 3.75
	// since the first number is a decimal, this analyzer is called
	// because the first token is a decimal, and it is most likely
	// to be an arithmetic expression
	if isBooleanOperator(statement[0]) {
		// Directly call the boolean analyzer
		AnalyzeBool(
			statement[1:],
			variables,
			functions,
			mods,
			statement[0],
		)

		*lastValue = wrapper.NewFluentObject(
			wrapper.ForceNewTypeWrapper(
				"bool",
				make([]wrapper.TypeWrapper, 0),
				types.BooleanType,
			),
			true,
		)
		return
	}

	// Check if there is an arithmetic operator
	checkArithmeticOperator(statement[0], statement)

	// Used to check if the next token is an operator
	expectingOperator := true

	lastStatement := make([]token.Token, 0)
	extractingIdentifier := false

	// Used to skip indexes
	skipToIndex := 0

	// The depth of parentheses
	parenDepth := 0

	for i, unit := range statement {
		if i < skipToIndex {
			continue
		}

		if extractingIdentifier {
			if isArithmeticOperator(unit) {
				checkStatement(
					&extractingIdentifier,
					&lastStatement,
					unit,
					i,
					variables,
					functions,
					mods,
				)

				continue
			}

			lastStatement = append(lastStatement, unit)
			continue
		}

		if isBooleanOperator(unit) {
			// Directly call the boolean analyzer
			AnalyzeBool(
				statement[1:],
				variables,
				functions,
				mods,
				unit,
			)

			*lastValue = wrapper.NewFluentObject(
				wrapper.ForceNewTypeWrapper(
					"bool",
					make([]wrapper.TypeWrapper, 0),
					types.BooleanType,
				),
				true,
			)
			return
		}

		if expectingOperator {
			if unit.GetType() == token.CloseParen {
				parenDepth--

				if parenDepth < 0 {
					logger.TokenError(
						unit,
						"Invalid arithmetic expression",
						"Unexpected closing parenthesis",
						"Check the arithmetic expression",
					)
				}

				continue
			}

			checkArithmeticOperator(unit, statement[i:])
			expectingOperator = false

			continue
		}

		tokenType := unit.GetType()

		switch tokenType {
		case token.OpenParen:
			parenDepth++
			continue
		case token.CloseParen:
			// If this case is reached, it means
			// the code has something like: "()", which is invalid
			logger.TokenError(
				unit,
				"Invalid arithmetic expression",
				"An arithmetic expression must contain at least one number or variable",
			)
		case token.Identifier:
			extractingIdentifier = true
			lastStatement = append(lastStatement, unit)
			continue
		case token.DecimalLiteral, token.NumLiteral:
			expectingOperator = true
			continue
		default:
			logger.TokenError(
				unit,
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
			mods,
		)
	}
}
