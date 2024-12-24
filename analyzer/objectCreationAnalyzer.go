package analyzer

import (
	"zyro/code"
	"zyro/code/mod"
	"zyro/code/wrapper"
	"zyro/logger"
	"zyro/stack"
	"zyro/token"
	"zyro/tokenUtil/splitter"
)

func AnalyzeObjectCreation(
	statement []token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.ZyroMod,
	startAt *int,
	lastValue *wrapper.ZyroObject,
	inferToType wrapper.TypeWrapper,
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
	module, modFound, sameFile := mod.FindMod(mods, modName.GetValue(), modName.GetFile())

	if !modFound {
		logger.TokenError(
			modName,
			"Undefined reference to module "+modName.GetValue(),
			"The module "+modName.GetValue()+" was not found in the current scope",
			"Import the module in the current scope",
		)
	}

	// Check access to the module
	if !module.IsPublic() && !sameFile {
		logger.TokenError(
			modName,
			"Module "+modName.GetValue()+" is not public",
			"Move the module to the current file or make it public",
		)
	}

	lookForParenAt := 2
	if len(module.GetTemplates()) > 0 {
		// At least: new MyObject<>() -> 6 tokens
		if len(statement) < 6 || statement[2].GetType() != token.LessThan {
			logger.TokenError(
				statement[2],
				"Invalid object creation",
				"An object creation with templates must be followed by a less than sign",
				"Add templates like: new MyObject<template1, template2>()",
			)
		}

		if statement[3].GetType() == token.GreaterThan {
			if inferToType.Compare(dummyNothingType) {
				// Constructor was called without templates
				logger.TokenError(
					statement[3],
					"Cannot infer type",
					"You must specify the templates for the object",
					"Add templates like: new MyObject<template1, template2>()",
				)
			}

			lookForParenAt = 4
		} else {
			// Extract the templates
			templatesRaw := splitter.ExtractTokensBefore(
				statement[2:],
				token.GreaterThan,
				true,
				token.LessThan,
				token.GreaterThan,
				true,
			)

			// len(templatesRaw) == 0 is impossible at this point

			inferToType = wrapper.NewTypeWrapper(
				templatesRaw,
				templatesRaw[0],
			)

			lookForParenAt = 4 + len(templatesRaw)
		}

	}

	// Validate the parentheses
	if statement[lookForParenAt].GetType() != token.OpenParen {
		logger.TokenError(
			statement[lookForParenAt],
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

	*lastValue = wrapper.NewZyroObject(
		module.BuildDummyWrapper(),
		module,
	)

	// Check if the module has any constructor
	constructor, constructorFound, constructorPublic := module.GetMethod(modName.GetValue())
	if !constructorFound {
		*startAt += 4
		// No constructor found, return the module
		return
	}

	// Check if the constructor is public
	if !constructorPublic && module.GetFile() != modName.GetFile() {
		logger.TokenError(
			modName,
			"Constructor "+modName.GetValue()+" is not public",
			"Move the constructor to the current file or make it public",
		)
	}

	// Parse the arguments
	argsRange := statement[(lookForParenAt + 1) : len(statement)-1]
	argsRaw := splitter.SplitTokens(
		argsRange,
		token.Comma,
		token.OpenParen,
		token.CloseParen,
	)

	*startAt += len(argsRaw) + 2 + lookForParenAt
	args := make([]wrapper.ZyroObject, len(argsRaw))
	for i, arg := range argsRaw {
		args[i] = AnalyzeStatement(
			arg,
			variables,
			functions,
			mods,
			dummyNothingType,
		)
	}

	AnalyzeMethod(
		constructor,
		functions,
		mods,
		lastValue,
		modName,
		true,
		args...,
	)

}
