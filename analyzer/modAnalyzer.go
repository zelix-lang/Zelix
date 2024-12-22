package analyzer

import (
	"regexp"
	"zyro/ansi"
	"zyro/code"
	"zyro/code/mod"
	"zyro/code/wrapper"
	"zyro/logger"
)

var pascalCaseRegex = regexp.MustCompile(`^[A-Z][a-z]+(?:[A-Z][a-z]+)*$`)

// AnalyzeMod analyzes the given mod template
func AnalyzeMod(
	mod mod.ZyroMod,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.ZyroMod,
) {
	if !pascalCaseRegex.MatchString(mod.GetName()) {
		logger.TokenWarning(
			mod.GetTrace(),
			"Module name is not in PascalCase",
			"Zyro uses PascalCase for modules' names",
			"Check "+ansi.Colorize("yellow", "[U-003]")+" in the style guide",
		)
	}

	// Analyze the mod's methods
	for _, method := range mod.GetMethods() {
		dummyObject := wrapper.NewZyroObject(mod.BuildDummyWrapper(), &mod)

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
