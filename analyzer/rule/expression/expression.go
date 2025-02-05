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
	"fluent/analyzer/rule/array"
	"fluent/analyzer/rule/call"
	"fluent/analyzer/rule/property"
	"fluent/analyzer/stack"
	"fluent/ast"
	"fluent/filecode"
	"fluent/filecode/module"
	"fluent/filecode/types"
)

// literalRules is a map of rules that represent literals
var literalRules = map[ast.Rule]bool{
	ast.StringLiteral:  true,
	ast.NumberLiteral:  true,
	ast.BooleanLiteral: true,
	ast.DecimalLiteral: true,
	ast.Array:          true,
}

func AnalyzeExpression(
	tree *ast.AST,
	trace *filecode.FileCode,
	variables *stack.ScopedStack,
	enforceHeapRequirement bool,
) (object.Object, error3.Error) {
	result := object.Object{
		Type: types.TypeWrapper{
			Children: &[]*types.TypeWrapper{},
		},
	}

	// Use a queue to analyze the expression
	queue := []queue2.ExpectedPair{
		{
			Expected: &types.TypeWrapper{
				Children: &[]*types.TypeWrapper{},
			},
			Got:               &result,
			Tree:              tree,
			HasMetDereference: false,
			ActualPointers:    0,
		},
	}

	// Use the queue
	for len(queue) > 0 {
		// Pop the first element
		element := queue[0]
		queue = queue[1:]

		// Used to skip nodes
		startAt := 0

		// Used to keep track of whether the current value
		// has nested expressions
		hasNested := false

		for _, node := range *element.Tree.Children {
			hasToBreak := false

			switch node.Rule {
			case ast.Pointer:
				startAt++
				// Increment the pointer count
				element.Got.Type.PointerCount++

				if element.HasMetDereference {
					element.ActualPointers++
				}
			case ast.Dereference:
				startAt++
				element.HasMetDereference = true

				// Decrement the pointer count
				element.ActualPointers--
				element.Got.Type.PointerCount--
			default:
				hasToBreak = true
			}

			if hasToBreak {
				break
			}
		}

		// Check for property access
		var lastPropValue module.Module
		if element.IsPropAccess {
			// Check if the lastPropValue is nil
			mod, castOk := element.Got.Value.(module.Module)
			if !castOk {
				return object.Object{}, error3.Error{
					Code:   error3.InvalidPropAccess,
					Line:   element.Tree.Line,
					Column: element.Tree.Column,
				}
			}

			// Clone the got object to keep track of the last property value
			lastPropValue = mod
		}

		// Get the child
		child := (*element.Tree.Children)[startAt]

		// See if the address of this value can be taken
		if literalRules[child.Rule] && element.Got.Type.PointerCount > 0 {
			return object.Object{}, error3.Error{
				Code:   error3.CannotTakeAddress,
				Line:   child.Line,
				Column: child.Column,
			}
		}

		switch child.Rule {
		case ast.StringLiteral:
			element.Got.Type.BaseType = "str"
			element.Got.Type.IsPrimitive = true
		case ast.NumberLiteral:
			element.Got.Type.BaseType = "num"
			element.Got.Type.IsPrimitive = true
		case ast.BooleanLiteral:
			element.Got.Type.BaseType = "bool"
			element.Got.Type.IsPrimitive = true
		case ast.DecimalLiteral:
			element.Got.Type.BaseType = "dec"
			element.Got.Type.IsPrimitive = true
		case ast.Identifier:
			// Check for property access
			if element.IsPropAccess {
				err := property.ProcessPropIdentifier(
					lastPropValue,
					&element,
					trace,
					child,
				)

				// Return the error if it is not nothing
				if err.Code != error3.Nothing {
					return object.Object{}, err
				}
				continue
			}

			// Check if the variable exists
			value := variables.Load(child.Value)

			if value == nil {
				return object.Object{}, error3.Error{
					Code:       error3.UndefinedReference,
					Additional: []string{*child.Value},
					Line:       element.Tree.Line,
					Column:     element.Tree.Column,
				}
			}

			oldPointerCount := element.Got.Type.PointerCount
			element.Got.Type = value.Value.Type
			element.Got.Type.PointerCount += oldPointerCount
			element.ActualPointers += value.Value.Type.PointerCount
			element.Got.Value = value.Value.Value
			element.Got.IsHeap = value.Value.IsHeap
		case ast.Array:
			err := array.AnalyzeArray(child, element.Expected, &queue)

			// Return the error if it is not nothing
			if err.Code != error3.Nothing {
				return object.Object{}, err
			}

			element.Got.Type = *element.Expected
		case ast.FunctionCall:
			// This will later be fully determined by the call analyzer
			element.Got.IsHeap = false

			// Pass the input to the function call analyzer
			err := call.AnalyzeFunctionCall(
				child,
				trace,
				&element,
				&queue,
				lastPropValue,
			)

			// Return the error if it is not nothing
			if err.Code != error3.Nothing {
				return object.Object{}, err
			}
		case ast.Expression:
			hasNested = true
			// Add the expression to the queue
			queue = append(queue, queue2.ExpectedPair{
				Expected:          element.Expected,
				Got:               element.Got,
				Tree:              child,
				HasMetDereference: element.HasMetDereference,
				ActualPointers:    element.ActualPointers,
			})

			element.Got.Type = *element.Expected
		case ast.PropertyAccess:
			// This will later be fully determined by the property access analyzer
			element.Got.IsHeap = false

			// Pass the input to the property access analyzer
			err := property.AnalyzePropertyAccess(
				child,
				&queue,
			)

			// Return the error if it is not nothing
			if err.Code != error3.Nothing {
				return object.Object{}, err
			}
		default:
		}

		// isInferred does not work here because it was defined
		// before the switch statement
		if element.Expected.BaseType == "(Infer)" {
			oldPointerCount := element.Expected.PointerCount
			oldArrayCount := element.Expected.ArrayCount

			*element.Expected = element.Got.Type
			element.Expected.PointerCount = oldPointerCount
			element.Expected.ArrayCount = oldArrayCount
		}

		// Check if the pointer count is negative
		if !hasNested && element.ActualPointers < 0 {
			return object.Object{}, error3.Error{
				Code:   error3.InvalidDereference,
				Line:   element.Tree.Line,
				Column: element.Tree.Column,
			}
		}

		// Check for type mismatch
		if element.Expected.BaseType != "" && !element.Expected.Compare(element.Got.Type) {
			return object.Object{}, error3.Error{
				Code:       error3.TypeMismatch,
				Line:       element.Tree.Line,
				Column:     element.Tree.Column,
				Additional: []string{element.Expected.Marshal(), element.Got.Type.Marshal()},
			}
		}

		// Check if the data escapes the function
		if enforceHeapRequirement && element.HeapRequired && !element.Got.IsHeap {
			return object.Object{}, error3.Error{
				Code:   error3.DataOutlivesStack,
				Line:   element.Tree.Line,
				Column: element.Tree.Column,
			}
		}

		if element.ModRequired && element.Got.Value == nil {
			return object.Object{}, error3.Error{
				Code:       error3.TypeMismatch,
				Line:       element.Tree.Line,
				Column:     element.Tree.Column,
				Additional: []string{"Module", element.Got.Type.Marshal()},
			}
		}
	}

	return result, error3.Error{}
}
