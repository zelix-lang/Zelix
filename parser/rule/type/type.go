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

package _type

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/queue"
	"fluent/parser/rule/identifier"
	"fluent/parser/util"
	"fluent/token"
)

// allowedIdentifiers is a map of token types that are allowed in a type.
var allowedIdentifiers = map[token.Type]struct{}{
	token.Identifier: {},
	token.String:     {},
	token.Num:        {},
	token.Bool:       {},
	token.Dec:        {},
	token.Nothing:    {},
}

// ProcessType processes a list of tokens to generate an abstract syntax tree (AST) for a type.
// It validates the tokens and handles different type constructs such as arrays, pointers, and templates.
//
// Parameters:
//   - input: A slice of tokens representing the type to be processed.
//   - trace: A token used for tracing the source location in case of errors.
//
// Returns:
//   - ast.AST: The generated abstract syntax tree for the type.
//   - error.Error: An error object containing details if the processing fails.
func ProcessType(input []token.Token, trace token.Token) (ast.AST, error.Error) {
	// Check the input's length
	if len(input) < 1 {
		return ast.AST{}, error.Error{
			Line:     trace.Line,
			Column:   trace.Column,
			File:     &trace.File,
			Expected: []ast.Rule{ast.Type},
		}
	}

	result := ast.AST{
		Rule:     ast.Type,
		Line:     input[0].Line,
		Column:   input[0].Column,
		File:     &input[0].File,
		Children: &[]*ast.AST{},
	}

	// Create a queue to process the tokens
	var processQueue []queue.Element

	// Add the first element to the queue
	processQueue = append(processQueue, queue.Element{Tokens: input, Parent: &result})

	for len(processQueue) > 0 {
		// Pop the first element
		element := processQueue[0]
		processQueue = processQueue[1:]

		input := element.Tokens
		result := element.Parent

		// Used to determine the validity of the type
		hasMetIdentifier := false

		// Used to determine what the parser is expecting
		expectingCloseBracket := false
		expectingOpenBracket := false

		// Used to know if the current type has pointers
		hasPointers := false

		// Used to check for invalid types
		baseTypeIsNothing := false

		// Used to skip indexes
		skip := 0

		for i, unit := range input {
			if i < skip {
				continue
			}

			// Parse array types
			if expectingOpenBracket || unit.TokenType == token.OpenBracket {
				if baseTypeIsNothing {
					return ast.AST{}, error.Error{
						Line:     unit.Line,
						Column:   unit.Column,
						File:     &unit.File,
						Expected: []ast.Rule{ast.Type},
					}
				}

				if !hasMetIdentifier || unit.TokenType != token.OpenBracket {
					return ast.AST{}, error.Error{
						Line:     unit.Line,
						Column:   unit.Column,
						File:     &unit.File,
						Expected: []ast.Rule{ast.ArrayType},
					}
				}

				// Push the array AST
				*result.Children = append(*result.Children, &ast.AST{
					Rule:     ast.ArrayType,
					Line:     unit.Line,
					Column:   unit.Column,
					File:     &unit.File,
					Children: &[]*ast.AST{},
				})

				// Set the flag
				expectingOpenBracket = false
				expectingCloseBracket = true
				continue
			}

			if expectingCloseBracket {
				if unit.TokenType != token.CloseBracket {
					return ast.AST{}, error.Error{
						Line:     unit.Line,
						Column:   unit.Column,
						File:     &unit.File,
						Expected: []ast.Rule{ast.ArrayType},
					}
				}

				// Reset the flag
				expectingCloseBracket = false
				expectingOpenBracket = true

				continue
			}

			if unit.TokenType == token.Ampersand || unit.TokenType == token.And {
				hasPointers = true
				if hasMetIdentifier {
					return ast.AST{}, error.Error{
						Line:     unit.Line,
						Column:   unit.Column,
						File:     &unit.File,
						Expected: []ast.Rule{ast.Type},
					}
				}

				// Push a pointer AST
				*result.Children = append(*result.Children, &ast.AST{
					Rule:     ast.Pointer,
					Line:     unit.Line,
					Column:   unit.Column,
					File:     &unit.File,
					Children: &[]*ast.AST{},
				})

				if unit.TokenType == token.And {
					// Add a pointer for the reference
					*result.Children = append(*result.Children, &ast.AST{
						Rule:     ast.Pointer,
						Line:     unit.Line,
						Column:   unit.Column,
						File:     &unit.File,
						Children: &[]*ast.AST{},
					})
				}

				continue
			}

			if _, ok := allowedIdentifiers[unit.TokenType]; ok {
				if hasMetIdentifier {
					return ast.AST{}, error.Error{
						Line:     unit.Line,
						Column:   unit.Column,
						File:     &unit.File,
						Expected: []ast.Rule{ast.Type},
					}
				}

				hasMetIdentifier = true

				// Check if the token is a primitive type or an identifier
				if unit.TokenType == token.Identifier {
					// Parse the identifier using the identifier parser
					child, _ := identifier.ProcessIdentifier(&unit)

					// Parse templates
					if i+2 < len(input) && input[i+1].TokenType == token.LessThan {
						// Check if the type is supposed to be inferred
						if input[i+2].TokenType == token.GreaterThan {
							*child.Children = append(*child.Children, &ast.AST{
								Rule:     ast.InferredType,
								Line:     input[i+2].Line,
								Column:   input[i+2].Column,
								File:     &input[i+2].File,
								Children: &[]*ast.AST{},
							})

							skip = i + 3
						} else {
							if i+3 > len(input) {
								return ast.AST{}, error.Error{
									Line:     input[i+2].Line,
									Column:   input[i+2].Column,
									File:     &input[i+2].File,
									Expected: []ast.Rule{ast.Templates},
								}
							}

							// Extract the tokens before the next GreaterThan
							extracted := util.ExtractTokensBefore(
								input[i+1:],
								[]token.Type{token.GreaterThan},
								true,
								token.LessThan,
								token.GreaterThan,
								true,
							)

							if extracted == nil {
								return ast.AST{}, error.Error{
									Line:     input[i+2].Line,
									Column:   input[i+2].Column,
									File:     &input[i+2].File,
									Expected: []ast.Rule{ast.Templates},
								}
							}

							// Strip the less than token
							extracted = extracted[1:]

							skip = i + len(extracted) + 4

							// Split the tokens by commas
							split := util.SplitTokens(
								extracted,
								token.Comma,
								[]token.Type{token.LessThan},
								[]token.Type{token.GreaterThan},
							)

							if split == nil {
								return ast.AST{}, error.Error{
									Line:     input[i+2].Line,
									Column:   input[i+2].Column,
									File:     &input[i+2].File,
									Expected: []ast.Rule{ast.Templates},
								}
							}

							// Parse the templates
							for _, tokens := range split {
								nestedType := ast.AST{
									Rule:     ast.Type,
									Line:     unit.Line,
									Column:   unit.Column,
									File:     &unit.File,
									Children: &[]*ast.AST{},
								}

								// Add a new element to the queue
								processQueue = append(processQueue, queue.Element{Tokens: tokens, Parent: &nestedType})

								// Append the nested type to the child
								*child.Children = append(*child.Children, &nestedType)
							}
						}
					}

					*result.Children = append(*result.Children, &child)
				} else {
					// Check for pointers to "nothing"
					baseTypeIsNothing = unit.TokenType == token.Nothing
					if baseTypeIsNothing && hasPointers {
						return ast.AST{}, error.Error{
							Line:     unit.Line,
							Column:   unit.Column,
							File:     &unit.File,
							Expected: []ast.Rule{ast.Type},
						}
					}

					// Push as primitive
					*result.Children = append(*result.Children, &ast.AST{
						Rule:     ast.Primitive,
						Line:     unit.Line,
						Column:   unit.Column,
						File:     &unit.File,
						Value:    &unit.Value,
						Children: &[]*ast.AST{},
					})
				}

				continue
			}

			// Invalid token
			return ast.AST{}, error.Error{
				Line:     unit.Line,
				Column:   unit.Column,
				File:     &unit.File,
				Expected: []ast.Rule{ast.Type},
			}
		}

		// Check if the type is valid
		if !hasMetIdentifier || expectingCloseBracket {
			return ast.AST{}, error.Error{
				Line:     input[0].Line,
				Column:   input[0].Column,
				File:     &input[0].File,
				Expected: []ast.Rule{ast.Identifier},
			}
		}
	}

	return result, error.Error{}
}
