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

package array

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/queue"
	"fluent/parser/util"
	"fluent/token"
)

// ProcessArray processes an array of tokens and returns an AST representation of the array.
// It also queues elements for further processing.
//
// Parameters:
// - input: A slice of tokens representing the array.
// - exprQueue: A pointer to a slice of queue elements for further processing.
//
// Returns:
// - ast.AST: The AST representation of the array.
// - error.Error: An error object if the processing fails.
func ProcessArray(input []token.Token, exprQueue *[]queue.Element) (*ast.AST, *error.Error) {
	// Check the input's length
	if len(input) < 2 {
		return nil, &error.Error{
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Expected: []ast.Rule{ast.Array},
		}
	}

	// Check the last token
	if input[len(input)-1].TokenType != token.CloseBracket {
		return nil, &error.Error{
			Line:     input[len(input)-1].Line,
			Column:   input[len(input)-1].Column,
			File:     &input[len(input)-1].File,
			Expected: []ast.Rule{ast.Array},
		}
	}

	result := ast.AST{
		Rule:     ast.Array,
		Line:     input[0].Line,
		Column:   input[0].Column,
		File:     &input[0].File,
		Children: &[]*ast.AST{},
	}

	// Check for empty arrays
	if len(input) == 2 {
		return &result, nil
	}

	// Split the elements by commas
	elements := util.SplitTokens(
		input[1:len(input)-1],
		token.Comma,
		[]token.Type{token.OpenBracket, token.GreaterThan, token.OpenParen},
		[]token.Type{token.CloseBracket, token.LessThan, token.CloseParen},
	)

	if elements == nil {
		return nil, &error.Error{
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Expected: []ast.Rule{ast.Expression},
		}
	}

	// Queue all elements for processing
	for _, element := range elements {
		// Create a new expression node
		expressionNode := ast.AST{
			Rule:     ast.Expression,
			Line:     element[0].Line,
			Column:   element[0].Column,
			File:     &element[0].File,
			Children: &[]*ast.AST{},
		}

		*exprQueue = append(*exprQueue, queue.Element{
			Tokens: element,
			Parent: &expressionNode,
		})

		// Append the expression to the result
		*result.Children = append(*result.Children, &expressionNode)
	}

	return &result, nil
}
