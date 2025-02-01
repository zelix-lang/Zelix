/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package conditional

import (
	"fluent/ast"
	"fluent/parser/error"
	"fluent/parser/queue"
	"fluent/token"
)

// ProcessElse processes an else block in the Fluent programming language.
// It creates an AST node for the else block and schedules the tokens for further processing.
//
// Parameters:
// - block: A slice of tokens representing the else block.
// - blockQueue: A pointer to a slice of queue elements for scheduling tokens.
//
// Returns:
// - An AST node representing the else block.
// - An error object if there is an issue with processing the else block.
func ProcessElse(block []token.Token, blockQueue *[]queue.Element) (ast.AST, error.Error) {
	firstToken := block[0]

	// Check the input length
	if len(block) < 2 {
		return ast.AST{}, error.Error{
			Line:     firstToken.Line,
			Column:   firstToken.Column,
			File:     &firstToken.File,
			Expected: []ast.Rule{ast.Else},
		}
	}

	// Strip the first token
	block = block[1:]
	firstToken = block[0]

	// Create a new else node
	elseNode := ast.AST{
		Rule:     ast.Else,
		File:     &firstToken.File,
		Children: &[]*ast.AST{},
		Line:     firstToken.Line,
		Column:   firstToken.Column,
	}

	// Create a new block node
	blockNode := ast.AST{
		Rule:     ast.Block,
		File:     &firstToken.File,
		Children: &[]*ast.AST{},
		Line:     firstToken.Line,
		Column:   firstToken.Column,
	}

	// Append the block node to the else node
	*elseNode.Children = append(*elseNode.Children, &blockNode)

	// Schedule the tokens
	*blockQueue = append(*blockQueue, queue.Element{
		Parent: &blockNode,
		Tokens: block[1:], // Exclude the open curly
	})

	return elseNode, error.Error{}
}
