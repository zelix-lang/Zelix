package fun

import (
	"surf/ast"
	"surf/code"
	args2 "surf/core/engine/args"
	_type "surf/core/engine/type"
	"surf/core/stack"
	"surf/logger"
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

	function, found, sameFile := ast.LocateFunction(*functions, firstToken.GetFile(), firstToken.GetValue())

	if !found {
		logger.TokenError(
			firstToken,
			"Undefined reference to function "+firstToken.GetValue(),
			"The function "+firstToken.GetValue()+" was not found in the file",
			"Add the function to the file",
		)
	}

	if !sameFile && !function.IsPublic() {
		logger.TokenError(
			firstToken,
			"Function "+firstToken.GetValue()+" is not public",
			"Move the function to the current file or make it public",
		)
	}

	CallFun(
		&function,
		runtime,
		functions,
		firstToken,
		args...,
	)
}
