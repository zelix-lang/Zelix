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

package tree

import (
	"fluent/ast"
	"strings"
)

// BlockMarshalElement represents an element in the block marshaling process.
type BlockMarshalElement struct {
	// Element is the AST node associated with this block.
	Element *ast.AST
	// Representation is a string builder for the block's representation.
	Representation *strings.Builder
	// ParentAddr is the address of the parent block.
	ParentAddr *string
	// RemainingAddr is the address of the remaining block.
	RemainingAddr *string
	// JumpToParent indicates whether to jump to the parent block.
	JumpToParent bool
	// Id is the identifier of the block.
	Id int
	// IsLast indicates whether this is the last block.
	IsLast bool
}
