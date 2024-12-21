package analyzer

import (
	"surf/code"
	"surf/core/stack"
	"surf/object"
	"surf/token"
)

// AnalyzeMethod analyzes a method call
func AnalyzeMethod(
	method code.Function,
	functions *map[string]map[string]*code.Function,
	mods *map[string]*code.SurfMod,
	lastValue *object.SurfObject,
	trace token.Token,
	args ...object.SurfObject,
) object.SurfObject {
	variables := stack.NewStack()
	currentMod := lastValue.GetValue().(*code.SurfMod)
	varTemplates := currentMod.GetVarDeclarations()

	for _, template := range varTemplates {
		AnalyzeVariableDeclaration(template[1:], variables, functions, mods, template[0].GetType() == token.Const)
	}
	// Add "this" to the stack
	variables.CreateScope()
	variables.Append("this", *lastValue, true)

	// Analyze the result
	result := AnalyzeFun(
		&method,
		functions,
		mods,
		trace,
		true,
		variables,
		args...,
	)

	// Remove "this" from the stack
	variables.DestroyScope()

	return result
}
