package analyzer

import (
	"fluent/ast"
	"fluent/code"
	"fluent/code/mod"
	"fluent/code/types"
	"fluent/code/wrapper"
	"fluent/logger"
	"fluent/stack"
	"fluent/token"
	"fluent/tokenUtil/splitter"
)

// boolOperators contains all the boolean operators
// use a map for O(1) lookup time
var boolOperators = map[token.Type]struct{}{
	token.And:   {},
	token.Or:    {},
	token.Equal: {},
}

// isBooleanOperator checks if the given token is a boolean operator
func isBooleanOperator(unit token.Token) bool {
	_, found := boolOperators[unit.GetType()]
	return found
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

	if unit.GetType() != token.BoolLiteral {
		logger.TokenError(
			unit,
			"Invalid boolean expression",
			"Expected a literal or a variable",
			"Check the boolean expression and its value",
		)
	}
}

// analyzeParenStatement analyzes the given statement inside parentheses
func analyzeParenStatement(
	remainingStatement []token.Token,
	startAt *int,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.FluentMod,
	trace token.Token,
) wrapper.TypeWrapper {
	// Recursively analyze the statement
	statementBeforeParen, _ := splitter.ExtractTokensBefore(
		remainingStatement,
		token.CloseParen,
		true,
		token.OpenParen,
		token.CloseParen,
		true,
	)

	*startAt += len(statementBeforeParen) + 1

	// Recursively analyze the statement inside the parentheses
	return AnalyzeBool(
		statementBeforeParen[1:],
		variables,
		functions,
		mods,
		trace,
	)
}

// checkOperator checks if the given token is an operator
func checkBoolOperator(unit token.Token, statement []token.Token) {
	if !isBooleanOperator(unit) || len(statement) == 1 {
		logger.TokenError(
			unit,
			"Invalid boolean operator",
			"Expected a boolean operator like '&&', '||'",
			"Construct the logical expression correctly",
			"Got: "+unit.GetValue(),
		)
	}
}

// analyzeAssignmentTypes checks if the given types are compatible
func analyzeAssignmentTypes(lastValue *wrapper.TypeWrapper, newValue wrapper.TypeWrapper, firstToken token.Token) {
	if !lastValue.Compare(dummyNothingType) && !lastValue.Compare(newValue) {
		logger.TokenError(
			firstToken,
			"Type mismatch",
			"Cannot compare or join different types",
			"This typically happens when you try to compare types that aren't compatible",
			"For example: 1.2 == 1 or \"hello\" == 1",
		)
	}
}

