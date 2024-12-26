package analyzer

import (
	"fluent/ansi"
	"fluent/ast"
	"fluent/code"
	"fluent/code/mod"
	"fluent/code/types"
	"fluent/code/wrapper"
	"fluent/logger"
	"fluent/stack"
	"fluent/token"
	"fluent/tokenUtil/checker"
	"fluent/tokenUtil/splitter"
	"strconv"
)

// AnalyzeVariableDeclaration analyzes the declaration of a variable
func AnalyzeVariableDeclaration(
	statement []token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.FluentMod,
	constant bool,
) {
	if len(statement) < 5 {
		logger.TokenError(
			statement[0],
			"Invalid variable declaration",
			"A variable declaration must have the form 'var name = value'",
			"Check the variable declaration",
		)
	}

	// The first token is the variable name
	varName := statement[0]

	if !snakeCaseRegex.MatchString(varName.GetValue()) {
		logger.TokenWarning(
			varName,
			"Variable name is not in snake_case",
			"Fluent uses snake_case for variable names",
			"Check "+ansi.Colorize("yellow", "[U-001]")+" in the style guide",
		)
	}

	if varName.GetType() != token.Identifier {
		logger.TokenError(
			varName,
			"Invalid variable name",
			"A variable name must be an identifier",
			"Change the variable name",
		)
	}

	// Check if the variable is already declared
	_, varFound := variables.Load(varName.GetValue())
	fun, funFound, sameFile := ast.LocateFunction(*functions, varName.GetFile(), varName.GetValue())

	if varFound || (funFound && (fun.IsPublic() || sameFile)) {
		logger.TokenError(
			varName,
			"Redefinition of value '"+varName.GetValue()+"'",
			"Change the variable name",
		)
	}

	colon := statement[1]

	if colon.GetType() != token.Colon {
		logger.TokenError(
			colon,
			"Invalid variable declaration",
			"A variable declaration must have the form 'var name : type = value'",
			"Check the variable declaration",
		)
	}

	// Extract the type
	varTypeTokens, _ := splitter.ExtractTokensBefore(
		statement[2:],
		token.Assign,
		false,
		token.Unknown,
		token.Unknown,
		true,
	)

	if len(varTypeTokens) == 0 {
		logger.TokenError(
			colon,
			"Invalid variable declaration",
			"A variable declaration must have the form 'var name : type = value'",
			"Check the variable declaration",
		)
	}

	// Check if the type is valid
	expectedWrapper := wrapper.NewTypeWrapper(varTypeTokens, varTypeTokens[0])
	var expectedType wrapper.FluentObject
	isMod := expectedWrapper.GetType() == types.ModType
	wasTypeInferred := false

	if isMod {
		module, found, sameFile := mod.FindMod(mods, varTypeTokens[0].GetValue(), varTypeTokens[0].GetFile())

		if !found {
			logger.TokenError(
				varTypeTokens[0],
				"Invalid type '"+varTypeTokens[0].GetValue()+"'",
				"The module "+varTypeTokens[0].GetValue()+" was not found in the current scope",
				"Import the module or define it in the current scope",
			)
		}

		if !sameFile && !module.IsPublic() {
			logger.TokenError(
				varTypeTokens[0],
				"Module "+varTypeTokens[0].GetValue()+" is not public",
				"Move the module to the current file or make it public",
			)
		}

		if len(module.GetTemplates()) != len(expectedWrapper.GetParameters()) {
			logger.TokenError(
				varTypeTokens[0],
				"Invalid number of templates",
				"The module "+varTypeTokens[0].GetValue()+" expects "+strconv.Itoa(len(module.GetTemplates()))+" templates",
				"Add the required templates",
			)
		}

		expectedType = wrapper.NewFluentObject(expectedWrapper, module)
		wasTypeInferred = true
	} else {
		if !checker.IsValidType(varTypeTokens[0].GetType()) {
			logger.TokenError(
				varTypeTokens[0],
				"Invalid type '"+varTypeTokens[0].GetValue()+"'",
				"Change the type to a valid one",
			)
		}

		expectedType = wrapper.NewFluentObject(expectedWrapper, nil)
		wasTypeInferred = false
	}

	// Analyze the statement
	AnalyzeType(
		statement[(len(varTypeTokens)+3):],
		variables,
		functions,
		mods,
		expectedType,
		!wasTypeInferred,
	)

	// Put the variable in the stack
	variables.Append(varName.GetValue(), expectedType, constant)
}
