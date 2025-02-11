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

package argument

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/rule/identifier"
	_type "fluent/parser/rule/type"
	"fluent/token"
)

// ProcessArguments parses a list of tokens representing function arguments
// and returns an AST (Abstract Syntax Tree) and an error if any.
// The function expects the input to be a list of tokens with at least three elements.
// It processes each token to extract argument names and types, handling nested types
// and ensuring proper syntax with colons and commas.
//
// Parameters:
// - input: A slice of token.Token representing the tokens to be parsed.
//
// Returns:
// - ast.AST: The resulting abstract syntax tree representing the parsed arguments.
// - error.Error: An error object containing details if parsing fails.
func ProcessArguments(input []token.Token) (ast.AST, error.Error) {
	// Check the input length
	if len(input) < 3 {
		return ast.AST{}, error.Error{
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Expected: []ast.Rule{ast.Parameter},
		}
	}

	result := ast.AST{
		Rule:     ast.Parameters,
		Line:     input[0].Line,
		Column:   input[0].Column,
		File:     &input[0].File,
		Children: &[]*ast.AST{},
	}

	// Used to know the level of nesting
	// in the argument list
	nestingLevel := 0

	// Current state of the parser
	expectingName := true
	expectingColon := false
	allowMore := true

	// Stores the argument's data
	var currentArgName string
	var currentArgLine int
	var currentArgColumn int
	var currentArgTypeTokens []token.Token

	inputLen := len(input)
	for i, unit := range input {
		// Check if we are expecting more arguments
		if !allowMore {
			return ast.AST{}, error.Error{
				Line:     unit.Line,
				Column:   unit.Column,
				File:     &unit.File,
				Expected: []ast.Rule{ast.Type},
			}
		}

		// Extract the argument name
		if expectingName {
			// Pass the unit to the identifier parser
			name, err := identifier.ProcessIdentifier(&unit)

			if err.IsError() {
				return ast.AST{}, err
			}

			// Set the current argument name
			currentArgName = *name.Value
			currentArgLine = name.Line
			currentArgColumn = name.Column
			expectingColon = true
			expectingName = false
			continue
		}

		// Check colons
		if expectingColon {
			if unit.TokenType != token.Colon {
				return ast.AST{}, error.Error{
					Line:     unit.Line,
					Column:   unit.Column,
					File:     &unit.File,
					Expected: []ast.Rule{ast.Type},
				}
			}

			expectingColon = false
			continue
		}

		if unit.TokenType == token.LessThan {
			nestingLevel++
		} else if unit.TokenType == token.GreaterThan {
			nestingLevel--

			if nestingLevel < 0 {
				return ast.AST{}, error.Error{
					Line:     unit.Line,
					Column:   unit.Column,
					File:     &unit.File,
					Expected: []ast.Rule{ast.Parameters},
				}
			}
		}

		// Detect the end of the argument
		if nestingLevel == 0 && (unit.TokenType == token.Comma || i == inputLen-1) {
			allowMore = i != inputLen-1

			// Append the argument's type to the result
			if !allowMore {
				currentArgTypeTokens = append(currentArgTypeTokens, unit)
			}

			argument := ast.AST{
				Rule:     ast.Parameter,
				Line:     currentArgLine,
				Column:   currentArgColumn,
				File:     &unit.File,
				Children: &[]*ast.AST{},
			}

			// Clone the value of the argument's name
			// to avoid memory issues
			clonedName := currentArgName

			// Push the argument's name to the result
			*argument.Children = append(*argument.Children, &ast.AST{
				Rule:     ast.Identifier,
				Value:    &clonedName,
				Line:     currentArgLine,
				Column:   currentArgColumn,
				Children: &[]*ast.AST{},
			})

			// Parse the argument's type and push it to the result
			argType, err := _type.ProcessType(currentArgTypeTokens, unit)

			if err.IsError() {
				return ast.AST{}, err
			}

			*argument.Children = append(*argument.Children, &argType)
			*result.Children = append(*result.Children, &argument)

			// Clear the current argument data
			currentArgName = ""
			currentArgTypeTokens = []token.Token{}
			expectingName = true
			continue
		}

		// Extract the argument type
		currentArgTypeTokens = append(currentArgTypeTokens, unit)
	}

	if nestingLevel > 0 {
		return ast.AST{}, error.Error{
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Expected: []ast.Rule{ast.Parameters},
		}
	}

	return result, error.Error{}
}
