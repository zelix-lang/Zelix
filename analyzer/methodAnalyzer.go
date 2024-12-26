package analyzer

import (
	"fluent/code"
	"fluent/code/mod"
	"fluent/code/wrapper"
	"fluent/stack"
	"fluent/token"
)

// AnalyzeMethod analyzes a method call
func AnalyzeMethod(
	method *code.Function,
	functions *map[string]map[string]*code.Function,
	mods *map[string]map[string]*mod.FluentMod,
	lastValue *wrapper.FluentObject,
	trace token.Token,
	checkArgs bool,
	args ...wrapper.FluentObject,
) wrapper.FluentObject {
	variables := stack.NewStack()

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
