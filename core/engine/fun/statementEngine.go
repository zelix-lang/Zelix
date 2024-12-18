package fun

import (
	"surf/ast"
	"surf/code"
	args2 "surf/core/engine/args"
	_type "surf/core/engine/type"
	"surf/core/stack"
	"surf/object"
	"surf/tokenUtil"
)

// CallStatement interprets a statement and executes it
func CallStatement(
	statement []code.Token,
	runtime map[string]func(...object.SurfObject),
	isStd bool,
	functions *map[string]map[string]ast.Function,
	variables *stack.Stack,
) {
	call := tokenUtil.ExtractTokensBefore(
		statement,
		code.CloseParen,
		true,
		code.OpenParen,
		code.CloseParen,
	)

	// At this point, only function invocations are allowed
	firstToken := statement[0]

	// Split by commas
	parameters, _ := args2.SplitArgs(call)
	var args []object.SurfObject

	for _, parameter := range parameters {
		args = append(args, _type.TranslateType(parameter[0], variables))
	}

	if isStd {
		stdFunc, found := runtime[firstToken.GetValue()]

		if found {
			stdFunc(args...)
			return
		}
	}

	function, _, _ := ast.LocateFunction(*functions, firstToken.GetFile(), firstToken.GetValue())

	CallFun(
		&function,
		runtime,
		functions,
		firstToken,
		args...,
	)
}
