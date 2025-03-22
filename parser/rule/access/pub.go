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

package access

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/rule/function"
	"fluent/parser/rule/mod"
	"fluent/token"
)

// ProcessPub processes a 'pub' declaration in the input tokens.
// It checks if the declaration is a function or module and delegates
// the processing to the appropriate parser.
//
// Parameters:
//   - input: A slice of token.Token representing the input tokens.
//
// Returns:
//   - ast.AST: The abstract syntax tree resulting from the parsing.
//   - error.Error: An error object if the parsing fails.
func ProcessPub(input []token.Token) (*ast.AST, *error.Error) {
	// Check the input
	if len(input) < 6 {
		return nil, &error.Error{
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Expected: []ast.Rule{ast.Function, ast.Module},
		}
	}

	// Check the 2nd token to determine if it is a function or module
	declaration := input[1].TokenType

	if declaration != token.Function && declaration != token.Mod {
		return nil, &error.Error{
			Line:     input[1].Line,
			Column:   input[1].Column,
			File:     &input[1].File,
			Expected: []ast.Rule{ast.Function, ast.Module},
		}
	}

	// Pass the tokens to the function or module parser
	if declaration == token.Function {
		// Exclude the 'pub' keyword
		return function.ProcessFunction(input[1:], true)
	}

	// Exclude the 'pub' keyword
	return mod.ProcessMod(input[1:], true)
}
