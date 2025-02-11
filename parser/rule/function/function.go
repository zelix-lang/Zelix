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

package function

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/rule/argument"
	block2 "fluent/parser/rule/block"
	"fluent/parser/rule/identifier"
	"fluent/parser/rule/template"
	_type "fluent/parser/rule/type"
	"fluent/parser/util"
	"fluent/token"
)

// ProcessFunction parses a function from the given tokens.
// It handles public and private functions, function names, templates, parameters, return types, and function bodies.
//
// Parameters:
// - input: A slice of tokens representing the function.
// - public: A boolean indicating if the function is public.
//
// Returns:
// - ast.AST: The abstract syntax tree representation of the function.
// - error.Error: An error object if parsing fails.
func ProcessFunction(input []token.Token, public bool) (ast.AST, error.Error) {
	// Check the input length
	inputLen := len(input)

	if inputLen < 5 {
		trace := input[0]

		return ast.AST{}, error.Error{
			Line:     trace.Line,
			Column:   trace.Column,
			File:     &trace.File,
			Expected: []ast.Rule{ast.Function},
		}
	}

	result := ast.AST{
		Rule:     ast.Function,
		Line:     input[0].Line,
		Column:   input[0].Column,
		File:     &input[0].File,
		Children: &[]*ast.AST{},
	}

	// Push access modifiers if needed
	if public {
		*result.Children = append(*result.Children, &ast.AST{
			Rule:     ast.Public,
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Children: &[]*ast.AST{},
		})
	}

	// Parse the function's name
	name, err := identifier.ProcessIdentifier(&input[1])

	if err.IsError() {
		return result, err
	}

	// Add the name to the result
	*result.Children = append(*result.Children, &ast.AST{
		Rule:     ast.Identifier,
		Value:    name.Value,
		Line:     name.Line,
		Column:   name.Column,
		Children: &[]*ast.AST{},
	})

	// See if the function has templates
	startAt := 2

	if input[2].TokenType == token.LessThan {
		// Extract the tokens before the next GreaterThan
		extracted := util.ExtractTokensBefore(
			input[2:],
			[]token.Type{token.GreaterThan},
			false,
			token.Unknown,
			token.Unknown,
			true,
		)

		if extracted == nil {
			return ast.AST{}, error.Error{
				Line:     input[2].Line,
				Column:   input[2].Column,
				File:     &input[2].File,
				Expected: []ast.Rule{ast.Templates},
			}
		}

		// Increment lookForParenAt
		startAt += len(extracted) + 1

		// Parse the templates
		templates, err := template.ProcessTemplates(extracted)

		if err.IsError() {
			return ast.AST{}, err
		}

		// Append the templates to the result
		*result.Children = append(*result.Children, &templates)
	}

	// Check for the opening parenthesis
	if input[startAt].TokenType != token.OpenParen {
		return ast.AST{}, error.Error{
			Line:     input[startAt].Line,
			Column:   input[startAt].Column,
			File:     &input[startAt].File,
			Expected: []ast.Rule{ast.Function},
		}
	}

	// Check if the function includes parameters
	if input[startAt+1].TokenType != token.CloseParen {
		// Extract everything before the closing parenthesis
		extracted := util.ExtractTokensBefore(
			input[startAt+1:],
			[]token.Type{token.CloseParen},
			false,
			token.Unknown,
			token.Unknown,
			true,
		)

		if extracted == nil {
			return ast.AST{}, error.Error{
				Line:     input[startAt].Line,
				Column:   input[startAt].Column,
				File:     &input[startAt].File,
				Expected: []ast.Rule{ast.Function},
			}
		}

		// +1 to skip the closing parenthesis
		startAt += len(extracted) + 2

		// Parse the arguments
		arguments, err := argument.ProcessArguments(extracted)

		if err.IsError() {
			return ast.AST{}, err
		}

		// Append the arguments to the result
		*result.Children = append(*result.Children, &arguments)
	} else {
		startAt += 2
	}

	if startAt >= inputLen {
		return ast.AST{}, error.Error{
			Line:     input[startAt-1].Line,
			Column:   input[startAt-1].Column,
			File:     &input[startAt-1].File,
			Expected: []ast.Rule{ast.Function},
		}
	}

	// Parse return types
	if input[startAt].TokenType == token.Arrow {
		// Extract everything before the opening curly brace
		extracted := util.ExtractTokensBefore(
			input[startAt+1:],
			[]token.Type{token.OpenCurly},
			false,
			token.Unknown,
			token.Unknown,
			true,
		)

		if extracted == nil {
			return ast.AST{}, error.Error{
				Line:     input[startAt].Line,
				Column:   input[startAt].Column,
				File:     &input[startAt].File,
				Expected: []ast.Rule{ast.Function},
			}
		}

		startAt += len(extracted) + 1

		// Parse the return types
		returnType, err := _type.ProcessType(extracted, input[startAt-1])

		if err.IsError() {
			return ast.AST{}, err
		}

		// Append the return types to the result
		*result.Children = append(*result.Children, &returnType)
	}

	if startAt >= inputLen {
		return ast.AST{}, error.Error{
			Line:     input[startAt-1].Line,
			Column:   input[startAt-1].Column,
			File:     &input[startAt-1].File,
			Expected: []ast.Rule{ast.Function},
		}
	}

	// Check for block opening
	if input[startAt].TokenType != token.OpenCurly {
		return ast.AST{}, error.Error{
			Line:     input[startAt].Line,
			Column:   input[startAt].Column,
			File:     &input[startAt].File,
			Expected: []ast.Rule{ast.Function},
		}
	}

	// Parse the block directly
	block, err := block2.ProcessBlock(input[startAt+1:])

	if err.IsError() {
		return ast.AST{}, err
	}

	// Append the block to the result
	*result.Children = append(*result.Children, &block)
	return result, error.Error{}
}
