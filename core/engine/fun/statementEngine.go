package fun

import (
	"surf/ast"
	"surf/code"
	args2 "surf/core/engine/args"
	_type "surf/core/engine/type"
	"surf/core/stack"
	"surf/object"
	"surf/token"
	"surf/tokenUtil"
)

// CallStatement interprets a statement and executes it
func CallStatement(
	statement []token.Token,
	runtime map[string]func(...object.SurfObject),
	isStd bool,
	functions *map[string]map[string]*code.Function,
	variables *stack.Stack,
) {
	call := tokenUtil.ExtractTokensBefore(
		statement,
		token.CloseParen,
		true,
		token.OpenParen,
		token.CloseParen,
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
		function,
		runtime,
		functions,
		firstToken,
		args...,
	)
}
