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
	"strings"
)

// WriteBoolLiteral writes a boolean literal to the parent InstructionTree.
// It writes "1 " if the value is true, "0 " otherwise.
//
// Parameters:
//   - child: a pointer to the AST node containing the boolean value.
//   - representation: a pointer to the strings.Builder where the value will be written.
func WriteBoolLiteral(child *ast.AST, representation *strings.Builder) {
	// Write 1 if the value is true, 0 otherwise
	if *child.Value == "true" {
		representation.WriteString("1 ")
	} else {
		representation.WriteString("0 ")
	}
}
