package analyzer

import (
	"regexp"
	"surf/ansi"
	"surf/code"
	"surf/logger"
	"surf/object"
)

var pascalCaseRegex = regexp.MustCompile(`^[A-Z][a-z]+(?:[A-Z][a-z]+)*$`)

// AnalyzeMod analyzes the given mod template
func AnalyzeMod(
	mod code.SurfMod,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*code.SurfMod,
) {
	if !pascalCaseRegex.MatchString(mod.GetName()) {
		logger.TokenWarning(
			mod.GetTrace(),
			"Module name is not in PascalCase",
			"Surf uses PascalCase for modules' names",
			"Check "+ansi.Colorize("yellow", "[U-003]")+" in the style guide",
		)
	}

	// Analyze the mod's methods
	for _, method := range mod.GetMethods() {
		dummyObject := object.NewSurfObject(object.ModType, &mod)

		AnalyzeMethod(
			*method,
			functions,
			mods,
			&dummyObject,
			mod.GetTrace(),
			false,
		)
	}

}
