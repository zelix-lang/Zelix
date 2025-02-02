/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package expression

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/object"
	queue2 "fluent/analyzer/queue"
	"fluent/analyzer/rule/call"
	"fluent/analyzer/stack"
	"fluent/ast"
	"fluent/filecode"
	"fluent/filecode/types"
)

func AnalyzeExpression(tree *ast.AST, trace *filecode.FileCode, variables stack.ScopedStack) (object.Object, error3.Error) {
	result := object.Object{}

	// Use a queue to analyze the expression
	queue := []queue2.ExpectedPair{
		{
			Expected: &types.TypeWrapper{},
			Got:      &object.Object{},
			Tree:     tree,
		},
	}

	// Use the queue
	for len(queue) > 0 {
		result := object.Object{}

		// Pop the first element
		element := queue[0]
		queue = queue[1:]

		// Used to skip tokens
		startAt := 0

		actualValuePointers := 0
		hasMetDereference := false

		for _, node := range *element.Tree.Children {
			hasToBreak := false

			switch node.Rule {
			case ast.Pointer:
				startAt++
				// Increment the pointer count
				result.Type.PointerCount++

				if hasMetDereference {
					actualValuePointers++
				}
			case ast.Dereference:
				startAt++
				hasMetDereference = true

				// Decrement the pointer count
				actualValuePointers--
				result.Type.PointerCount--
			default:
				hasToBreak = true
			}

			if hasToBreak {
				break
			}
		}

		child := (*element.Tree.Children)[startAt]

		switch child.Rule {
		case ast.Identifier:
			// Check if the variable exists
			value := variables.Load(child.Value)

			if value == nil {
				return object.Object{}, error3.Error{
					Code:       error3.UndefinedReference,
					Additional: *child.Value,
					Line:       element.Tree.Line,
					Column:     element.Tree.Column,
				}
			}
		case ast.FunctionCall:
			// Pass the input to the function call analyzer
			err := call.AnalyzeFunctionCall(child, trace, variables, &result, &queue)

			// Return the error if it is not nothing
			if err.Code != error3.Nothing {
				return object.Object{}, err
			}
		case ast.Expression:
			// Add the expression to the queue
			queue = append(queue, queue2.ExpectedPair{
				Expected: element.Expected,
				Got:      element.Got,
				Tree:     child,
			})
		default:
		}

		// Check if the pointer count is negative
		if actualValuePointers < 0 {
			return object.Object{}, error3.Error{
				Code:   error3.InvalidDereference,
				Line:   element.Tree.Line,
				Column: element.Tree.Column,
			}
		}

		// Check for type mismatch
		if element.Expected.BaseType != "" && !element.Expected.Compare(result.Type) {
			return object.Object{}, error3.Error{
				Code:   error3.TypeMismatch,
				Line:   element.Tree.Line,
				Column: element.Tree.Column,
			}
		}

		element.Got = &result
	}

	return result, error3.Error{}
}
