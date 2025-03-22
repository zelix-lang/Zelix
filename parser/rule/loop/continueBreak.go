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
	"fluent/ast"
	"fluent/parser/error"
	"fluent/token"
)

// ProcessContinueOrBreak processes a continue or break statement and returns the corresponding AST node.
// Parameters:
// - statement: a slice of tokens representing the statement.
// Returns:
// - ast.AST: the generated AST node.
// - error.Error: an error object if the statement is invalid.
func ProcessContinueOrBreak(statement []token.Token) (*ast.AST, *error.Error) {
	// Get the first token
	el := statement[0]

	var rule ast.Rule

	if el.TokenType == token.Continue {
		rule = ast.Continue
	} else {
		rule = ast.Break
	}

	// Create a new AST node
	node := ast.AST{
		Rule:     rule,
		Line:     el.Line,
		Column:   el.Column,
		File:     &el.File,
		Children: &[]*ast.AST{},
	}

	// Ensure there are no tokens after the keyword
	if len(statement) > 1 {
		return nil, &error.Error{
			Line:     el.Line,
			Column:   el.Column,
			File:     &el.File,
			Expected: []ast.Rule{ast.Statement},
		}
	}

	// Return the AST node
	return &node, nil
}
