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
	function2 "fluent/filecode/function"
	"fluent/filecode/module"
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
// - queueElement: The current element in the queue.
// - exprQueue: The queue to schedule parameter analysis.
//
// Returns:
// - An error3.Error if there is an issue with the function call, otherwise an empty error3.Error.
func AnalyzeFunctionCall(
	tree *ast.AST,
	trace *filecode.FileCode,
	queueElement *queue2.ExpectedPair,
	exprQueue *[]queue2.ExpectedPair,
	lastPropValue *module.Module,
) error3.Error {
	// Get the function's name
	functionName := (*tree.Children)[0].Value

	var function function2.Function
	var found bool

	if !queueElement.IsPropAccess {
		function, found = trace.Functions[*functionName]
	} else {
		function, found = lastPropValue.Functions[*functionName]
	}

	// Check if the function was found (and whether the current function has permission to call it)
	if !found || (!function.Public && trace.Path != function.Path) {
		return error3.Error{
			Line:       tree.Line,
			Column:     tree.Column,
			Code:       error3.UndefinedReference,
			Additional: []string{*functionName},
		}
	}

	result := queueElement.Got
	expected := queueElement.Expected

	// Update the result accordingly
	oldPointerCount := result.Type.PointerCount
	oldArrayCount := result.Type.ArrayCount
	returnType := function.ReturnType

	// Heap and return type
	result.IsHeap = returnType.PointerCount > 0
	result.Type = returnType

	// Do not enforce heap tracing if the current function is heap_alloc
	enforceHeap := function.Name != "heap_alloc" && !function.IsStd

	// Honor pointers and arrays
	result.Type.PointerCount += oldPointerCount
	result.Type.ArrayCount += oldArrayCount
	queueElement.ActualPointers += returnType.PointerCount

	if !returnType.IsPrimitive {
		// Check for generics
		if _, found := function.Templates[returnType.BaseType]; found {
			// Make sure there is an expected type
			if expected.BaseType == "" {
				result.Type.BaseType = returnType.BaseType
				result.Type.Children = returnType.Children
			} else {
				oldPointerCount := result.Type.PointerCount
				oldArrayCount := result.Type.ArrayCount

				// Make the expression analyzer infer the type
				result.Type = *expected

				if result.Type.BaseType != "(Infer)" {
					result.Type.PointerCount = oldPointerCount
					result.Type.ArrayCount = oldArrayCount
				}
			}
		} else {
			// Get the module
			mod, found := trace.Modules[returnType.BaseType]

			// Check if the module was found
			if !found {
				return error3.Error{
					Line:       tree.Line,
					Column:     tree.Column,
					Code:       error3.UndefinedReference,
					Additional: []string{returnType.BaseType},
				}
			}

			// Update the result
			result.Value = mod
		}
	}

	// See if the function call has any parameters
	if len(*tree.Children) > 1 {
		paramsNode := (*tree.Children)[1]

		// Check that the function call has the correct number of parameters
		if len(*paramsNode.Children) != len(function.Params) {
			return error3.Error{
				Line:       tree.Line,
				Column:     tree.Column,
				Code:       error3.ParameterCountMismatch,
				Additional: []string{strconv.Itoa(len(function.Params))},
			}
		}

		// Schedule all the parameters for analysis
		i := 0
		for _, param := range function.Params {
			// Get the parameter's value
			value := (*paramsNode.Children)[i]
			paramType := param.Type
			paramNodes := (*value.Children)[0]
			isParamHeap := paramType.PointerCount > 0
			enforceHeapInParam := false

			if !param.Type.IsPrimitive {
				// Check for generics
				if _, found := function.Templates[param.Type.BaseType]; found {
					// Check if this param has the return type's generic
					if returnType.Compare(paramType) {
						enforceHeapInParam = result.IsHeap
						if expected.BaseType == "" {
							paramType = types.TypeWrapper{
								BaseType:     "(Infer)",
								PointerCount: param.Type.PointerCount,
								ArrayCount:   param.Type.ArrayCount,
							}
						} else {
							paramType = *expected
						}
					} else {
						paramType.BaseType = "(Infer)"
					}
				}
			}

			*exprQueue = append(*exprQueue, queue2.ExpectedPair{
				Tree: paramNodes,
				Got: &object.Object{
					Type: types.TypeWrapper{
						Children: &[]*types.TypeWrapper{},
					},
					IsHeap: isParamHeap,
				},
				Expected:     &paramType,
				HeapRequired: enforceHeapInParam && enforceHeap,
			})

			i++
		}
	} else {
		// Check for parameter count mismatch
		if len(function.Params) > 0 {
			return error3.Error{
				Line:       tree.Line,
				Column:     tree.Column,
				Code:       error3.ParameterCountMismatch,
				Additional: []string{strconv.Itoa(len(function.Params))},
			}
		}
	}

	return error3.Error{}
}
