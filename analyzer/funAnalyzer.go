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
	"strconv"
)

// Store the functions that have already been checked
// Use a map for O(1) lookup
var calledFunctions = map[string]struct{}{}

// checkParamType checks if the given parameter type is valid
func checkParamType(
	paramType wrapper.TypeWrapper,
	trace token.Token,
) {
	if paramType.GetType() == types.NothingType {
		logger.TokenError(
			trace,
			"Invalid parameter type",
			"Parameters cannot be of type 'nothing'",
			"Check the function definition",
		)
	}
}

// AnalyzeFun analyzes the given function
func AnalyzeFun(
	function *code.Function,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.FluentMod,
	trace token.Token,
	checkArgs bool,
	variables *stack.Stack,
	args ...wrapper.FluentObject,
) wrapper.FluentObject {
	// Still check params and return type
	AnalyzeGeneric(function.GetReturnType(), mods, trace, true)

	// Create a new StaticStack
	actualParams := function.GetParameters()

	if checkArgs {
		if len(args) != len(function.GetParameters()) {
			logger.TokenError(
				trace,
				"Invalid number of arguments",
				"This function expects "+strconv.Itoa(len(function.GetParameters()))+" arguments",
				"Add the missing arguments",
			)
		}

		// Store the arguments in the variables
		for i, param := range actualParams {
			expected := param.GetType()
			value := args[i]

			checkParamType(expected, trace)
			if !expected.Compare(value.GetType()) {
				valueType := value.GetType()
				logger.TokenError(
					trace,
					"Mismatched parameter types",
					"This function did not expect this parameter this time",
					"Change the parameters of the function call",
					"Expected: "+expected.Marshal(),
					"Got: "+valueType.Marshal(),
				)
			}

			variables.Append(param.GetName(), value, false)
		}
	} else {
		// Store the parameters without checking to avoid undefined references
		for _, param := range actualParams {
			expected := param.GetType()
			checkParamType(expected, trace)
			dummyObj := wrapper.NewFluentObject(expected, nil)

			variables.Append(param.GetName(), dummyObj, false)
		}
	}

	functionReturnType := function.GetReturnType()
	returnValue := wrapper.NewFluentObject(
		functionReturnType,
		nil,
	)

	// Check if the function has already been called
	if _, ok := calledFunctions[function.GetName()]; ok {
		// Stack overflow is caught at runtime, skip this check
		// as it may target a recursive function incorrectly
		return returnValue
	}

	// Mark the function as called
	calledFunctions[function.GetName()] = struct{}{}

	// Used to skip tokens
	skipToIndex := 0

	//  Used to know if we're currently parsing a conditional statement
	inConditional := false
	conditionalCount := 0

	// Used to keep track of function's validity
	elseIfAllowed := false
	hasReturn := false

	for i, unit := range function.GetBody() {
		if i < skipToIndex {
			continue
		}

		tokenType := unit.GetType()

		// If and while share the same logic, so we can combine them
		if tokenType == token.If || tokenType == token.ElseIf || tokenType == token.Else || tokenType == token.While {
			if tokenType != token.While && tokenType != token.Else {
				inConditional = true
				conditionalCount++
			}

			// Extract the bool declaration
			declaration, _ := splitter.ExtractTokensBefore(
				function.GetBody()[i:],
				token.OpenCurly,
				false,
				token.Unknown,
				token.Unknown,
				true,
			)

			if tokenType == token.Else {
				if !elseIfAllowed {
					logger.TokenError(
						declaration[0],
						"Unexpected conditional",
						"An 'else if' statement must be preceded by an 'if' or 'else if' statement",
						"Check the conditional statements",
					)
				}

				elseIfAllowed = false
				AnalyzeElse(declaration[1:])
				// Skip to the end of the if statement
				skipToIndex = i + len(declaration) + 1

				continue
			}

			AnalyzeBool(declaration[1:], variables, functions, mods, unit)
			variables.CreateScope()

			// Skip to the end of the if statement
			skipToIndex = i + len(declaration) + 1
			continue
		} else if tokenType == token.For {
			elseIfAllowed = false

			// Extract the loop declaration
			declaration, _ := splitter.ExtractTokensBefore(
				function.GetBody()[i:],
				token.OpenCurly,
				false,
				token.Unknown,
				token.Unknown,
				true,
			)

			AnalyzeFor(declaration[1:], functions, mods, variables)
			skipToIndex = i + len(declaration) + 1
			continue
		} else if tokenType == token.Identifier || tokenType == token.Let || tokenType == token.Const || tokenType == token.New {
			elseIfAllowed = false

			// Extract the statement
			statement, _ := splitter.ExtractTokensBefore(
				function.GetBody()[i:],
				token.Semicolon,
				// Don't handle nested statements here
				false,
				token.Unknown,
				token.Unknown,
				true,
			)

			skipToIndex = i + len(statement) + 1

			// Analyze the statement
			if tokenType == token.Let || tokenType == token.Const {
				AnalyzeVariableDeclaration(
					statement[1:],
					variables,
					functions,
					mods,
					tokenType == token.Const,
					statement[0],
				)
				continue
			}

			AnalyzeStatement(statement, variables, functions, mods, dummyNothingType)
			continue
		} else if tokenType == token.CloseCurly {
			variables.DestroyScope(unit, function.IsInMod())

			if inConditional {
				conditionalCount--
				inConditional = conditionalCount > 0

				// See if there is a next token and if it is an else if
				if i+1 < len(function.GetBody()) {
					nextToken := function.GetBody()[i+1]
					elseIfAllowed = nextToken.GetType() == token.ElseIf || nextToken.GetType() == token.Else
				}
			}

			continue
		} else if tokenType == token.Return {
			// Extract the statement
			statement, _ := splitter.ExtractTokensBefore(
				function.GetBody()[i:],
				token.Semicolon,
				// Don't handle nested statements here
				false,
				token.Unknown,
				token.Unknown,
				true,
			)

			skipToIndex = i + len(statement) + 1

			// Analyze the statement
			AnalyzeReturn(statement[1:], variables, functions, mods, returnValue)
			if !hasReturn && !inConditional {
				hasReturn = true
			}

			continue
		}

		logger.TokenError(
			unit,
			"Unexpected token",
			"Expected an identifier or a statement",
			"Check the function body",
		)
	}

	if !hasReturn && functionReturnType.GetType() != types.NothingType {
		logger.TokenError(
			trace,
			"Missing return statement",
			"This function must return a value",
			"Add a return statement at the end of the function",
		)
	}

	// Destroy the scope
	variables.DestroyScope(trace, function.IsInMod())

	return returnValue
}
