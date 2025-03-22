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

package reassignment

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/rule/expression"
	"fluent/analyzer/stack"
	"fluent/ast"
	"fluent/filecode"
	"fluent/filecode/types/wrapper"
)

// AnalyzeReassignment analyzes the reassignment of variables in the given AST.
// It checks the left and right expressions for type compatibility and returns an error if any issues are found.
//
// Parameters:
// - tree: the AST to analyze
// - variables: the stack of scoped variables
// - trace: the file code trace
// - assignments: A map of all the assignments throughout the function
//
// Returns:
// - error3.Error: an error object indicating the result of the analysis
func AnalyzeReassignment(
	tree *ast.AST,
	variables *stack.ScopedStack,
	trace *filecode.FileCode,
	collectedAssignments *[]*ast.AST,
) *error3.Error {
	// Get the tree's children
	children := *tree.Children

	// Get the left expression
	leftExpr := children[0]
	rightExpr := children[1]

	// Check if we are supposed to collect reassignments
	if collectedAssignments != nil {
		leftChildren := *leftExpr.Children
		leftNode := leftChildren[0]

		// Check for property access
		if leftNode.Rule == ast.PropertyAccess {
			leftChildren = *leftNode.Children
			// Get the candidate
			candidateExpr := leftChildren[0]
			candidateChildren := *candidateExpr.Children
			candidate := candidateChildren[0]

			if len(leftChildren) == 2 && candidate.Rule == ast.Identifier && *candidate.Value == "this" {
				*collectedAssignments = append(*collectedAssignments, leftNode)
			}
		}
	}

	// Analyze the property access
	obj, err := expression.AnalyzeExpression(
		leftExpr,
		trace,
		variables,
		false,
		&wrapper.TypeWrapper{
			Children: &[]*wrapper.TypeWrapper{},
		},
		true,
		false,
	)

	// Return the err if needed
	if err != nil {
		return err
	}

	// Define the expected type
	expected := obj.Type

	// Analyze the right expression
	obj, err = expression.AnalyzeExpression(
		rightExpr,
		trace,
		variables,
		false,
		&expected,
		false,
		true,
	)

	// Return the err if needed
	if err != nil {
		return err
	}

	return nil
}
