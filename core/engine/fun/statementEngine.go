package fun

import (
	"zyro/ast"
	"zyro/code"
	args2 "zyro/core/engine/args"
	_type "zyro/core/engine/type"
	"zyro/core/stack"
	"zyro/object"
	"zyro/token"
	"zyro/tokenUtil"
)

// CallStatement interprets a statement and executes it
func CallStatement(
	statement []token.Token,
	runtime map[string]func(...object.ZyroObject),
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
		true,
	)

	// At this point, only function invocations are allowed
	firstToken := statement[0]

	// Split by commas
	parameters, _ := args2.SplitArgs(call)
	var args []object.ZyroObject

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
		args...,
	)
}
