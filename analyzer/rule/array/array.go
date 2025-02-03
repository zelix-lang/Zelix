/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package array

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/object"
	queue2 "fluent/analyzer/queue"
	"fluent/ast"
	"fluent/filecode/types"
)

func AnalyzeArray(
	tree *ast.AST,
	expected *types.TypeWrapper,
	exprQueue *[]queue2.ExpectedPair,
) error3.Error {
	// Arrays that appear directly as expressions cannot
	// have their type inferred
	if expected.ArrayCount < 1 {
		return error3.Error{
			Code:   error3.CannotInferType,
			Line:   tree.Line,
			Column: tree.Column,
		}
	}

	// Check if the tree has any children
	children := *tree.Children

	if len(children) < 1 {
		// No children, return (Infer the type)
		return error3.Error{}
	}

	// Clone the expected type to determine the expected type
	// individually for each element in the array
	clone := types.TypeWrapper{
		PointerCount: expected.PointerCount,
		ArrayCount:   expected.ArrayCount - 1, // Remove one array count
		Children:     expected.Children,
		BaseType:     expected.BaseType,
		Trace:        expected.Trace,
		IsPrimitive:  expected.IsPrimitive,
	}

	// Schedule all the children to be analyzed
	for _, child := range children {
		*exprQueue = append(*exprQueue, queue2.ExpectedPair{
			Expected: &clone,
			Got: &object.Object{
				Type: types.TypeWrapper{
					Children: &[]*types.TypeWrapper{},
				},
			},
			Tree:              child,
			HasMetDereference: false,
			ActualPointers:    0,
		})
	}

	return error3.Error{}
}
