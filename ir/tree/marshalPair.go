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
	"fluent/filecode/types/wrapper"
)

// MarshalPair represents a pair of a child AST node and its parent instruction tree.
// It also includes metadata about the relationship.
type MarshalPair struct {
	Child    *ast.AST            // Child is the AST node.
	Parent   *InstructionTree    // Parent is the instruction tree containing the child.
	IsInline bool                // IsInline indicates if the pair is inline.
	Counter  int                 // Counter is used for tracking purposes.
	IsParam  bool                // IsParam indicates if the pair is a parameter.
	Expected wrapper.TypeWrapper // Expected is the expected type of the child node.
}
