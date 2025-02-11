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

package generic

import (
	"fluent/ast"
	"fluent/parser/error"
	_type "fluent/parser/rule/type"
	"fluent/parser/util"
	"fluent/token"
)

// ProcessGenerics processes a list of tokens representing generics and returns an AST and an error.
// It checks for inferred types and splits the tokens based on commas and angle brackets.
// It then iterates over the generics, parsing each type and appending it to the result's children.
//
// Parameters:
//   - input: A slice of tokens representing the generics.
//
// Returns:
//   - ast.AST: The abstract syntax tree representing the generics.
//   - error.Error: An error object if any parsing error occurs.
func ProcessGenerics(input []token.Token) (ast.AST, error.Error) {
	// Checked for inferred types
	if len(input) == 0 {
		return ast.AST{
			Rule:     ast.InferredType,
			Children: &[]*ast.AST{},
		}, error.Error{}
	}

	result := ast.AST{
		Rule:     ast.Generics,
		File:     &input[0].File,
		Line:     input[0].Line,
		Column:   input[0].Column,
		Children: &[]*ast.AST{},
	}

	// Split the tokens
	generics := util.SplitTokens(
		input,
		token.Comma,
		[]token.Type{token.LessThan},
		[]token.Type{token.GreaterThan},
	)

	// Iterate over the generics
	for _, generic := range generics {
		// Call the type parser
		genericType, parsingError := _type.ProcessType(generic, input[0])

		if parsingError.IsError() {
			return ast.AST{}, parsingError
		}

		// Append the generic type to the result's children
		*result.Children = append(*result.Children, &genericType)
	}

	return result, error.Error{}
}
