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

package template

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/rule/identifier"
	"fluent/token"
)

// ProcessTemplates parses a list of tokens into an AST representing templates.
// It expects the input to contain at least two tokens, with the last token being a GreaterThan token.
// If the input is invalid or a parsing error occurs, it returns an error.
//
// Parameters:
//   - input: A slice of tokens to be parsed.
//
// Returns:
//   - ast.AST: The resulting abstract syntax tree.
//   - error.Error: An error object if a parsing error occurs.
func ProcessTemplates(input []token.Token) (ast.AST, error.Error) {
	// Check the input length
	if len(input) < 2 { // The GreaterThan token (at the end) is not included
		return ast.AST{}, error.Error{
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Expected: []ast.Rule{ast.Templates},
		}
	}

	result := ast.AST{
		Rule:     ast.Templates,
		Line:     input[0].Line,
		Column:   input[0].Column,
		File:     &input[0].File,
		Children: &[]*ast.AST{},
	}

	expectingComma := false

	// Iterate over the tokens to parse the templates
	for i := 1; i < len(input); i++ {
		// Check for commas
		if expectingComma {
			if input[i].TokenType != token.Comma {
				return ast.AST{}, error.Error{
					Line:     input[i].Line,
					Column:   input[i].Column,
					File:     &input[i].File,
					Expected: []ast.Rule{ast.Templates},
				}
			}

			// Reset the flag
			expectingComma = false
			continue
		}

		// Extract the identifier
		template, parsingError := identifier.ProcessIdentifier(&input[i])

		// Check for parsing errors
		if parsingError.IsError() {
			return ast.AST{}, parsingError
		}

		// Append the template to the result
		*result.Children = append(*result.Children, &template)

		// Set the flag to expect a comma
		expectingComma = true
	}

	// The loop should never end with the flag set to false
	if !expectingComma {
		return ast.AST{}, error.Error{
			Line:     input[len(input)-1].Line,
			Column:   input[len(input)-1].Column,
			File:     &input[len(input)-1].File,
			Expected: []ast.Rule{ast.Templates},
		}
	}

	return result, error.Error{}
}
