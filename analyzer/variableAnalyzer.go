package analyzer

import (
	"regexp"
	"surf/ansi"
	"surf/ast"
	"surf/code"
	"surf/core/stack"
	"surf/logger"
	"surf/tokenUtil"
)

// A regex to match camelCase variable names
var snakeCaseRegex, _ = regexp.Compile("^[a-z]+(_[a-z0-9]+)*$")

// AnalyzeVariableDeclaration analyzes the declaration of a variable
func AnalyzeVariableDeclaration(
	statement []code.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*ast.Function,
	mods *map[string]*ast.SurfMod,
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
			"Surf uses snake_case for variable names",
			"Check "+ansi.Colorize("yellow", "[U-001]")+" in the style guide",
		)
	}

	if varName.GetType() != code.Identifier {
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

	if colon.GetType() != code.Colon {
		logger.TokenError(
			colon,
			"Invalid variable declaration",
			"A variable declaration must have the form 'var name : type = value'",
			"Check the variable declaration",
		)
	}

	// Extract the type
	varTypeTokens := tokenUtil.ExtractTokensBefore(
		statement[2:],
		code.Assign,
		false,
		code.Unknown,
		code.Unknown,
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
	expectedType := tokenUtil.FromRawType(varTypeTokens[0], variables)
	if !tokenUtil.IsValidType(varTypeTokens[0].GetType()) {
		logger.TokenError(
			varTypeTokens[0],
			"Invalid type '"+varTypeTokens[0].GetValue()+"'",
			"Change the type to a valid one",
		)
	}

	// Analyze the statement
	value := AnalyzeStatement(
		statement[(len(varTypeTokens)+3):],
		variables,
		functions,
		mods,
	)

	if value.GetType() != expectedType.GetType() {
		logger.TokenError(
			varTypeTokens[0],
			"Type mismatch",
			"The variable type does not match the value type",
			"Change the value type",
		)
	}

	// Put the variable in the stack
	variables.Append(varName.GetValue(), expectedType)
}
