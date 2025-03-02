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

package conditional

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/queue"
	"fluent/analyzer/rule/expression"
	"fluent/analyzer/stack"
	"fluent/ast"
	"fluent/filecode"
	"fluent/filecode/types/wrapper"
)

// Define a global boolean type
var globalBool = wrapper.TypeWrapper{
	BaseType:    "bool",
	IsPrimitive: true,
	Children:    &[]*wrapper.TypeWrapper{},
}

// ProcessSingleConditional analyzes a single conditional expression and schedules
// the main block for further analysis if the expression is valid.
//
// Parameters:
// - children: A slice of AST nodes representing the conditional expression and main block.
// - trace: The current file code trace.
// - variables: The scoped stack of variables.
// - blockQueue: A queue to schedule blocks for further analysis.
// - scopeIds: The scope IDs of the parent block.
// - inLoop: A flag indicating whether the conditional is inside a loop.
//
// Returns:
// - An error3.Error indicating the result of the analysis.
func ProcessSingleConditional(
	children []*ast.AST,
	trace *filecode.FileCode,
	variables *stack.ScopedStack,
	blockQueue *[]queue.BlockQueueElement,
	scopeIds []int,
	inLoop bool,
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
		false,
	)

	if err.Code != error3.Nothing {
		return err
	}

	// Create a new scope for this conditional
	newScopeId := variables.NewScope()

	scopeIds = append(scopeIds, newScopeId)
	// Schedule the block for analysis
	*blockQueue = append(*blockQueue, queue.BlockQueueElement{
		Block:  mainBlock,
		ID:     scopeIds,
		InLoop: inLoop,
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
// - scopeIds: The scope IDs of the parent block.
// - inLoop: A flag indicating whether the if-else structure is inside a loop.
//
// Returns:
// - An error3.Error indicating the result of the analysis.
func AnalyzeIf(
	tree *ast.AST,
	trace *filecode.FileCode,
	variables *stack.ScopedStack,
	blockQueue *[]queue.BlockQueueElement,
	scopeIds []int,
	inLoop bool,
) error3.Error {
	// Get the expression and main block
	children := *tree.Children
	err := ProcessSingleConditional(children, trace, variables, blockQueue, scopeIds, inLoop)

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
			err := ProcessSingleConditional(children, trace, variables, blockQueue, scopeIds, inLoop)

			if err.Code != error3.Nothing {
				return err
			}
		case ast.Else:
			// Create a new scope for this conditional
			newScopeId := variables.NewScope()

			// Schedule the block for analysis
			scopeIds = append(scopeIds, newScopeId)
			*blockQueue = append(*blockQueue, queue.BlockQueueElement{
				Block:  children[0],
				ID:     scopeIds,
				InLoop: inLoop,
			})
		default:
		}
	}

	return error3.Error{}
}
