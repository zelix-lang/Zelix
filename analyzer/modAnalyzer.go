package analyzer

import (
	"fluent/ansi"
	"fluent/code"
	"fluent/code/mod"
	"fluent/code/wrapper"
	"fluent/logger"
	"fluent/stack"
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

	dummyStack := stack.NewStack()
	// Check for variables redeclaration
	for _, tokens := range mod.GetVarDeclarations() {
		AnalyzeVariableDeclaration(tokens[1:], dummyStack, functions, mods, false)
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
