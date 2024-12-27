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

	// Extract the type
	expectedWrapper := wrapper.NewTypeWrapper(varTypeTokens, varTypeTokens[0])
	var expectedType wrapper.TypeWrapper
	isMod := expectedWrapper.GetType() == types.ModType

	if !isMod && !checker.IsValidType(varTypeTokens[0].GetType()) {
		logger.TokenError(
			varTypeTokens[0],
			"Invalid type '"+varTypeTokens[0].GetValue()+"'",
			"Change the type to a valid one",
		)
	}

	expectedType = wrapper.NewTypeWrapper(varTypeTokens, varTypeTokens[0])
	AnalyzeGeneric(expectedType, mods, varTypeTokens[0])

	// +2 for the var name + the colon
	// +1 for the equals sign
	valueTokens := statement[3+len(varTypeTokens):]
	// Interpret the statement to get a value
	value := AnalyzeStatement(
		valueTokens,
		variables,
		functions,
		mods,
		expectedType,
	)

	variables.Append(varName.GetValue(), value, constant)
}
