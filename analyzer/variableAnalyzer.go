package analyzer

import (
	"zyro/ansi"
	"zyro/ast"
	"zyro/code"
	"zyro/code/mod"
	"zyro/code/types"
	"zyro/code/wrapper"
	"zyro/core/stack"
	"zyro/logger"
	"zyro/token"
	"zyro/tokenUtil/checker"
	"zyro/tokenUtil/splitter"
)

// AnalyzeVariableDeclaration analyzes the declaration of a variable
func AnalyzeVariableDeclaration(
	statement []token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.ZyroMod,
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
			"Zyro uses snake_case for variable names",
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
	_, funFound, _ := ast.LocateFunction(*functions, varName.GetFile(), varName.GetValue())

	if varFound || funFound {
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
	varTypeTokens := splitter.ExtractTokensBefore(
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
	expectedWrapper := wrapper.NewTypeWrapper(varTypeTokens, varTypeTokens[0], false)
	var expectedType wrapper.ZyroObject
	isMod := expectedWrapper.GetType() == types.ModType

	if isMod {
		module, found, sameFile := mod.FindMod(mods, varTypeTokens[0].GetValue(), varTypeTokens[0].GetFile())

		if !sameFile && !module.IsPublic() {
			logger.TokenError(
				varTypeTokens[0],
				"Module "+varTypeTokens[0].GetValue()+" is not public",
				"Move the module to the current file or make it public",
			)
		}

		if !found {
			logger.TokenError(
				varTypeTokens[0],
				"Invalid type '"+varTypeTokens[0].GetValue()+"'",
				"Change the type to a valid one",
			)
		}

		expectedType = wrapper.NewZyroObject(module.BuildDummyWrapper(), module)
	} else {
		if !checker.IsValidType(varTypeTokens[0].GetType()) {
			logger.TokenError(
				varTypeTokens[0],
				"Invalid type '"+varTypeTokens[0].GetValue()+"'",
				"Change the type to a valid one",
			)
		}

		expectedType = wrapper.NewZyroObject(expectedWrapper, nil)
	}

	// Analyze the statement
	AnalyzeType(
		statement[(len(varTypeTokens)+3):],
		variables,
		functions,
		mods,
		expectedType,
	)

	// Put the variable in the stack
	variables.Append(varName.GetValue(), expectedType, constant)
}
