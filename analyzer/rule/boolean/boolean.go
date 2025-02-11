/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package boolean

import (
	"fluent/analyzer/object"
	"fluent/analyzer/queue"
	"fluent/ast"
	"fluent/filecode/types"
)

// Create a global "bool" type that can be reused
var boolType = types.TypeWrapper{
	Children:    &[]*types.TypeWrapper{},
	IsPrimitive: true,
	BaseType:    "bool",
}

// AnalyzeBoolean analyzes a boolean expression in the AST.
// It sets the current element's type to bool and schedules the candidate for evaluation.
// It also schedules expressions for further analysis based on the logical operators.
//
// Parameters:
// - input: The AST node representing the boolean expression.
// - currentElement: The current element being analyzed.
// - exprQueue: The queue of expressions to be analyzed.
func AnalyzeBoolean(
	input *ast.AST,
	currentElement *queue.ExpectedPair,
	exprQueue *[]queue.ExpectedPair,
) {
	// Set the current element's gotten type to bool
	currentElement.Got.Type.BaseType = "bool"

	// Get the input's children
	children := *input.Children
	childrenLen := len(children)

	// Schedule the candidate for evaluation
	candidateElement := queue.ExpectedPair{
		Expected: &types.TypeWrapper{
			Children: &[]*types.TypeWrapper{},
		},
		Got: &object.Object{
			Type: types.TypeWrapper{
				Children: &[]*types.TypeWrapper{},
			},
		},
		Tree: children[0],
	}

	*exprQueue = append(*exprQueue, candidateElement)

	for i := 1; i < childrenLen; i++ {
		el := children[i]

		if el.Rule == ast.BooleanOperator {
			continue
		}

		// Determine the expected type
		var expected *types.TypeWrapper

		prev := children[i-1]

		// Handle logical OR
		if prev.Value != nil && *prev.Value == "||" {
			expected = &boolType
		} else {
			expected = &candidateElement.Got.Type
		}

		// Schedule the expression for analyzing
		*exprQueue = append(*exprQueue, queue.ExpectedPair{
			Expected: expected,
			Got: &object.Object{
				Type: types.TypeWrapper{
					Children: &[]*types.TypeWrapper{},
				},
			},
			Tree: el,
		})
	}
}
