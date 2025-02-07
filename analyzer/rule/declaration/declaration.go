/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package declaration

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/format"
	"fluent/analyzer/rule/expression"
	"fluent/analyzer/rule/value"
	"fluent/analyzer/stack"
	"fluent/analyzer/variable"
	"fluent/ast"
	"fluent/filecode"
	"fluent/filecode/types"
)

// AnalyzeDeclaration analyzes a declaration statement in the AST.
// It checks for undefined references, ensures the name is in snake_case,
// and analyzes the expression associated with the declaration.
//
// Parameters:
// - statement: The AST node representing the declaration statement.
// - scope: The scoped stack for variable management.
// - trace: The file code trace for error reporting.
// - generics: A map of generics used in the declaration.
//
// Returns:
// - An error if there is an issue with the declaration.
// - A warning if the name is not in snake_case.
func AnalyzeDeclaration(
	statement *ast.AST,
	scope *stack.ScopedStack,
	trace *filecode.FileCode,
	generics *map[string]bool,
) (error3.Error, error3.Error) {
	// Get the children of the statement
	children := *statement.Children
	isConst := *children[0].Value == "const"
	nameNode := *children[1]
	typeNode := *children[2]
	expr := *children[3]

	// Convert the type node to a TypeWrapper
	typeWrapper := types.ConvertToTypeWrapper(typeNode)

	// Check for undefined references in the type
	err := value.AnalyzeUndefinedReference(trace, typeWrapper, generics)

	if err.Code != error3.Nothing {
		return err, error3.Error{}
	}

	var warning error3.Error
	// Check that the name is in snake_case
	if !format.CheckCase(nameNode.Value, format.SnakeCase) {
		warning = error3.Error{
			Code:   error3.NameShouldBeSnakeCase,
			Line:   nameNode.Line,
			Column: nameNode.Column,
		}
	}

	// Pass the expression to the expression analyzer
	obj, err := expression.AnalyzeExpression(&expr, trace, scope, false, &typeWrapper, false)

	if err.Code != error3.Nothing {
		return err, error3.Error{}
	}

	// Save the object in the variables
	scope.Append(*nameNode.Value, variable.Variable{
		Constant: isConst,
		Value:    obj,
	})

	return error3.Error{}, warning
}
