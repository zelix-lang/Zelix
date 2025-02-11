package queue

import (
	"fluent/ast"
	"fluent/token"
)

// Element represents an element in the queue with associated tokens and a parent AST node.
type Element struct {
	// Tokens is a slice of Token objects associated with this queue element.
	Tokens []token.Token
	// Parent is a pointer to the parent AST node.
	Parent *ast.AST
}
