/*
   The Fluent Programming Language
   -----------------------------------------------------
   This code is released under the GNU GPL v3 license.
   For more information, please visit:
   https://www.gnu.org/licenses/gpl-3.0.html
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package call

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/object"
	"fluent/analyzer/property"
	queue2 "fluent/analyzer/queue"
	"fluent/ast"
	"fluent/filecode"
	function2 "fluent/filecode/function"
	"fluent/filecode/types/wrapper"
	"strconv"
)

// AnalyzeFunctionCall analyzes a function call or object creation in the AST.
// It checks if the function exists, verifies the number of parameters,
// and updates the result object with the function's return type.
//
// Parameters:
// - tree: The AST of the function call.
// - trace: The file code trace containing function and module definitions.
// - queueElement: The current element in the queue.
// - exprQueue: The queue to schedule parameter analysis.
// - isObjectCreation: A boolean indicating whether the tree represents an object creation.
//
// Returns:
// - An error3.Error if there is an issue with the function call, otherwise an empty error3.Error.
func AnalyzeFunctionCall(
	tree *ast.AST,
	trace *filecode.FileCode,
	queueElement *queue2.ExpectedPair,
	exprQueue *[]queue2.ExpectedPair,
	isObjectCreation bool,
) *error3.Error {
	// Get the function's name
	functionName := (*tree.Children)[0].Value
	// Determine if the function has parameters
	hasParams := len(*tree.Children) > 1

	var generics map[string]bool
	var function *function2.Function
	var found bool
	var returnType wrapper.TypeWrapper

	if isObjectCreation {
		// Find the module inside the trace's module
		mod, ok := trace.Modules[*functionName]

		if !ok || (!mod.Public && trace.Path != mod.Path) {
			return &error3.Error{
				Line:       tree.Line,
				Column:     tree.Column,
				Code:       error3.UndefinedReference,
				Additional: []string{*functionName},
			}
		}

		// See if the module has a constructor
		constructor, ok := mod.Functions[mod.Name]

		if !ok {
			if hasParams {
				return &error3.Error{
					Line:   tree.Line,
					Column: tree.Column,
					Code:   error3.DoesNotHaveConstructor,
				}
			}

			queueElement.Got.Type.BaseType = mod.Name
			queueElement.Got.Value = mod
			return nil
		}

		function, found = constructor, true
		generics = mod.Templates
		returnType = wrapper.TypeWrapper{
			BaseType: mod.Name,
			Children: &[]*wrapper.TypeWrapper{},
		}
	} else if queueElement.IsPropAccess {
		lastPropValue := property.EvaluateLastPropValue(queueElement)
		if lastPropValue == nil {
			return &error3.Error{
				Code:   error3.InvalidPropAccess,
				Line:   queueElement.Tree.Line,
				Column: queueElement.Tree.Column,
			}
		}

		function, found = lastPropValue.Functions[*functionName]
		generics = function.Templates
		returnType = function.ReturnType

		if !returnType.IsPrimitive {
			*queueElement.LastPropValue = trace.Modules[returnType.BaseType]
		}
	} else {
		function, found = trace.Functions[*functionName]
		generics = function.Templates
		returnType = function.ReturnType
	}

	// Check if the function was found (and whether the current function has permission to call it)
	if !found || (!function.Public && trace.Path != function.Path) {
		return &error3.Error{
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
		if _, found := generics[returnType.BaseType]; found {
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
				return &error3.Error{
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
	preventParamAnalysis := false
	if hasParams {
		paramsNode := (*tree.Children)[1]

		// Check that the function call has the correct number of parameters
		if function.IsStd && function.Name == "panic" {
			if len(*paramsNode.Children) == 1 {
				preventParamAnalysis = true
			} else {
				return &error3.Error{
					Line:       tree.Line,
					Column:     tree.Column,
					Code:       error3.ParameterCountMismatch,
					Additional: []string{"1"},
				}
			}
		}

		if len(*paramsNode.Children) != len(function.Params) {
			if !preventParamAnalysis {
				return &error3.Error{
					Line:       tree.Line,
					Column:     tree.Column,
					Code:       error3.ParameterCountMismatch,
					Additional: []string{strconv.Itoa(len(function.Params))},
				}
			}
		}

		// Schedule all the parameters for analysis
		if !preventParamAnalysis {
			i := 0

			// Prevent inferring twice the same generics
			inferredGenerics := make(map[string]*object.Object)

			for _, param := range function.Params {
				// Get the parameter's value
				value := (*paramsNode.Children)[i]
				paramType := param.Type
				paramNodes := (*value.Children)[0]
				isParamHeap := paramType.PointerCount > 0
				enforceHeapInParam := false
				passedResult := object.Object{
					Type: wrapper.TypeWrapper{
						Children: &[]*wrapper.TypeWrapper{},
					},
					IsHeap: isParamHeap,
				}
				passedType := &paramType

				if !param.Type.IsPrimitive {
					// Check for generics
					if _, found := generics[param.Type.BaseType]; found {
						// See if we have seen this generic before
						if inferredGenerics[param.Type.BaseType] != nil {
							passedType = &inferredGenerics[param.Type.BaseType].Type
						} else {
							inferredGenerics[param.Type.BaseType] = &passedResult
							// Check if this param has the return type's generic
							if returnType.Compare(paramType) {
								enforceHeapInParam = result.IsHeap
								if expected.BaseType == "" {
									passedType = &wrapper.TypeWrapper{
										BaseType:     "(Infer)",
										PointerCount: param.Type.PointerCount,
										ArrayCount:   param.Type.ArrayCount,
									}
								} else {
									passedType = expected
								}
							} else {
								passedType.BaseType = "(Infer)"
							}
						}
					} else {
						passedType = &paramType
					}
				} else {
					passedType = &paramType
				}

				*exprQueue = append(*exprQueue, queue2.ExpectedPair{
					Tree:         paramNodes,
					Got:          &passedResult,
					Expected:     passedType,
					HeapRequired: enforceHeapInParam && enforceHeap,
					IsParam:      true,
				})

				i++
			}
		}
	} else {
		// Check for parameter count mismatch
		if len(function.Params) > 0 {
			return &error3.Error{
				Line:       tree.Line,
				Column:     tree.Column,
				Code:       error3.ParameterCountMismatch,
				Additional: []string{strconv.Itoa(len(function.Params))},
			}
		}
	}

	return nil
}
