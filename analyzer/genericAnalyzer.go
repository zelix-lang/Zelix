package analyzer

import (
	"fluent/code/mod"
	"fluent/code/types"
	"fluent/code/wrapper"
	"fluent/logger"
	"fluent/token"
)

// AnalyzeGenericInitialization analyzes the given generic initialization
func AnalyzeGenericInitialization(
	generics []wrapper.TypeWrapper,
	mods *map[string]map[string]*mod.FluentMod,
	trace token.Token,
) {
	for _, param := range generics {
		baseType := param.GetType()

		if baseType != types.ModType {
			logger.TokenError(
				trace,
				"Invalid type",
				"Generic Parameters cannot have",
			)
		}

		baseName := param.GetBaseType()
		module, found, sameFile := mod.FindMod(mods, baseName, trace.GetFile())
		if found && (module.IsPublic() || sameFile) {
			logger.TokenError(
				trace,
				"Invalid type",
				"Generic parameters cannot be modules",
			)
		}
	}
}

// AnalyzeGeneric analyzes the given generic templates
func AnalyzeGeneric(
	genericTemplate wrapper.TypeWrapper,
	mods *map[string]map[string]*mod.FluentMod,
	trace token.Token,
	isReturnType bool,
) {
	baseType := genericTemplate.GetType()
	params := genericTemplate.GetParameters()

	if baseType == types.NothingType && !isReturnType {
		logger.TokenError(
			trace,
			"Invalid type",
			"Generic parameters cannot be of type 'nothing'",
		)
	} else if baseType == types.ModType {
		// len(params) == 0 is always false, skip that condition
		// and directly check the types
		baseName := genericTemplate.GetBaseType()
		module, found, sameFile := mod.FindMod(mods, baseName, trace.GetFile())

		if !found {
			logger.TokenError(
				trace,
				"Undefined reference to module "+baseName,
				"The module "+baseName+" was not found in the current scope",
				"Import the module in the current scope",
			)
		}

		if !sameFile && !module.IsPublic() {
			logger.TokenError(
				trace,
				"Module "+baseName+" is not public",
				"Move the module to the current file or make it public",
			)
		}

		if len(module.GetTemplates()) != len(params) {
			logger.TokenError(
				trace,
				"Invalid number of templates",
				"The number of generic parameters does not match the number of templates",
				"Check the number of generic parameters",
			)
		}

		for _, param := range params {
			AnalyzeGeneric(param, mods, trace, isReturnType)
		}
	} else if len(params) > 0 {
		logger.TokenError(
			trace,
			"Invalid type",
			"Primitive types cannot have generic parameters",
			"Remove the generic parameters",
			"Create a module to hold the generic parameters",
		)
	}
}
