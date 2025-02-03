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
	"fluent/ast"
	"fluent/filecode"
	"fluent/filecode/types"
	"strconv"
)

// AnalyzeFunctionCall analyzes a function call in the AST.
// It checks if the function exists, verifies the number of parameters,
// and updates the result object with the function's return type.
//
// Parameters:
// - tree: The AST of the function call.
// - trace: The file code trace containing function and module definitions.
// - result: The object to store the result of the analysis.
// - exprQueue: The queue to schedule parameter analysis.
//
// Returns:
// - An error3.Error if there is an issue with the function call, otherwise an empty error3.Error.
func AnalyzeFunctionCall(
	tree *ast.AST,
	trace *filecode.FileCode,
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

	// Check that the function call has the correct number of parameters
	if len(*tree.Children) != len(function.Params)+1 {
		return error3.Error{
			Line:       tree.Line,
			Column:     tree.Column,
			Code:       error3.ParameterCountMismatch,
			Additional: strconv.Itoa(len(function.Params)),
		}
	}

	// Update the result accordingly
	returnType := function.ReturnType
	result.IsHeap = returnType.PointerCount > 1
	result.Type = returnType

	// Check for modules
	if !returnType.IsPrimitive {
		// Get the module
		module, found := trace.Modules[returnType.BaseType]

		// Check if the module was found
		if !found {
			return error3.Error{
				Line:       tree.Line,
				Column:     tree.Column,
				Code:       error3.UndefinedReference,
				Additional: returnType.BaseType,
			}
		}

		// Update the result
		result.Value = module
	}

	paramsNode := (*tree.Children)[1]

	// Schedule all the parameters for analysis
	i := 0
	for _, param := range function.Params {
		// Get the parameter's value
		value := (*paramsNode.Children)[i]

		*exprQueue = append(*exprQueue, queue2.ExpectedPair{
			Tree: (*value.Children)[0],
			Got: &object.Object{
				Type: types.TypeWrapper{
					Children: &[]*types.TypeWrapper{},
				},
			},
			Expected: &param.Type,
		})

		i++
	}

	return error3.Error{}
}
