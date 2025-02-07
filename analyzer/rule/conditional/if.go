/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package conditional

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/queue"
	"fluent/analyzer/rule/expression"
	"fluent/analyzer/stack"
	"fluent/ast"
	"fluent/filecode"
	"fluent/filecode/types"
)

// Define a global boolean type
var globalBool = types.TypeWrapper{
	BaseType:    "bool",
	IsPrimitive: true,
	Children:    &[]*types.TypeWrapper{},
}

// ProcessSingleConditional analyzes a single conditional expression and schedules
// the main block for further analysis if the expression is valid.
//
// Parameters:
// - children: A slice of AST nodes representing the conditional expression and main block.
// - trace: The current file code trace.
// - variables: The scoped stack of variables.
// - blockQueue: A queue to schedule blocks for further analysis.
//
// Returns:
// - An error3.Error indicating the result of the analysis.
func ProcessSingleConditional(
	children []*ast.AST,
	trace *filecode.FileCode,
	variables *stack.ScopedStack,
	blockQueue *[]queue.BlockQueueElement,
) error3.Error {
	expr := children[0]
	mainBlock := children[1]

	_, err := expression.AnalyzeExpression(
		expr,
		trace,
		variables,
		false,
		&globalBool,
		false,
	)

	if err.Code != error3.Nothing {
		return err
	}

	// Create a new scope for this conditional
	newScopeId := variables.NewScope()

	// Schedule the block for analysis
	*blockQueue = append(*blockQueue, queue.BlockQueueElement{
		Block: mainBlock,
		ID:    newScopeId,
	})
	return error3.Error{}
}

// AnalyzeIf analyzes an if-else conditional structure in the AST.
//
// Parameters:
// - tree: The AST node representing the if-else structure.
// - trace: The current file code trace.
// - variables: The scoped stack of variables.
// - blockQueue: A queue to schedule blocks for further analysis.
//
// Returns:
// - An error3.Error indicating the result of the analysis.
func AnalyzeIf(
	tree *ast.AST,
	trace *filecode.FileCode,
	variables *stack.ScopedStack,
	blockQueue *[]queue.BlockQueueElement,
) error3.Error {
	// Get the expression and main block
	children := *tree.Children
	err := ProcessSingleConditional(children, trace, variables, blockQueue)

	if err.Code != error3.Nothing {
		return err
	}

	// Transverse the AST to analyze elseif and else nodes
	for i := 2; i < len(children); i++ {
		child := children[i]
		children := *child.Children

		// Determine what to do based on the child's rule
		switch child.Rule {
		case ast.ElseIf:
			err := ProcessSingleConditional(children, trace, variables, blockQueue)

			if err.Code != error3.Nothing {
				return err
			}
		case ast.Else:
			// Create a new scope for this conditional
			newScopeId := variables.NewScope()

			// Schedule the block for analysis
			*blockQueue = append(*blockQueue, queue.BlockQueueElement{
				Block: children[0],
				ID:    newScopeId,
			})
		default:
		}
	}

	return error3.Error{}
}
