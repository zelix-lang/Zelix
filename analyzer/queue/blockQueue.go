/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package queue

import "fluent/ast"

// BlockQueueElement represents an element in the block queue.
type BlockQueueElement struct {
	// Block is the block that is being stored in the queue.
	Block *ast.AST
	// ID is the scope ID that the block belongs to.
	ID int
}