// analyzeExpression analyzes the given expression
// and returns the resulting value
func analyzeBoolExpression(
	statement []token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.FluentMod,
	startAt *int,
	expectingOperator *bool,
	expectingArithmeticOperator *bool,
	lastValue *wrapper.TypeWrapper,
) {
	firstToken := statement[*startAt]
	firstTokenType := firstToken.GetType()

	if firstTokenType == token.OpenParen {
		newValue := analyzeParenStatement(
			statement,
			startAt,
			variables,
			functions,
			mods,
			firstToken,
		)

		// Wait for the next operator
		newType := newValue.GetType()

		// Check for type mismatch
		analyzeAssignmentTypes(lastValue, newValue, firstToken)

		*lastValue = newValue
		*expectingArithmeticOperator = newType == types.IntType || newType == types.DecimalType
		*expectingOperator = !*expectingArithmeticOperator

		return
	}

	if firstTokenType == token.BoolLiteral {
		analyzeSingleBool(firstToken, variables, functions, mods)
		newValue := wrapper.ForceNewTypeWrapper(
			"bool",
			make([]wrapper.TypeWrapper, 0),
			types.BooleanType,
		)

		// Check for type mismatch
		analyzeAssignmentTypes(lastValue, newValue, firstToken)

		*lastValue = newValue

		return
	}

	if firstTokenType == token.Identifier {
		// Retrieve variables
		variable, varFound := variables.Load(firstToken.GetValue())
		fun, funFound, isPublic := ast.LocateFunction(*functions, firstToken.GetFile(), firstToken.GetValue())

		if !varFound && !funFound {
			logger.TokenError(
				firstToken,
				"Undefined reference to variable "+firstToken.GetValue(),
				"The variable "+firstToken.GetValue()+" was not found in the current scope",
				"Add the variable to the current scope",
			)
		}

		// Check if the function is public
		if funFound && !isPublic && fun.GetTrace().GetFile() != firstToken.GetFile() {
			logger.TokenError(
				firstToken,
				"Function "+firstToken.GetValue()+" is not public",
				"Move the function to the current file or make it public",
			)
		}

		if varFound {
			value := variable.GetValue()
			*lastValue = value.GetType()
			*startAt += 1
		} else {
			// Parse the function call
			argsRaw, _ := splitter.ExtractTokensBefore(
				statement[*startAt:],
				token.CloseParen,
				true,
				token.OpenParen,
				token.CloseParen,
				true,
			)

			// Add the close parenthesis
			argsRaw = append(argsRaw, statement[*startAt+len(argsRaw)])

			// Add the length of the function call to skip those tokens
			*startAt += len(argsRaw)

			// Don't do any further processing and call analyzeStatement
			// for it to parse the function call
			newVal := AnalyzeStatement(
				argsRaw,
				variables,
				functions,
				mods,
				dummyNothingType,
			)

			// Check for type mismatch
			analyzeAssignmentTypes(lastValue, newVal.GetType(), firstToken)

			*lastValue = newVal.GetType()
		}

		// Reassign for the next iteration
		firstToken = statement[*startAt]
	}

	if firstTokenType == token.StringLiteral {
		newVal := wrapper.ForceNewTypeWrapper(
			"str",
			make([]wrapper.TypeWrapper, 0),
			types.StringType,
		)

		// Check for type mismatch
		analyzeAssignmentTypes(lastValue, newVal, firstToken)

		*lastValue = newVal
		return
	}

	if firstTokenType == token.NumLiteral || firstTokenType == token.DecimalLiteral {
		rawType := types.IntType
		typeName := "num"

		if firstTokenType == token.DecimalLiteral {
			rawType = types.DecimalType
			typeName = "dec"
		}

		*expectingArithmeticOperator = true
		newVal := wrapper.ForceNewTypeWrapper(
			typeName,
			make([]wrapper.TypeWrapper, 0),
			rawType,
		)

		// Check for type mismatch
		analyzeAssignmentTypes(lastValue, newVal, firstToken)

		*lastValue = newVal
		return
	}

	logger.TokenError(
		firstToken,
		"Invalid boolean expression",
		"Expected a boolean literal, a variable or an arithmetic expression",
		"Reconstruct this boolean expression to be valid",
		"Remember that creating new objects directly within a boolean expression is not allowed",
		"For example, 'new Object().property > 0' will always fail",
		"In such case, please store the object in a variable and then use it in the boolean expression",
	)
}

// AnalyzeBool analyzes the given boolean expression
func AnalyzeBool(
	statement []token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.FluentMod,
	trace token.Token,
) wrapper.TypeWrapper {
	if len(statement) == 0 {
		logger.TokenError(
			trace,
			"Unexpected token",
			"Expected a boolean expression",
			"Build a boolean expression like 'true' or 'false'",
		)
	}

	lastValue := dummyNothingType

	if len(statement) == 1 {
		analyzeSingleBool(statement[0], variables, functions, mods)
		return lastValue
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

	expectingOperator := false
	// Boolean expressions can also have arithmetic expressions within them
	expectingArithmeticOperator := false

	// Store arithmetic expressions to invoke arithmetic analyzer
	lastArithmeticExpression := make([]token.Token, 0)

	// The remaining statement after the negations
	remainingStatement := statement[startAt:]

	for i, unit := range remainingStatement {
		if i < startAt {
			continue
		}

		if expectingOperator {
			checkBoolOperator(unit, remainingStatement[i:])

			if unit.GetType() != token.Equal {
				// Don't reset the last value for the equality operator
				// as it can be used to compare different types
				lastValue = dummyNothingType
			}

			expectingArithmeticOperator = false
			expectingOperator = false
			continue
		}

		if expectingArithmeticOperator {
			if isBooleanOperator(unit) {
				// Directly analyze the arithmetic expression
				AnalyzeArithmetic(
					lastArithmeticExpression,
					variables,
					functions,
					mods,
				)

				// Reset the slice
				lastArithmeticExpression = make([]token.Token, 0)
				expectingArithmeticOperator = false
				continue
			}

			lastArithmeticExpression = append(lastArithmeticExpression, unit)
			continue
		}

		analyzeBoolExpression(
			remainingStatement[i:],
			variables,
			functions,
			mods,
			&startAt,
			&expectingOperator,
			&expectingArithmeticOperator,
			&lastValue,
		)

		expectingOperator = true
	}

	return lastValue
}
