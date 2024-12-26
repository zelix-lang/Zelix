package analyzer

import (
	"fluent/ansi"
	"fluent/code"
	"fluent/code/mod"
	"fluent/code/wrapper"
	"fluent/logger"
	"fluent/token"
	"regexp"
)

var pascalCaseRegex = regexp.MustCompile(`^[A-Z][a-z]+(?:[A-Z][a-z]+)*$`)

// AnalyzeMod analyzes the given mod template
func AnalyzeMod(
	mod mod.FluentMod,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.FluentMod,
) {
	if !pascalCaseRegex.MatchString(mod.GetName()) {
		logger.TokenWarning(
			mod.GetTrace(),
			"Module name is not in PascalCase",
			"Fluent uses PascalCase for modules' names",
			"Check "+ansi.Colorize("yellow", "[U-003]")+" in the style guide",
		)
	}

	if !mod.IsInitialized() {
		variables := mod.GetVariables()
		varDeclarations := mod.GetVarDeclarations()

		// Construct the module's variables stack
		for _, varDecl := range varDeclarations {
			AnalyzeVariableDeclaration(
				varDecl[1:],
				variables,
				functions,
				mods,
				varDecl[0].GetType() == token.Const,
			)
		}
	}

	// Analyze the mod's methods
	for _, method := range mod.GetMethods() {
		dummyObject := wrapper.NewFluentObject(mod.BuildDummyWrapper(), &mod)

		AnalyzeMethod(
			method,
			functions,
			mods,
			&dummyObject,
			mod.GetTrace(),
			false,
		)
	}

}
