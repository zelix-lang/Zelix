package analyzer

import (
	"surf/code"
	"surf/core/stack"
	"surf/object"
	"surf/token"
)

func AnalyzePropAccess(
	prop []token.Token,
	variables *stack.Stack,
	functions *map[string]map[string]*code.Function,
	mods *map[string]*code.SurfMod,
	lastValue *object.SurfObject,
) {

}
