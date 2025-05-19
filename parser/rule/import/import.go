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

package _import

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/token"
)

// ProcessImport processes a list of tokens representing an import statement.
// It expects exactly two tokens: the 'import' keyword and a string literal.
// If the tokens are not as expected, it returns an error.
// If the tokens are valid, it returns an AST node representing the import statement.
//
// Parameters:
//
//	tokens []token.Token - A slice of tokens representing the import statement.
//
// Returns:
//
//	ast.AST - The AST node representing the import statement.
//	error.Error - An error object if the tokens are invalid.
func ProcessImport(tokens []token.Token) (*ast.AST, *error.Error) {
	// The first token is always an import keyword
	// We are not going to process it

	if len(tokens) != 2 {
		return nil, &error.Error{
			Line:     tokens[0].Line,
			Column:   tokens[0].Column,
			File:     &tokens[0].File,
			Expected: []ast.Rule{ast.StringLiteral},
		}
	}

	// Basic validation for the string literal
	if tokens[1].TokenType != token.StringLiteral {
		return nil, &error.Error{
			Line:     tokens[1].Line,
			Column:   tokens[1].Column,
			File:     &tokens[1].File,
			Expected: []ast.Rule{ast.StringLiteral},
		}
	}

	// Directly parse the import rule
	child := ast.AST{
		Rule: ast.Import,
		Children: &[]*ast.AST{
			{
				Rule:     ast.StringLiteral,
				Value:    tokens[1].Value,
				Line:     tokens[1].Line,
				Column:   tokens[1].Column,
				File:     &tokens[1].File,
				Children: &[]*ast.AST{},
			},
		},
		Line:   tokens[0].Line,
		Column: tokens[0].Column,
		File:   &tokens[0].File,
	}

	return &child, nil
}
