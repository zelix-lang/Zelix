package analyzer

import (
	"surf/code"
	"surf/core/stack"
	"surf/logger"
	"surf/object"
	"surf/token"
	"surf/tokenUtil"
)

func AnalyzeObjectCreation(
	statement []token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*code.SurfMod,
	startAt *int,
	lastValue *object.SurfObject,
) {
	// The statement should have at least 4 tokens:
	// new MyObject()
	if len(statement) < 4 {
		logger.TokenError(
			statement[0],
			"Invalid object creation",
			"An object creation must be followed by an identifier and parentheses",
			"Check the object creation",
		)
	}

	// At this point, the first token is always "new"
	// no need to check it
	modName := statement[1]
	mod, modFound, sameFile := code.FindMod(mods, modName.GetValue(), modName.GetFile())

	if !modFound {
		logger.TokenError(
			modName,
			"Undefined reference to module "+modName.GetValue(),
			"The module "+modName.GetValue()+" was not found in the current scope",
			"Import the module in the current scope",
		)
	}

	// Check access to the module
	if !mod.IsPublic() && !sameFile {
		logger.TokenError(
			modName,
			"Module "+modName.GetValue()+" is not public",
			"Move the module to the current file or make it public",
		)
	}

	// Validate the parentheses
	if statement[2].GetType() != token.OpenParen {
		logger.TokenError(
			statement[2],
			"Invalid object creation",
			"An object creation must be followed by parentheses",
			"Check the object creation",
		)
	}

	if statement[len(statement)-1].GetType() != token.CloseParen {
		logger.TokenError(
			statement[len(statement)-1],
			"Invalid object creation",
			"An object creation must end with a closing parenthesis",
			"Check the object creation",
		)
	}

	*lastValue = object.NewSurfObject(
		object.ModType,
		mod,
	)

	// Check if the module has any constructor
	constructor, constructorFound, constructorPublic := mod.GetMethod(modName.GetValue())
	if !constructorFound {
		*startAt += 4
		// No constructor found, return the module
		return
	}

	// Check if the constructor is public
	if !constructorPublic && mod.GetFile() != modName.GetFile() {
		logger.TokenError(
			modName,
			"Constructor "+modName.GetValue()+" is not public",
			"Move the constructor to the current file or make it public",
		)
	}

	// Parse the arguments
	argsRange := statement[3 : len(statement)-1]
	argsRaw := tokenUtil.SplitTokens(
		argsRange,
		token.Comma,
		token.OpenParen,
		token.CloseParen,
	)

	*startAt += len(argsRaw) + 4
	args := make([]object.SurfObject, len(argsRaw))
	for i, arg := range argsRaw {
		args[i] = AnalyzeStatement(
			arg,
			variables,
			functions,
			mods,
		)
	}

	AnalyzeMethod(
		*constructor,
		functions,
		mods,
		lastValue,
		modName,
		args...,
	)

}
