/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package signed

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/queue"
	"fluent/parser/util"
	"fluent/token"
)

// Used to avoid building this string every time it is needed
var dummyOrString = "||"

// ProcessSignedOp processes a signed operation in the input tokens.
// It handles negation, operators, and expressions, and constructs an AST.
//
// Parameters:
// - input: A slice of tokens to be processed.
// - expressionRule: The rule representing an expression.
// - operatorRule: The rule representing an operator.
// - allowedSigns: A map of allowed token types for operators.
// - allowedSignsSlice: A slice of allowed token types for operators.
// - candidate: A pointer to an AST node representing the candidate expression.
// - exprQueue: A pointer to a slice of queue elements for expressions.
// - hasSkippedTokens: A boolean indicating if tokens have been skipped.
// - isBoolean: A boolean indicating if the expression is a boolean.
// - includeSigns: A boolean indicating if the signs should be included.
//
// Returns:
// - An AST representing the processed expression.
// - An error if the processing fails.
func ProcessSignedOp(
	input []token.Token,
	expressionRule ast.Rule,
	operatorRule ast.Rule,
	allowedSigns map[token.Type]struct{},
	allowedSignsSlice []token.Type,
	candidate *ast.AST,
	exprQueue *[]queue.Element,
	hasSkippedTokens bool,
	isBoolean bool,
	includeSigns bool,
) (ast.AST, error.Error) {
	// Check the input length
	if len(input) < 1 {
		return ast.AST{}, error.Error{
			Line:     candidate.Line,
			Column:   candidate.Column,
			File:     candidate.File,
			Expected: []ast.Rule{expressionRule},
		}
	}

	result := ast.AST{
		Rule:     expressionRule,
		Line:     input[0].Line,
		Column:   input[0].Column,
		File:     &input[0].File,
		Children: &[]*ast.AST{},
	}

	if candidate != nil && candidate.Rule != ast.Program {
		*result.Children = append(*result.Children, &ast.AST{
			Rule:     ast.Expression,
			Line:     candidate.Line,
			Column:   candidate.Column,
			File:     candidate.File,
			Children: &[]*ast.AST{candidate},
		})
	}

	// The expression parser has processed a part of the expression
	// Therefore, we receive all tokens from the operator and so on
	expectingOperator := hasSkippedTokens
	inNegation := false

	// Used to skip tokens
	skip := 0

	for i, unit := range input {
		if i < skip {
			continue
		}

		if unit.TokenType == token.Not {
			if !isBoolean {
				return ast.AST{}, error.Error{
					Line:     unit.Line,
					Column:   unit.Column,
					File:     &unit.File,
					Expected: []ast.Rule{expressionRule},
				}
			}

			if expectingOperator && hasSkippedTokens {
				return ast.AST{}, error.Error{
					Line:     unit.Line,
					Column:   unit.Column,
					File:     &unit.File,
					Expected: []ast.Rule{expressionRule},
				}
			}

			hasSkippedTokens = true
			inNegation = true
			expectingOperator = false

			// Add the operator to the result
			*result.Children = append(*result.Children, &ast.AST{
				Rule:     operatorRule,
				Line:     unit.Line,
				Column:   unit.Column,
				File:     &unit.File,
				Value:    &unit.Value,
				Children: &[]*ast.AST{},
			})

			continue
		}

		if expectingOperator {
			if _, ok := allowedSigns[unit.TokenType]; !ok {
				return ast.AST{}, error.Error{
					Line:     unit.Line,
					Column:   unit.Column,
					File:     &unit.File,
					Expected: []ast.Rule{operatorRule},
				}
			}

			if includeSigns {
				// Add the operator to the result
				*result.Children = append(*result.Children, &ast.AST{
					Rule:     operatorRule,
					Line:     unit.Line,
					Column:   unit.Column,
					File:     &unit.File,
					Value:    &unit.Value,
					Children: &[]*ast.AST{},
				})
			}

			// Reset the flag
			expectingOperator = false
			continue
		}

		// Extract the tokens before the next operator
		extracted := util.ExtractTokensBefore(
			input[i:],
			allowedSignsSlice,
			true,
			token.OpenParen,
			token.CloseParen,
			false,
		)

		if extracted == nil {
			return ast.AST{}, error.Error{
				Line:     unit.Line,
				Column:   unit.Column,
				File:     &unit.File,
				Expected: []ast.Rule{expressionRule},
			}
		}

		firstElement := extracted[0]
		// Increment the skip
		skip = i + len(extracted)

		// Add the tokens to the expression queue
		expressionNode := ast.AST{
			Rule:     ast.Expression,
			Line:     firstElement.Line,
			Column:   firstElement.Column,
			File:     &firstElement.File,
			Children: &[]*ast.AST{},
		}

		if isBoolean {
			// Skip the tokens by OR for precedence
			split := util.SplitTokens(
				extracted,
				token.Or,
				[]token.Type{token.OpenParen},
				[]token.Type{token.CloseParen},
			)
			splitLen := len(split)

			// Used to handle nested boolean expressions
			var pushTo *ast.AST

			if splitLen > 1 {
				pushTo = &ast.AST{
					Rule:     ast.BooleanExpression,
					Line:     firstElement.Line,
					Column:   firstElement.Column,
					File:     &firstElement.File,
					Children: &[]*ast.AST{},
				}

				*expressionNode.Children = append(*expressionNode.Children, pushTo)
			}

			for j, tokens := range split {
				// Create a nested expression
				nestedExpression := ast.AST{
					Rule:     ast.Expression,
					Line:     firstElement.Line,
					Column:   firstElement.Column,
					File:     &firstElement.File,
					Children: &[]*ast.AST{},
				}

				*exprQueue = append(*exprQueue, queue.Element{
					Tokens: tokens,
					Parent: &nestedExpression,
				})

				// Append the nested expression to the result
				if pushTo != nil {
					*pushTo.Children = append(*pushTo.Children, &nestedExpression)
				} else {
					*expressionNode.Children = append(*expressionNode.Children, &nestedExpression)
				}

				// Add the OR operator if it is not the last expression
				if j < splitLen-1 {
					orRule := ast.AST{
						Rule:     operatorRule,
						Line:     firstElement.Line,
						Column:   firstElement.Column,
						File:     &firstElement.File,
						Value:    &dummyOrString,
						Children: &[]*ast.AST{},
					}

					if pushTo != nil {
						*pushTo.Children = append(*pushTo.Children, &orRule)
					} else {
						*expressionNode.Children = append(*expressionNode.Children, &orRule)
					}
				}
			}
		} else {
			*exprQueue = append(*exprQueue, queue.Element{
				Tokens: extracted,
				Parent: &expressionNode,
			})
		}

		// Append the expression to the result
		*result.Children = append(*result.Children, &expressionNode)

		// Set the flags
		inNegation = false
		expectingOperator = true
	}

	// The loop should never end not expecting an operator
	if !expectingOperator || inNegation {
		return ast.AST{}, error.Error{
			Line:     input[len(input)-1].Line,
			Column:   input[len(input)-1].Column,
			File:     &input[len(input)-1].File,
			Expected: []ast.Rule{expressionRule},
		}
	}

	// Used to skip tokens
	return result, error.Error{}
}
