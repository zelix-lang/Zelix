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

package loop

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/queue"
	"fluent/parser/rule/expression"
	"fluent/parser/rule/identifier"
	"fluent/parser/util"
	"fluent/token"
)

// ProcessForLoop processes a 'for' loop in the Fluent programming language.
// It parses the loop expression and the block of code to be executed within the loop.
//
// Parameters:
// - input: A slice of tokens representing the 'for' loop.
// - blockQueue: A pointer to a slice of queue elements to which the loop block will be added.
//
// Returns:
// - An AST node representing the 'for' loop.
// - An error object if the parsing fails.
func ProcessForLoop(input []token.Token, blockQueue *[]queue.Element) (*ast.AST, *error.Error) {
	// Check the input's length
	if len(input) < 7 {
		return nil, &error.Error{
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Expected: []ast.Rule{ast.For},
		}
	}

	result := ast.AST{
		Rule:     ast.For,
		File:     &input[0].File,
		Children: &[]*ast.AST{},
		Line:     input[0].Line,
		Column:   input[0].Column,
	}

	// Strip the first token
	input = input[1:]

	// Get the loop expression
	loopExpression := util.ExtractTokensBefore(
		input,
		[]token.Type{token.OpenCurly},
		false,
		token.Unknown,
		token.Unknown,
		true,
	)

	if loopExpression == nil {
		return nil, &error.Error{
			Line:     input[0].Line,
			Column:   input[0].Column,
			File:     &input[0].File,
			Expected: []ast.Rule{ast.For},
		}
	}

	// Split the tokens by an "in" keyword
	loopExpressionSplit := util.SplitTokens(
		loopExpression,
		token.In,
		[]token.Type{},
		[]token.Type{},
	)

	// Check the loop expression's length
	if len(loopExpressionSplit) != 2 {
		return nil, &error.Error{
			Line:     loopExpression[0].Line,
			Column:   loopExpression[0].Column,
			File:     &loopExpression[0].File,
			Expected: []ast.Rule{ast.For},
		}
	}

	// Split the 1st item of the loop expression to get the range expression
	rangeExpressionSplit := util.SplitTokens(
		loopExpressionSplit[0],
		token.To,
		[]token.Type{},
		[]token.Type{},
	)

	// Check the range expression's length
	if len(rangeExpressionSplit) != 2 {
		return nil, &error.Error{
			Line:     loopExpressionSplit[0][0].Line,
			Column:   loopExpressionSplit[0][0].Column,
			File:     &loopExpressionSplit[0][0].File,
			Expected: []ast.Rule{ast.For},
		}
	}

	// Invoke the expression parser to parse the range expression
	rangeExpressionLeftNode, parsingError := expression.ProcessExpression(rangeExpressionSplit[0])

	if parsingError != nil {
		return nil, parsingError
	}

	rangeExpressionRightNode, parsingError := expression.ProcessExpression(rangeExpressionSplit[1])

	if parsingError != nil {
		return nil, parsingError
	}

	// Append the range expression to the result's children
	*result.Children = append(*result.Children, rangeExpressionLeftNode)
	*result.Children = append(*result.Children, rangeExpressionRightNode)

	// The 2nd item is the loop expression should always be an identifier
	if len(loopExpressionSplit[1]) != 1 {
		return nil, &error.Error{
			Line:     loopExpressionSplit[1][1].Line,
			Column:   loopExpressionSplit[1][1].Column,
			File:     &loopExpressionSplit[1][1].File,
			Expected: []ast.Rule{ast.Identifier},
		}
	}

	id, parsingError := identifier.ProcessIdentifier(&loopExpressionSplit[1][0])

	if parsingError != nil {
		return nil, parsingError
	}

	// Append the identifier to the result's children
	*result.Children = append(*result.Children, id)

	// Get the actual block of the loop
	block := input[len(loopExpression)+1:]

	// Create a new block node
	blockNode := ast.AST{
		Rule:     ast.Block,
		File:     &input[0].File,
		Line:     input[0].Line,
		Column:   input[0].Column,
		Children: &[]*ast.AST{},
	}

	// Add the block to the block queue
	*blockQueue = append(*blockQueue, queue.Element{Tokens: block, Parent: &blockNode})

	// Append the block node to the result's children
	*result.Children = append(*result.Children, &blockNode)

	return &result, nil
}
