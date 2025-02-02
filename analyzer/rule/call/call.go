/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package call

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/object"
	queue2 "fluent/analyzer/queue"
	"fluent/analyzer/stack"
	"fluent/ast"
	"fluent/filecode"
)

func AnalyzeFunctionCall(
	tree *ast.AST,
	trace *filecode.FileCode,
	variables stack.ScopedStack,
	result *object.Object,
	exprQueue *[]queue2.ExpectedPair,
) error3.Error {
	// Get the function's name
	functionName := (*tree.Children)[0].Value
	function, found := trace.Functions[*functionName]

	// Check if the function was found (and whether the current function has permission to call it)
	if !found || (!function.Public && trace.Path != function.Path) {
		return error3.Error{
			Line:       tree.Line,
			Column:     tree.Column,
			Code:       error3.UndefinedReference,
			Additional: *functionName,
		}
	}

	// See if the function call has any parameters
	if len(*tree.Children) < 2 {
		return error3.Error{}
	}

	return error3.Error{}
}
