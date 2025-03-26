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
	"fluent/ir/variable"
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
// - usedNumbers: A pointer to a pool.StringPool for managing used numbers.
// - variables: A pointer to a map of variable names to IRVariable pointers.
func RetrieveStaticVal(
	fileCodeId int,
	expr *ast.AST,
	representation *strings.Builder,
	usedStrings *pool.StringPool,
	usedNumbers *pool.StringPool,
	variables *map[string]*variable.IRVariable,
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
		case ast.Identifier:
			// Write the variable's address
			representation.WriteString((*variables)[*child.Value].Addr)
			representation.WriteString(" ")
			return true
		case ast.BooleanLiteral:
			// Write the boolean's value
			WriteBoolLiteral(child, representation)
			return true
		case ast.NumberLiteral, ast.DecimalLiteral:
			// Get the number's value
			num := *child.Value

			// See if the number's value is either 0 or 1
			if num == "0" {
				// Write the __FALSE constant
				representation.WriteString("__FALSE")
				representation.WriteString(" ")
				return true
			} else if num == "1" {
				// Write the __TRUE constant
				representation.WriteString("__TRUE")
				representation.WriteString(" ")
				return true
			}

			// Request an address for this number
			address := usedNumbers.RequestAddress(fileCodeId, num)
			representation.WriteString(address)
			representation.WriteString(" ")

			return true
		default:
		}
	}

	return false
}
