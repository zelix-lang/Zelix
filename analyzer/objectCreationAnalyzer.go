package analyzer

import (
	"fluent/code"
	"fluent/code/mod"
	"fluent/code/wrapper"
	"fluent/logger"
	"fluent/stack"
	"fluent/token"
	"fluent/tokenUtil/splitter"
)

func AnalyzeObjectCreation(
	statement []token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.FluentMod,
	startAt *int,
	lastValue *wrapper.FluentObject,
	inferToType wrapper.TypeWrapper,
) {
	finalType := inferToType

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
			templatesRaw, _ := splitter.ExtractTokensBefore(
				statement[1:],
				token.OpenParen,
				true,
				token.LessThan,
				token.GreaterThan,
				true,
			)

			// len(templatesRaw) == 0 is impossible at this point
			finalType = wrapper.NewTypeWrapper(
				templatesRaw,
				templatesRaw[0],
			)

			AnalyzeGeneric(finalType, mods, templatesRaw[0])
			lookForParenAt = len(templatesRaw) + 1
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

	generics := make(map[string]wrapper.TypeWrapper)
	builtInMods := make(map[string]*mod.FluentMod)

	for i, template := range module.GetTemplates() {
		param := finalType.GetParameters()[i]
		generics[template.GetBaseType()] = param
		// Add dummy generics to the mods so things like:
		// let a: T = abc;
		// don't throw an error
		newMod := mod.NewFluentMod(
			make(map[string]*wrapper.FluentObject),
			make(map[string]*code.Function),
			make(map[string]*code.Function),
			template.GetBaseType(),
			modName.GetFile(),
			make([][]token.Token, 0),
			true,
			modName,
			[]wrapper.TypeWrapper{},
		)

		builtInMods[template.GetBaseType()] = &newMod
	}

	(*mods)["built-in"] = builtInMods
	withoutGenerics := module.BuildWithoutGenerics(generics)

	*lastValue = wrapper.NewFluentObject(
		withoutGenerics.BuildDummyWrapper(),
		&withoutGenerics,
	)

	// Check if the module has any constructor
	constructor, constructorFound, constructorPublic := withoutGenerics.GetMethod(modName.GetValue())

	// Parse the arguments
	argsRange := statement[(lookForParenAt + 1) : len(statement)-1]
	argsRaw := splitter.SplitTokens(
		argsRange,
		token.Comma,
		token.OpenParen,
		token.CloseParen,
	)

	if !constructorFound {
		if len(argsRaw) > 0 {
			logger.TokenError(
				modName,
				"Constructor "+modName.GetValue()+" not found",
				"The constructor "+modName.GetValue()+" does not exist",
				"Check the constructor name",
			)
		}

		*startAt += lookForParenAt + 2
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

	*startAt += 2 + lookForParenAt
	for i, arg := range argsRaw {
		*startAt += len(arg)
		if i < len(argsRaw)-1 {
			// +1 for the comma
			*startAt += 1
		}
	}

	args := make([]wrapper.FluentObject, len(argsRaw))
	for i, arg := range argsRaw {
		args[i] = AnalyzeStatement(
			arg,
			variables,
			functions,
			mods,
			dummyNothingType,
		)
	}

	if !inferToType.Compare(dummyNothingType) && !finalType.Compare(inferToType) {
		logger.TokenError(
			modName,
			"Invalid object creation",
			"The object creation does not match the expected type",
			"Check the object creation",
		)
	}

	AnalyzeMethod(
		constructor,
		functions,
		mods,
		lastValue,
		modName,
		true,
		finalType,
		args...,
	)

	// Delete the built-in mods
	delete(*mods, "built-in")

}
