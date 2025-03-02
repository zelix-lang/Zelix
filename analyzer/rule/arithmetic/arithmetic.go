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

package arithmetic

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/object"
	"fluent/analyzer/queue"
	"fluent/ast"
	"fluent/filecode/types/wrapper"
)

// AnalyzeArithmetic analyzes an arithmetic expression in the AST.
// It checks if the type of the expression can be inferred and if it matches the expected type.
// If the type cannot be inferred or does not match, it returns an error.
// Otherwise, it processes the expression and updates the expression queue.
//
// Parameters:
// - input: The AST node representing the arithmetic expression.
// - currentElement: The current expected pair being analyzed.
// - exprQueue: The queue of expected pairs to be processed.
//
// Returns:
// - An error3.Error indicating the result of the analysis.
func AnalyzeArithmetic(
	input *ast.AST,
	currentElement *queue.ExpectedPair,
	exprQueue *[]queue.ExpectedPair,
) error3.Error {
	expected := currentElement.Expected

	// Check if we can infer the type of the expression
	if expected.BaseType == "" {
		return error3.Error{
			Code:   error3.CannotInferType,
			Line:   input.Line,
			Column: input.Column,
		}
	}

	// Check if the expected is either a num or a dec
	if expected.BaseType != "num" && expected.BaseType != "dec" && expected.BaseType != "(Infer)" {
		return error3.Error{
			Code:       error3.TypeMismatch,
			Line:       input.Line,
			Column:     input.Column,
			Additional: []string{expected.BaseType, "arithmetic"},
		}
	}

	// Handle inferred types
	var candidateType *wrapper.TypeWrapper
	if expected.BaseType == "(Infer)" {
		candidateType = &wrapper.TypeWrapper{
			Children: &[]*wrapper.TypeWrapper{},
		}
	} else {
		// Clone the expected element to avoid memory issues
		candidateType = &wrapper.TypeWrapper{
			PointerCount: expected.PointerCount,
			ArrayCount:   expected.ArrayCount,
			Children:     expected.Children,
			BaseType:     expected.BaseType,
			Trace:        expected.Trace,
			IsPrimitive:  expected.IsPrimitive,
		}
	}

	startAt := 0
	// Get the expression's children
	children := *input.Children

	// Evaluate the candidate only if the type is inferred
	if expected.BaseType == "(Infer)" {
		// Push the candidate to determine the expression's type
		candidate := children[0]
		candidateElement := object.Object{
			Type: wrapper.TypeWrapper{
				Children: &[]*wrapper.TypeWrapper{},
			},
		}

		// Prevent nesting problems
		if candidate.Rule != ast.Expression {
			startAt = 1
		}

		for candidate.Rule == ast.Expression {
			newCandidate := (*candidate.Children)[0]

			if newCandidate.Rule == ast.Expression {
				candidate = newCandidate
				continue
			}

			if newCandidate.Rule != ast.ArithmeticExpression {
				break
			}

			candidate = newCandidate
		}

		*exprQueue = append(*exprQueue, queue.ExpectedPair{
			Expected:     candidateType,
			Got:          &candidateElement,
			Tree:         candidate,
			IsArithmetic: true,
			IsParam:      true,
		})

		candidateType = &candidateElement.Type
	}

	// Push the rest of the expression
	for i := startAt; i < len(children); i++ {
		element := children[i]

		// Skip sings
		if element.Rule == ast.ArithmeticSign {
			continue
		}

		*exprQueue = append(*exprQueue, queue.ExpectedPair{
			Expected: candidateType,
			Got: &object.Object{
				Type: wrapper.TypeWrapper{
					Children: &[]*wrapper.TypeWrapper{},
				},
			},
			Tree:    element,
			IsParam: true,
		})
	}

	currentElement.Got.Type.BaseType = currentElement.Expected.BaseType
	return error3.Error{}
}
