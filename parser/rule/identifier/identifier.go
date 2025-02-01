/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package identifier

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/token"
)

// ProcessIdentifier processes a token and returns an AST node if the token is an identifier.
// If the token is not an identifier, it returns an error.
//
// Parameters:
// - unit: A pointer to the token to be processed.
//
// Returns:
// - ast.AST: The AST node representing the identifier.
// - error.Error: An error object if the token is not an identifier.
func ProcessIdentifier(unit *token.Token) (ast.AST, error.Error) {
	// See if the token is an identifier
	if unit.TokenType != token.Identifier {
		return ast.AST{}, error.Error{
			Line:     unit.Line,
			Column:   unit.Column,
			File:     &unit.File,
			Expected: []ast.Rule{ast.Identifier},
		}
	}

	// Create the identifier AST
	return ast.AST{
		Rule:     ast.Identifier,
		Value:    &unit.Value,
		Line:     unit.Line,
		Column:   unit.Column,
		File:     &unit.File,
		Children: &[]*ast.AST{},
	}, error.Error{}
}
