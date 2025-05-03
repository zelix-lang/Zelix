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
	"fluent/ir/pool"
	"strings"
)

// RetrieveStaticVal processes a static value from the AST and writes its representation.
// It returns true if the value was successfully processed, otherwise false.
//
// Parameters:
// - fileCodeId: An integer representing the file code ID.
// - expr: A pointer to the AST node representing the expression.
// - representation: A pointer to a strings.Builder to write the representation.
// - usedStrings: A pointer to a pool.StringPool for managing used strings.
func RetrieveStaticVal(
	fileCodeId int,
	expr *ast.AST,
	representation *strings.Builder,
	usedStrings *pool.StringPool,
) bool {
	// Get the expression's children
	exprChildren := *expr.Children
	child := exprChildren[0]

	// Check if we can reuse a string
	if len(exprChildren) == 1 {
		switch child.Rule {
		case ast.StringLiteral:
			representation.WriteString(
				usedStrings.RequestAddress(
					fileCodeId,
					*child.Value,
				),
			)

			representation.WriteString(" ")
			return true
		case ast.BooleanLiteral:
			// Write the boolean's value
			WriteBoolLiteral(child, representation)
			return true
		default:
		}
	}

	return false
}
