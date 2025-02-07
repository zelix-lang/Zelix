/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package property

import (
	"fluent/analyzer/object"
	"fluent/analyzer/queue"
	"fluent/ast"
	"fluent/filecode/types"
)

// AnalyzePropertyAccess analyzes the property access in the given AST.
// It schedules the candidate and other children for evaluation.
//
// Parameters:
// - input: A pointer to the AST to be analyzed.
// - exprQueue: A pointer to a slice of ExpectedPair to which the analysis results will be appended.
func AnalyzePropertyAccess(
	input *ast.AST,
	exprQueue *[]queue.ExpectedPair,
) {
	children := *input.Children

	// Scheduling the candidate for evaluation
	candidateResult := object.Object{
		Type: types.TypeWrapper{
			Children: &[]*types.TypeWrapper{},
		},
	}

	*exprQueue = append(*exprQueue, queue.ExpectedPair{
		Expected: &types.TypeWrapper{
			Children: &[]*types.TypeWrapper{},
		},
		Got:         &candidateResult,
		Tree:        children[0],
		ModRequired: true,
	})

	// Schedule the other children for evaluation
	for i := 1; i < len(children); i++ {
		*exprQueue = append(*exprQueue, queue.ExpectedPair{
			Expected: &types.TypeWrapper{
				Children: &[]*types.TypeWrapper{},
			},
			Got:          &candidateResult,
			Tree:         children[i],
			IsPropAccess: true,
		})
	}
}
