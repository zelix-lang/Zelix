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
	trace2 "fluent/filecode/trace"
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
// - parentName: The name of the module that holds this function, empty if none.
//
// Returns:
// - An error if there is an issue with the declaration.
// - A warning if the name is not in snake_case.
func AnalyzeDeclaration(
	statement *ast.AST,
	scope *stack.ScopedStack,
	trace *filecode.FileCode,
	generics *map[string]bool,
	parentName string,
) (*error3.Error, *error3.Error) {
	// Get the children of the statement
	children := *statement.Children
	isConst := *children[0].Value == "const"
	nameNode := *children[1]
	typeNode := *children[2]
	expr := *children[3]

	// Check if the var name is already defined
	if scope.Load(nameNode.Value) != nil {
		return &error3.Error{
			Code:       error3.Redefinition,
			Line:       nameNode.Line,
			Column:     nameNode.Column,
			Additional: []string{*nameNode.Value},
		}, nil
	}

	// Convert the type node to a TypeWrapper
	typeWrapper := types.ConvertToTypeWrapper(typeNode)

	// Check for undefined references in the type
	err := value.AnalyzeUndefinedReference(trace, typeWrapper, generics)

	if err != nil {
		return err, nil
	}

	// Check for illegal self-reference
	if typeWrapper.BaseType == parentName {
		exprChildren := *expr.Children
		candidate := exprChildren[0]
		if candidate.Rule != ast.Identifier || (candidate.Rule == ast.Identifier && *candidate.Value != "this") {
			return &error3.Error{
				Code:   error3.SelfReference,
				Line:   expr.Line,
				Column: expr.Column,
			}, nil
		}
	}

	var warning *error3.Error
	// Check that the name is in snake_case
	if !format.CheckCase(nameNode.Value, format.SnakeCase) {
		warning = &error3.Error{
			Code:   error3.NameShouldBeSnakeCase,
			Line:   nameNode.Line,
			Column: nameNode.Column,
		}
	}

	// Pass the expression to the expression analyzer
	obj, err := expression.AnalyzeExpression(&expr, trace, scope, false, &typeWrapper, false, true)

	if err != nil {
		return err, nil
	}

	// Save the object in the variables
	scope.Append(*nameNode.Value, variable.Variable{
		Constant: isConst,
		Value:    *obj,
		Trace: trace2.Trace{
			Line:   nameNode.Line,
			Column: nameNode.Column,
		},
	})

	return nil, warning
}
