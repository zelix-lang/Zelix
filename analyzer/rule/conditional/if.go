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

// processSingleConditional analyzes a single conditional expression and schedules
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
func processSingleConditional(
	children []*ast.AST,
	trace *filecode.FileCode,
	variables *stack.ScopedStack,
	blockQueue *[]ast.AST,
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

	// Schedule the block for analysis
	*blockQueue = append(*blockQueue, *mainBlock)
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
	blockQueue *[]ast.AST,
) error3.Error {
	// Get the expression and main block
	children := *tree.Children
	err := processSingleConditional(children, trace, variables, blockQueue)

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
			err := processSingleConditional(children, trace, variables, blockQueue)

			if err.Code != error3.Nothing {
				return err
			}
		case ast.Else:
			// Schedule the block for analysis
			*blockQueue = append(*blockQueue, *children[0])
		default:
		}
	}

	return error3.Error{}
}
