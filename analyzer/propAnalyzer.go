package analyzer

import (
	"fluent/code"
	"fluent/code/mod"
	"fluent/code/types"
	"fluent/code/wrapper"
	"fluent/logger"
	"fluent/stack"
	"fluent/token"
	"fluent/tokenUtil/splitter"
)

// AnalyzePropAccess analyzes the given property access
func AnalyzePropAccess(
	prop []token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.FluentMod,
	lastValue *wrapper.FluentObject,
	isFunCall *bool,
	isAssignment bool,
) {
	// No need to check for prop's length
	// the token splitter ensures that
	// there parts are not empty

	// Check the last value is a mod
	valueTypeWrapper := (*lastValue).GetType()
	if valueTypeWrapper.GetType() != types.ModType {
		logger.TokenError(
			prop[0],
			"Illegal property access",
			"Cannot access properties of a non-object",
			"Check the object type",
		)
	}

	module := lastValue.GetValue().(*mod.FluentMod)
	propName := prop[0]

	// Only identifiers are allowed as property names
	if propName.GetType() != token.Identifier {
		logger.TokenError(
			propName,
			"Illegal property name",
			"Property names must be identifiers",
			"Check the property name",
		)
	}

	// Check for single property access
	// i.e.: object.property
	if len(prop) == 1 {
		val, found := variables.Load(propName.GetValue())

		if !found {
			logger.TokenError(
				prop[0],
				"Property not found",
				"Check the property name",
			)
		}

		// Properties are private by default
		// so we have to check access here
		if module.GetFile() != propName.GetFile() {
			logger.TokenError(
				propName,
				"Illegal access",
				"Cannot access private properties",
				"Use setters and getters to access or modify the property",
			)
		}

		if isAssignment && val.IsConstant() {
			logger.TokenError(
				prop[0],
				"Cannot assign to constant",
				"This property is constant",
				"Remove the assignment or change the property's constant status",
			)
		}

		*lastValue = val.GetValue()
		return
	}

	// Check for constructor call
	if propName.GetValue() == module.GetName() {
		logger.TokenError(
			propName,
			"Cannot access the constructor upon initialization",
			"Define other methods to modularize your code",
		)
	}

	// Find the method
	method, found, public := module.GetMethod(propName.GetValue())

	if !found {
		logger.TokenError(
			propName,
			"Undefined reference to method "+propName.GetValue(),
			"Check for typos or define the method",
		)
	}

	// Check access
	if !public && module.GetFile() != propName.GetFile() {
		logger.TokenError(
			propName,
			"Illegal access",
			"Cannot access private methods",
			"Use public methods to access the method or change the method's visibility",
		)
	}

	// Method call
	// i.e.: object.method()
	if prop[1].GetType() != token.OpenParen || prop[len(prop)-1].GetType() != token.CloseParen {
		logger.TokenError(
			prop[1],
			"Invalid operation",
			"Invalid operation after identifier",
			"A function or method call was expected here",
		)
	}

	argsRaw := prop[:len(prop)-1]
	argsSplit, _ := splitter.SplitArgs(argsRaw)
	funArgs := make([]wrapper.FluentObject, len(argsSplit))

	for i, arg := range argsSplit {
		funArgs[i] = AnalyzeStatement(
			arg,
			variables,
			functions,
			mods,
			dummyNothingType,
		)
	}

	// Update metadata
	*isFunCall = true

	AnalyzeMethod(
		method,
		functions,
		mods,
		lastValue,
		propName,
		true,
		funArgs...,
	)
}
