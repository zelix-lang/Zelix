package analyzer

import (
	"zyro/code"
	"zyro/code/mod"
	"zyro/code/wrapper"
	"zyro/stack"
	"zyro/token"
)

// AnalyzeMethod analyzes a method call
func AnalyzeMethod(
	method *code.Function,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.ZyroMod,
	lastValue *wrapper.ZyroObject,
	trace token.Token,
	checkArgs bool,
	inferToType wrapper.TypeWrapper,
	args ...wrapper.ZyroObject,
) wrapper.ZyroObject {
	variables := stack.NewStack()
	currentMod := lastValue.GetValue().(*mod.ZyroMod)
	varTemplates := currentMod.GetVarDeclarations()

	for _, template := range varTemplates {
		AnalyzeVariableDeclaration(template[1:], variables, functions, mods, template[0].GetType() == token.Const)
	}
	// Add "this" to the stack
	variables.CreateScope()
	variables.Append("this", *lastValue, true)

	// Analyze the result
	result := AnalyzeFun(
		method,
		functions,
		mods,
		trace,
		checkArgs,
		variables,
		args...,
	)

	// Remove "this" from the stack
	variables.DestroyScope()

	return result
}
