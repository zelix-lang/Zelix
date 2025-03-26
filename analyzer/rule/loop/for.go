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

package loop

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/object"
	"fluent/analyzer/queue"
	"fluent/analyzer/rule/expression"
	"fluent/analyzer/stack"
	"fluent/analyzer/variable"
	"fluent/ast"
	"fluent/filecode"
	"fluent/filecode/types/wrapper"
)

// AnalyzeFor analyzes a for-loop in the AST.
// It checks the range expressions, ensures the variable name is not redefined,
// and appends the block to the block queue.
//
// Parameters:
// - tree: The AST of the for-loop.
// - trace: The file code trace.
// - variables: The scoped stack of variables.
// - blockQueue: The queue of AST blocks.
// - scopeIds: The scope IDs of the parent block.
//
// Returns:
// - An error3.Error indicating success or failure of the analysis.
func AnalyzeFor(
	tree *ast.AST,
	trace *filecode.FileCode,
	variables *stack.ScopedStack,
	blockQueue *[]queue.BlockQueueElement,
	scopeIds []int,
) (*error3.Error, *string, *variable.Variable) {
	children := *tree.Children

	// Get the range expression
	fromExpr := children[0]
	toExpr := children[1]
	varName := children[2]
	block := children[3]

	// Check if the var name is already defined
	if variables.Load(varName.Value) != nil {
		return &error3.Error{
			Code:       error3.Redefinition,
			Line:       varName.Line,
			Column:     varName.Column,
			Additional: []string{*varName.Value},
		}, nil, nil
	}

	// Analyze the left expression
	leftObj, err := expression.AnalyzeExpression(
		fromExpr,
		trace,
		variables,
		false,
		&wrapper.TypeWrapper{
			Children: &[]*wrapper.TypeWrapper{},
		},
		false,
		false,
	)

	if err != nil {
		return err, nil, nil
	}

	if leftObj.Type.BaseType != "num" && leftObj.Type.BaseType != "dec" {
		return &error3.Error{
			Code:       error3.TypeMismatch,
			Line:       varName.Line,
			Column:     varName.Column,
			Additional: []string{"num or dec", leftObj.Type.Marshal()},
		}, nil, nil
	}

	// Also analyze the right expression
	_, err = expression.AnalyzeExpression(
		toExpr,
		trace,
		variables,
		false,
		&leftObj.Type,
		false,
		false,
	)

	if err != nil {
		return err, nil, nil
	}

	// Append the block to the block queue
	scopeIds = append(scopeIds, variables.Count)
	*blockQueue = append(*blockQueue, queue.BlockQueueElement{
		Block: block,
		// Predict the next ID
		ID:     scopeIds,
		InLoop: true,
	})

	return nil, varName.Value, &variable.Variable{
		Constant: true,
		Value: object.Object{
			Type:   leftObj.Type,
			Value:  nil,
			IsHeap: leftObj.IsHeap,
		},
	}
}
