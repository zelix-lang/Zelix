/*
   The Fluent Programming Language
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
	"fluent/parser/rule/expression"
	"fluent/parser/util"
	"fluent/token"
)

// ProcessDeclaration processes a declaration from a list of tokens.
// It expects the input to contain at least 6 tokens and splits the tokens by the equal sign.
// If the input is valid, it processes the variable and its value, returning the resulting AST and any parsing errors.
//
// Parameters:
//   - input: a slice of token.Token representing the tokens to be processed.
//
// Returns:
//   - ast.AST: the resulting abstract syntax tree of the declaration.
//   - error.Error: any error encountered during parsing.
func ProcessDeclaration(input []token.Token) (ast.AST, error.Error) {
	if len(input) < 6 {
		return ast.AST{}, error.Error{
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Expected: []ast.Rule{ast.Declaration},
		}
	}

	split := util.SplitTokens(
		input,
		token.Assign,
		make([]token.Type, 0),
		make([]token.Type, 0),
	)

	if split == nil {
		return ast.AST{}, error.Error{
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Expected: []ast.Rule{ast.Declaration},
		}
	}

	if len(split) != 2 {
		return ast.AST{}, error.Error{
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Expected: []ast.Rule{ast.Declaration},
		}
	}

	variable, parsingError := ProcessIncompleteDeclaration(split[0])

	if parsingError.IsError() {
		return ast.AST{}, parsingError
	}

	variable.Rule = ast.Declaration

	value, parsingError := expression.ProcessExpression(split[1])

	if parsingError.IsError() {
		return ast.AST{}, parsingError
	}

	*variable.Children = append(*variable.Children, &value)

	return variable, error.Error{}
}
