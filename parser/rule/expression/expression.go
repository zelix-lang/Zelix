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

package expression

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/queue"
	"fluent/parser/rule/array"
	"fluent/parser/rule/call"
	"fluent/parser/rule/identifier"
	"fluent/parser/rule/object"
	"fluent/parser/rule/signed"
	"fluent/parser/shared/arithmetic"
	"fluent/parser/shared/boolean"
	"fluent/parser/util"
	"fluent/token"
)

// processSingleToken processes a single token and returns an AST (Abstract Syntax Tree) node and an error.
// It handles different types of literals and identifiers.
//
// Parameters:
// - unit: A pointer to the token to be processed.
//
// Returns:
// - ast.AST: The resulting abstract syntax tree node.
// - error.Error: An error object if any parsing error occurs.
func processSingleToken(unit *token.Token) (*ast.AST, *error.Error) {
	switch unit.TokenType {
	case token.BoolLiteral:
		return &ast.AST{
			Rule:   ast.BooleanLiteral,
			Value:  unit.Value,
			Line:   unit.Line,
			Column: unit.Column,
			File:   &unit.File,
		}, nil
	case token.NumLiteral:
		return &ast.AST{
			Rule:   ast.NumberLiteral,
			Value:  unit.Value,
			Line:   unit.Line,
			Column: unit.Column,
			File:   &unit.File,
		}, nil
	case token.DecimalLiteral:
		return &ast.AST{
			Rule:   ast.DecimalLiteral,
			Value:  unit.Value,
			Line:   unit.Line,
			Column: unit.Column,
			File:   &unit.File,
		}, nil
	case token.StringLiteral:
		return &ast.AST{
			Rule:   ast.StringLiteral,
			Value:  unit.Value,
			Line:   unit.Line,
			Column: unit.Column,
			File:   &unit.File,
		}, nil
	// The only other case allowed is an identifier
	default:
		id, parsingError := identifier.ProcessIdentifier(unit)

		if parsingError != nil {
			return nil, parsingError
		}

		return id, nil
	}
}

