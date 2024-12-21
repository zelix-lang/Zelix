package analyzer

import (
	"strconv"
	"surf/code"
	"surf/core/stack"
	"surf/logger"
	"surf/object"
	"surf/token"
	"surf/tokenUtil"
	"surf/util"
	"time"
)

// checkParamType checks if the given parameter type is valid
func checkParamType(
	paramType object.SurfObject,
	trace token.Token,
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
	function *code.Function,
	functions *map[string]map[string]*code.Function,
	mods *map[string]*code.SurfMod,
	trace token.Token,
	checkArgs bool,
	variables *stack.Stack,
	args ...object.SurfObject,
) object.SurfObject {
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
			expected := tokenUtil.FromRawType(actualParams[param][0], mods)
			value := args[i]

			checkParamType(expected, trace)
			if value.GetType() != expected.GetType() {
				logger.TokenError(
					trace,
					"Mismatched parameter types",
					"This function did not expect this parameter this time",
					"Change the parameters of the function call",
				)
			}

			variables.Append(param, value, false)
		}
	} else {
		// Store the parameters without checking to avoid undefined references
		argsKeys := util.MapKeys(actualParams)

		for _, param := range argsKeys {
			expected := tokenUtil.FromRawType(actualParams[param][0], mods)
			checkParamType(expected, trace)
			variables.Append(param, expected, false)
		}
	}

	// Beyond this point, standard functions no longer
	// need to be evaluated
	if function.IsStd() {
		return object.NewSurfObject(object.NothingType, nil)
	}

	// Used to skip tokens
	skipToIndex := 0

	for i, unit := range function.GetBody() {
		if i < skipToIndex {
			continue
		}

		tokenType := unit.GetType()

		if tokenType == token.Identifier || tokenType == token.Let || tokenType == token.Const || tokenType == token.New {
			// Extract the statement
			statement := tokenUtil.ExtractTokensBefore(
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
				AnalyzeVariableDeclaration(statement[1:], variables, functions, mods, tokenType == token.Const)
				continue
			}

			AnalyzeStatement(statement, variables, functions, mods)
			continue
		}

		logger.TokenError(
			unit,
			"Unexpected token",
			"Expected an identifier or a statement",
			"Check the function body",
		)
	}

	// TODO! Parse return statements

	// Destroy the scope
	variables.DestroyScope()
	return object.NewSurfObject(object.NothingType, nil)
}
