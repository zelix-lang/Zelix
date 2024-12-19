package analyzer

import (
	"surf/ast"
	"surf/code"
	"surf/core/stack"
	"surf/logger"
	"surf/tokenUtil"
)

// AnalyzeVariableDeclaration analyzes the declaration of a variable
func AnalyzeVariableDeclaration(
	statement []code.Token,
	variables *stack.StaticStack,
	functions *map[string]map[string]*ast.Function,
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
	)

	if value != expectedType {
		logger.TokenError(
			varTypeTokens[0],
			"Type mismatch",
			"The variable type does not match the value type",
			"Change the value type",
		)
	}
}
