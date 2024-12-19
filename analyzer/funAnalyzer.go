package analyzer

import (
	"strconv"
	"surf/ast"
	"surf/code"
	"surf/core/stack"
	"surf/logger"
	"surf/object"
	"surf/tokenUtil"
	"surf/util"
	"time"
)

// checkParamType checks if the given parameter type is valid
func checkParamType(
	paramType object.SurfObject,
	trace code.Token,
) {
	if paramType.GetType() == object.NothingType {
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
	function *ast.Function,
	functions *map[string]map[string]*ast.Function,
	trace code.Token,
	checkArgs bool,
	args ...object.SurfObjectType,
) object.SurfObjectType {
	function.SetTimesCalled(function.GetTimesCalled() + 1)
	function.SetLastCalled(time.Now())

	// Check for stack overflows
	if function.GetTimesCalled() > 1000 && time.Since(function.GetLastCalled()).Seconds() < 1 {
		logger.TokenError(
			trace,
			"Stack overflow",
			"This function has overflown its stack",
			"Check for infinite loops",
		)
	}

	// Create a new StaticStack
	variables := stack.NewStack()
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
		argsKeys := util.MapKeys(actualParams)

		for i, param := range argsKeys {
			expected := tokenUtil.FromRawType(actualParams[param][0], variables)
			value := args[i]

			checkParamType(expected, trace)
			if value != expected.GetType() {
				logger.TokenError(
					trace,
					"Mismatched parameter types",
					"This function did not expect this parameter this time",
					"Change the parameters of the function call",
				)
			}

			variables.Append(param, object.NewSurfObject(value, nil))
		}
	} else {
		// Store the parameters without checking to avoid undefined references
		argsKeys := util.MapKeys(actualParams)

		for _, param := range argsKeys {
			expected := tokenUtil.FromRawType(actualParams[param][0], variables)
			checkParamType(expected, trace)
			variables.Append(param, expected)
		}
	}

	// Beyond this point, standard functions no longer
	// need to be evaluated
	if function.IsStd() {
		return object.NothingType
	}

	// Used to skip tokens
	skipToIndex := 0

	for i, token := range function.GetBody() {
		if i < skipToIndex {
			continue
		}

		tokenType := token.GetType()

		if tokenType == code.Identifier || tokenType == code.Let {
			// Extract the statement
			statement := tokenUtil.ExtractTokensBefore(
				function.GetBody()[i:],
				code.Semicolon,
				// Don't handle nested statements here
				false,
				code.Unknown,
				code.Unknown,
			)

			skipToIndex = i + len(statement) + 1

			// Analyze the statement
			if tokenType == code.Let {
				AnalyzeVariableDeclaration(statement[1:], variables, functions)
				continue
			}

			AnalyzeStatement(statement, variables, functions)
			continue
		}

		logger.TokenError(
			token,
			"Unexpected token",
			"Expected an identifier or a statement",
			"Check the function body",
		)
	}

	// TODO! Parse return statements
	return object.NothingType
}
