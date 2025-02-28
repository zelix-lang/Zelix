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
	"fluent/ir/tree"
)

func RetrieveVarOrStr(
	fileCodeId int,
	expr *ast.AST,
	parent *tree.InstructionTree,
	usedStrings *pool.StringPool,
	usedNumbers *pool.StringPool,
	variables map[string]string,
) bool {
	// Get the expression's children
	exprChildren := *expr.Children

	// Check if we can reuse a string
	if len(exprChildren) == 1 {
		switch exprChildren[0].Rule {
		case ast.StringLiteral:
			strLiteral := exprChildren[0]
			parent.Representation.WriteString(
				usedStrings.RequestAddress(
					fileCodeId,
					*strLiteral.Value,
				),
			)

			parent.Representation.WriteString(" ")
			return true
		case ast.Identifier:
			// Write the variable's address
			parent.Representation.WriteString(variables[*expr.Value])
			parent.Representation.WriteString(" ")
			return true
		case ast.BooleanLiteral:
			// Write the boolean's value
			WriteBoolLiteral(exprChildren[0], parent)
			return true
		case ast.NumberLiteral, ast.DecimalLiteral:
			// Get the number's value
			num := *exprChildren[0].Value

			// See if the number's value is either 0 or 1
			if num == "0" {
				// Write the __FALSE constant
				parent.Representation.WriteString("__FALSE")
				parent.Representation.WriteString(" ")
				return true
			} else if num == "1" {
				// Write the __TRUE constant
				parent.Representation.WriteString("__TRUE")
				parent.Representation.WriteString(" ")
				return true
			}

			// Request an address for this number
			address := usedNumbers.RequestAddress(fileCodeId, num)
			parent.Representation.WriteString(address)
			parent.Representation.WriteString(" ")

			return true
		default:
		}
	}

	return false
}
