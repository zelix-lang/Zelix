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
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/rule/identifier"
	_type "fluent/parser/rule/type"
	"fluent/token"
)

// ProcessIncompleteDeclaration processes an incomplete declaration from the given tokens.
// It expects a minimum of 4 tokens: a declaration type (let/const), an identifier, a colon, and a type.
// If the input is valid, it returns an AST representing the incomplete declaration.
// If the input is invalid, it returns an error.
//
// Parameters:
//   - input: a slice of tokens representing the declaration.
//
// Returns:
//   - ast.AST: the abstract syntax tree representing the incomplete declaration.
//   - error.Error: an error object if the input is invalid.
func ProcessIncompleteDeclaration(input []token.Token) (*ast.AST, *error.Error) {
	// Check the input's length
	if len(input) < 4 { // Minimum length: let x: T;
		return nil, &error.Error{
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Expected: []ast.Rule{ast.Declaration},
		}
	}

	result := ast.AST{
		Rule:     ast.IncompleteDeclaration,
		Line:     input[0].Line,
		Column:   input[0].Column,
		File:     &input[0].File,
		Children: &[]*ast.AST{},
	}

	// Add the declaration type
	switch input[0].TokenType {
	case token.Let, token.Const:
		// Create a declaration type node
		declarationType := ast.AST{
			Rule:     ast.DeclarationType,
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Children: &[]*ast.AST{},
			Value:    &input[0].Value,
		}

		// Add the declaration type to the result
		*result.Children = append(*result.Children, &declarationType)
	default:
		return nil, &error.Error{
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Expected: []ast.Rule{ast.Declaration},
		}
	}

	// The 2nd element must be an identifier
	name, parsingError := identifier.ProcessIdentifier(&input[1])

	if parsingError != nil {
		return nil, parsingError
	}

	// Add the identifier to the result
	*result.Children = append(*result.Children, name)

	// The 3rd element must be a colon
	if input[2].TokenType != token.Colon {
		return nil, &error.Error{
			Line:     input[2].Line,
			Column:   input[2].Column,
			File:     &input[2].File,
			Expected: []ast.Rule{ast.Declaration},
		}
	}

	// Extract the type
	typeTokens := input[3:]

	// Call the type parser
	typeAST, parsingError := _type.ProcessType(typeTokens, input[0])

	if parsingError != nil {
		return nil, parsingError
	}

	// Add the type to the result
	*result.Children = append(*result.Children, typeAST)

	// Return the result
	return &result, nil
}
