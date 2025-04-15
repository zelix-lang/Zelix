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

package property

import (
	"fluent/analyzer/object"
	"fluent/analyzer/queue"
	"fluent/ast"
	"fluent/filecode/types/wrapper"
)

// AnalyzePropertyAccess analyzes the property access in the given AST.
// It schedules the candidate and other children for evaluation.
//
// Parameters:
// - currentElement: A pointer to the current ExpectedPair in the queue.
// - child: A pointer to the AST node to be analyzed.
// - exprQueue: A pointer to a slice of ExpectedPair to which the analysis results will be appended.
// - isPropReassignment: A boolean indicating whether the current element comes from a property reassignment.
func AnalyzePropertyAccess(
	currentElement *queue.ExpectedPair,
	child *ast.AST,
	exprQueue *[]queue.ExpectedPair,
	isPropReassignment bool,
) {
	children := *child.Children

	// Scheduling the candidate for evaluation
	candidateResult := object.Object{
		Type: wrapper.TypeWrapper{
			Children: &[]*wrapper.TypeWrapper{},
		},
	}

	currentElement.Got.Value = &candidateResult.Value

	// Temporarily save the children to be appended in a new slice
	newChildren := make([]queue.ExpectedPair, len(children))

	newChildren[0] = queue.ExpectedPair{
		Expected: &wrapper.TypeWrapper{
			Children: &[]*wrapper.TypeWrapper{},
		},
		Got:         &candidateResult,
		Tree:        children[0],
		ModRequired: true,
	}

	// Schedule the other children for evaluation
	childrenLen := len(children) - 1
	for i := 1; i <= childrenLen; i++ {
		newChildren[i] = queue.ExpectedPair{
			Expected: &wrapper.TypeWrapper{
				Children: &[]*wrapper.TypeWrapper{},
			},
			Got:                currentElement.Got,
			Tree:               children[i],
			IsPropAccess:       true,
			IsPropReassignment: isPropReassignment && i == childrenLen,
			LastPropValue:      &candidateResult.Value,
		}

		if i == childrenLen {
			newChildren[i].Expected = currentElement.Expected
		}
	}

	*exprQueue = append(newChildren, *exprQueue...)
}
