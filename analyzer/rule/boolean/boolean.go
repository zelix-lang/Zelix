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

package boolean

import (
	"fluent/analyzer/object"
	"fluent/analyzer/queue"
	"fluent/ast"
	"fluent/filecode/types/wrapper"
)

// Create a global "bool" type that can be reused
var boolType = wrapper.TypeWrapper{
	Children:    &[]*wrapper.TypeWrapper{},
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
	startAt := 0

	// Handle negations
	for children[startAt].Rule == ast.BooleanOperator {
		startAt++
	}

	// Schedule the candidate for evaluation
	candidateElement := queue.ExpectedPair{
		Expected: &wrapper.TypeWrapper{
			Children: &[]*wrapper.TypeWrapper{},
		},
		Got: &object.Object{
			Type: wrapper.TypeWrapper{
				Children: &[]*wrapper.TypeWrapper{},
			},
		},
		Tree:    children[startAt],
		IsParam: true,
	}

	*exprQueue = append(*exprQueue, candidateElement)

	for i := startAt + 1; i < childrenLen; i++ {
		el := children[i]

		if el.Rule == ast.BooleanOperator {
			continue
		}

		// Determine the expected type
		var expected *wrapper.TypeWrapper

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
				Type: wrapper.TypeWrapper{
					Children: &[]*wrapper.TypeWrapper{},
				},
			},
			Tree:    el,
			IsParam: true,
		})
	}
}