// ProcessExpression processes a list of tokens and returns an AST (Abstract Syntax Tree) and an error.
// It uses a queue to avoid recursion and handles various types of expressions including single tokens,
// reassignments, pointers, dereferences, parentheses, negations, function calls, identifiers, booleans,
// arithmetic expressions, and property accesses.
//
// Parameters:
// - input: A slice of tokens to be processed.
//
// Returns:
// - ast.AST: The resulting abstract syntax tree.
// - error.Error: An error object if any parsing error occurs.
func ProcessExpression(input []token.Token) (*ast.AST, *error.Error) {
	// Use a queue to avoid recursion
	var processQueue []queue.Element

	// Check for empty input
	if len(input) == 0 {
		return nil, &error.Error{
			Line:     0,
			Column:   0,
			Expected: []ast.Rule{ast.Expression},
		}
	}

	result := ast.AST{
		Rule:     ast.Expression,
		Line:     input[0].Line,
		Column:   input[0].Column,
		File:     &input[0].File,
		Children: &[]*ast.AST{},
	}

	processQueue = append(processQueue, queue.Element{Tokens: input, Parent: &result})

	for len(processQueue) > 0 {
		// Pop the first element from the queue
		element := processQueue[0]
		processQueue = processQueue[1:]

		// Subtract the data from the element
		input := element.Tokens
		parent := element.Parent

		inputLen := len(input)

		// Check for empty input
		if inputLen == 0 {
			return nil, &error.Error{
				Line:     0,
				Column:   0,
				Expected: []ast.Rule{ast.Expression},
			}
		}

		startAt := 0

		// Parse pointers and dereferences
		for i, unit := range input {
			hasToBreak := false
			tokenType := unit.TokenType

			switch tokenType {
			case token.Asterisk, token.Ampersand, token.And:
				var operation ast.AST

				if tokenType == token.Asterisk {
					// Create a new dereference node
					operation = ast.AST{
						Rule:     ast.Dereference,
						Line:     unit.Line,
						Column:   unit.Column,
						File:     &unit.File,
						Children: &[]*ast.AST{},
					}
				} else {
					// Create a new pointer node
					operation = ast.AST{
						Rule:     ast.Pointer,
						Line:     unit.Line,
						Column:   unit.Column,
						File:     &unit.File,
						Children: &[]*ast.AST{},
					}
				}

				// Append the operation to the parent
				*parent.Children = append(*parent.Children, &operation)

				if tokenType == token.And {
					// Append a double pointer to the parent
					*parent.Children = append(*parent.Children, &operation)
				}

				startAt = i + 1
			default:
				hasToBreak = true
			}

			if hasToBreak {
				break
			}
		}

		if startAt > 0 {
			// Strip the tokens that were processed
			input = input[startAt:]
			inputLen = len(input)
			startAt = 0
		}

		// Check for arrays
		if input[0].TokenType == token.OpenBracket {
			// Pass the tokens to the array parser
			arr, parsingError := array.ProcessArray(input, &processQueue)

			if parsingError != nil {
				return nil, parsingError
			}

			// Add the array to the parent
			*parent.Children = append(*parent.Children, &ast.AST{
				Rule:     ast.Expression,
				Children: &[]*ast.AST{arr},
				Line:     arr.Line,
				Column:   arr.Column,
				File:     arr.File,
			})

			// Avoid further processing
			continue
		}

		// Check for single tokens
		if inputLen == 1 {
			child, parsingError := processSingleToken(&input[0])

			if parsingError != nil {
				return nil, parsingError
			}

			*parent.Children = append(*parent.Children, child)
			continue
		}

		// A candidate that can be later re-utilized if a property access is detected
		var candidate *ast.AST

		// Check for parentheses
		if input[0].TokenType == token.OpenParen {
			// Extract the tokens inside the parentheses
			nestedExpression := util.ExtractTokensBefore(
				input,
				[]token.Type{token.CloseParen},
				true,
				token.OpenParen,
				token.CloseParen,
				true,
			)

			if nestedExpression == nil {
				return nil, &error.Error{
					Line:     input[0].Line,
					Column:   input[0].Column,
					File:     &input[0].File,
					Expected: []ast.Rule{ast.Expression},
				}
			}

			// -1 to exclude the closing parenthesis
			if len(nestedExpression) == inputLen-1 {
				// Queue the nested expression
				processQueue = append(
					processQueue,
					queue.Element{
						// Strip the open parenthesis
						Tokens: nestedExpression[1:],
						Parent: parent,
					},
				)

				// Avoid processing the nested expression again
				continue
			} else {
				// Skip the tokens of the subexpression
				startAt = len(nestedExpression) + 1

				// Create a new expression node
				child := ast.AST{
					Rule:     ast.Expression,
					Line:     input[0].Line,
					Column:   input[0].Column,
					File:     &input[0].File,
					Children: &[]*ast.AST{},
				}

				candidate = &child

				// Queue the nested expression
				processQueue = append(
					processQueue,
					queue.Element{
						// Strip the open parenthesis
						Tokens: nestedExpression[1:],
						Parent: &child,
					},
				)
			}
		}

		// See if there are more tokens
		if startAt > 0 && startAt >= inputLen {
			continue
		}

		if startAt == 0 {
			// Check for negations
			if input[0].TokenType == token.Not {
				// Pass all the tokens directly to the signed op parser
				expression, parsingError := signed.ProcessSignedOp(
					input,
					ast.BooleanExpression,
					ast.BooleanOperator,
					boolean.Operators,
					boolean.OperatorsSlice,
					nil,
					&processQueue,
					false,
					true,
					true,
				)

				if parsingError != nil {
					return nil, parsingError
				}

				// Append the negation to the parent
				*parent.Children = append(*parent.Children, expression)
				continue
			}

			// Check for function calls or identifiers
			if input[0].TokenType == token.New || input[1].TokenType == token.OpenParen {
				// Extract the function call's tokens
				callTokens := util.ExtractTokensBefore(
					input,
					[]token.Type{token.CloseParen},
					true,
					token.OpenParen,
					token.CloseParen,
					true,
				)

				if callTokens == nil {
					return nil, &error.Error{
						Line:     input[0].Line,
						Column:   input[0].Column,
						File:     &input[0].File,
						Expected: []ast.Rule{ast.FunctionCall},
					}
				}

				startAt = len(callTokens) + 1
				// Append the closing parenthesis
				callTokens = append(callTokens, input[startAt-1])

				var nestedResult *ast.AST
				var parsingError *error.Error

				// Pass the tokens to the appropriate parser
				if input[0].TokenType == token.New {
					nestedResult, parsingError = object.ProcessObjectCreation(callTokens, &processQueue)
				} else {
					nestedResult, parsingError = call.ProcessFunctionCall(callTokens, &processQueue)
				}

				if parsingError != nil {
					return nil, parsingError
				}

				candidate = nestedResult
			} else {
				startAt = 1

				// Parse the first token
				child, parsingError := processSingleToken(&input[0])

				if parsingError != nil {
					return nil, parsingError
				}

				candidate = child
			}
		}

		// See if there are more tokens
		if startAt >= inputLen {
			// Push the candidate if needed
			if candidate != nil && candidate.Rule != ast.Program {
				*parent.Children = append(*parent.Children, candidate)
			}

			continue
		}

		// Get the token type of the next token
		remainingTokens := input[startAt:]
		nextTokenType := input[startAt].TokenType

		// Check for booleans
		if _, ok := boolean.Operators[nextTokenType]; ok {
			var expression *ast.AST
			var parsingError *error.Error

			if nextTokenType == token.Or {
				expression, parsingError = signed.ProcessSignedOp(
					input,
					ast.BooleanExpression,
					ast.BooleanOperator,
					boolean.Operators,
					boolean.OperatorsSlice,
					nil,
					&processQueue,
					false,
					true,
					true,
				)
			} else {
				expression, parsingError = signed.ProcessSignedOp(
					remainingTokens,
					ast.BooleanExpression,
					ast.BooleanOperator,
					boolean.Operators,
					boolean.OperatorsSlice,
					candidate,
					&processQueue,
					true,
					true,
					true,
				)
			}

			if parsingError != nil {
				return nil, parsingError
			}

			*parent.Children = append(*parent.Children, expression)
			continue
		}

		// Parse expressions like 2 + 2 == 4
		isBoolean := util.TokenSliceContains(remainingTokens, boolean.Operators)
		var tokensBeforeOperator []token.Token

		if isBoolean {
			tokensBeforeOperator = util.ExtractTokensBefore(
				remainingTokens,
				boolean.OperatorsSlice,
				false,
				token.Unknown,
				token.Unknown,
				true,
			)
		}

		// Check for arithmetic
		if _, ok := arithmetic.Operators[nextTokenType]; ok {
			// Edge case: Check for "a.b.c + 25"
			if element.IsPropAccess {
				// Get the result's children
				resultChildren := *result.Children
				// Get the last child
				lastChild := resultChildren[len(resultChildren)-1]
				// Delete the last child
				*result.Children = resultChildren[:len(resultChildren)-1]
				// Append the candidate to the last child
				lastChildExpr := (*lastChild.Children)[len(*lastChild.Children)-1]
				*lastChildExpr.Children = append(*lastChildExpr.Children, candidate)
				// Update the candidate
				candidate = lastChild
			}

			passedTokens := remainingTokens

			if isBoolean {
				passedTokens = tokensBeforeOperator
			}

			expression, parsingError := signed.ProcessSignedOp(
				passedTokens,
				ast.ArithmeticExpression,
				ast.ArithmeticSign,
				arithmetic.Operators,
				arithmetic.OperatorsSlice,
				candidate,
				&processQueue,
				true,
				false,
				true,
			)

			if parsingError != nil {
				return nil, parsingError
			}

			// Extend the queue
			if isBoolean {
				candidate = expression
			} else {
				if element.IsPropAccess {
					*result.Children = append(*result.Children, expression)
				} else {
					*parent.Children = append(*parent.Children, expression)
				}
			}

			continue
		}

		if isBoolean {
			// Reassign the remaining tokens
			remainingTokens = remainingTokens[len(tokensBeforeOperator):]

			// Parse the boolean expression
			expression, parsingError := signed.ProcessSignedOp(
				remainingTokens,
				ast.BooleanExpression,
				ast.BooleanOperator,
				boolean.Operators,
				boolean.OperatorsSlice,
				candidate,
				&processQueue,
				true,
				true,
				true,
			)

			if parsingError != nil {
				return nil, parsingError
			}

			*parent.Children = append(*parent.Children, expression)
			continue
		}

		if nextTokenType != token.Dot {
			return nil, &error.Error{
				Line:     input[startAt].Line,
				Column:   input[startAt].Column,
				File:     &input[startAt].File,
				Expected: []ast.Rule{ast.PropertyAccess},
			}
		}

		// Parse the property access
		propertyAccess, parsingError := signed.ProcessSignedOp(
			remainingTokens,
			ast.PropertyAccess,
			ast.PropertyAccess,
			map[token.Type]struct{}{token.Dot: {}},
			[]token.Type{token.Dot},
			candidate,
			&processQueue,
			true,
			false,
			false,
		)

		if parsingError != nil {
			return nil, parsingError
		}

		// Append the property access to the parent
		*parent.Children = append(*parent.Children, propertyAccess)
	}

	if len(*result.Children) == 1 && (*result.Children)[0].Rule == ast.Expression {
		return (*result.Children)[0], nil
	}

	return &result, nil
}
