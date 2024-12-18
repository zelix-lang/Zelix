package analyzer

import (
	"surf/ast"
	"surf/code"
	"surf/core/engine/args"
	"surf/core/stack"
	"surf/logger"
	"surf/object"
	"surf/tokenUtil"
)

// AnalyzeIdentifier analyzes the given identifier
// and next tokens
func AnalyzeIdentifier(
	statement []code.Token,
	variables *stack.StaticStack,
	functions *map[string]map[string]*ast.Function,
	startAt *int,
	lastValueType *object.SurfObjectType,
	isArithmetic *bool,
	isFunCall *bool,
) {
	firstToken := statement[0]

	// Check if the identifier is a variable
	// or a function call
	variable, varFound := variables.Load(firstToken.GetValue())
	function, funFound, sameFile := ast.LocateFunction(*functions, firstToken.GetFile(), firstToken.GetValue())

	if !varFound && !funFound {
		logger.TokenError(
			firstToken,
			"Undefined reference to variable "+firstToken.GetValue(),
			"The variable "+firstToken.GetValue()+" was not found in the current scope",
			"Declare the variable in the current scope",
		)
	}

	if funFound {
		if !sameFile && !function.IsPublic() {
			logger.TokenError(
				firstToken,
				"Function "+firstToken.GetValue()+" is not public",
				"Move the function to the current file or make it public",
			)
		}

		statementLen := len(statement)
		if statementLen < 3 || statement[1].GetType() != code.OpenParen {
			logger.TokenError(
				firstToken,
				"Invalid function call",
				"A function call must be followed by parentheses",
				"Add parentheses to the function call",
			)
		}

		// Parse and check the arguments
		// Extract all the tokens of the function invocation
		call := tokenUtil.ExtractTokensBefore(
			statement,
			code.CloseParen,
			true,
			code.OpenParen,
			code.CloseParen,
		)

		// ExtractTokensBefore also checks the end closing parenthesis
		// is also met, so no need to check it here

		argumentsRaw, skipped := args.SplitArgs(call)
		arguments := make([]object.SurfObjectType, len(argumentsRaw))

		for i, argument := range argumentsRaw {
			arguments[i] = AnalyzeStatement(
				argument,
				variables,
				functions,
			)
		}

		*startAt = skipped

		*lastValueType = AnalyzeFun(
			function,
			functions,
			firstToken,
			true,
			arguments...,
		)
	} else {
		*lastValueType = variable
		*startAt = 1
	}

	if *lastValueType == object.IntType || *lastValueType == object.DecimalType {
		*isArithmetic = true
	}

	*isFunCall = funFound
}
