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

package ret

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/rule/expression"
	"fluent/token"
)

// ProcessReturn processes a return statement in the Fluent programming language.
// It takes a slice of tokens representing the statement and returns an AST node
// and an error if the statement is invalid.
//
// Parameters:
//   - statement: A slice of token.Token representing the return statement.
//
// Returns:
//   - ast.AST: The abstract syntax tree node representing the return statement.
//   - error.Error: An error object if the statement is invalid, otherwise an empty error.
func ProcessReturn(statement []token.Token) (*ast.AST, *error.Error) {
	// There is no need to check the input's length
	// "return;" is valid anyway, which means the input is 1 token long
	// hence, the input's length is always valid

	result := ast.AST{
		Rule:     ast.Return,
		Line:     statement[0].Line,
		Column:   statement[0].Column,
		File:     &statement[0].File,
		Children: &[]*ast.AST{},
	}

	// See if the statement has any expression
	if len(statement) > 1 {
		// Parse the expression
		expr, err := expression.ProcessExpression(statement[1:])

		if err != nil {
			return nil, err
		}

		// Add the expression to the result
		*result.Children = append(*result.Children, expr)
	}

	return &result, nil
}
