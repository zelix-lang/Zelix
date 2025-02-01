/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package reassignment

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/queue"
	"fluent/parser/util"
	"fluent/token"
)

// FindAndProcessReassignment identifies and processes reassignment operations in the input tokens.
// It returns an AST node representing the assignment, an error if any, and a boolean indicating success.
//
// Parameters:
// - input: A slice of tokens to be processed.
// - exprQueue: A pointer to a queue of elements for further processing.
//
// Returns:
// - ast.AST: The resulting AST node for the assignment.
// - error.Error: An error object if any error occurs.
// - bool: A boolean indicating whether the reassignment was successfully processed.
func FindAndProcessReassignment(input []token.Token, exprQueue *[]queue.Element) (ast.AST, error.Error, bool) {
	firstToken := input[0]
	firstTokenType := firstToken.TokenType

	// Check if the input has reassignment tokens
	if firstTokenType != token.Let && firstTokenType != token.Const && !util.TokenSliceContains(input, map[token.Type]struct{}{token.Assign: {}}) {
		return ast.AST{}, error.Error{}, false
	}

	// Split the tokens by reassignment operators
	split := util.SplitTokens(
		input,
		token.Assign,
		make([]token.Type, 0),
		make([]token.Type, 0),
	)

	// The split must have 2 elements
	if len(split) != 2 {
		return ast.AST{}, error.Error{
			Line:     firstToken.Line,
			Column:   firstToken.Column,
			File:     &firstToken.File,
			Expected: []ast.Rule{ast.Assignment},
		}, false
	}

	// Create a new assignment node
	assignmentNode := ast.AST{
		Rule:     ast.Assignment,
		Line:     firstToken.Line,
		Column:   firstToken.Column,
		File:     &firstToken.File,
		Children: &[]*ast.AST{},
	}

	// Create an expression node for both sides of the assigment
	expressionLeft := ast.AST{
		Rule:     ast.Expression,
		Line:     firstToken.Line,
		Column:   firstToken.Column,
		File:     &firstToken.File,
		Children: &[]*ast.AST{},
	}

	expressionRight := ast.AST{
		Rule:     ast.Expression,
		Line:     firstToken.Line,
		Column:   firstToken.Column,
		File:     &firstToken.File,
		Children: &[]*ast.AST{},
	}

	// Queue both sides
	*exprQueue = append(*exprQueue, queue.Element{
		Parent: &expressionLeft,
		Tokens: split[0],
	})

	*exprQueue = append(*exprQueue, queue.Element{
		Parent: &expressionRight,
		Tokens: split[1],
	})

	// Push both sides to the assignment node
	*assignmentNode.Children = append(*assignmentNode.Children, &expressionLeft)
	*assignmentNode.Children = append(*assignmentNode.Children, &expressionRight)
	return assignmentNode, error.Error{}, true
}
