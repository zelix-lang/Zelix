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

package value

import (
	"fluent/ast"
	"fluent/ir/tree"
)

// WriteBoolLiteral writes a boolean literal to the parent InstructionTree.
// It writes "__TRUE" if the value is true, "__FALSE" otherwise.
//
// Parameters:
//   - child: a pointer to the AST node containing the boolean value.
//   - parent: a pointer to the InstructionTree where the value will be written.
func WriteBoolLiteral(child *ast.AST, parent *tree.InstructionTree) {
	// Write 1 if the value is true, 0 otherwise
	if *child.Value == "true" {
		parent.Representation.WriteString("__TRUE")
	} else {
		parent.Representation.WriteString("__FALSE")
	}
	parent.Representation.WriteString(" ")
}
