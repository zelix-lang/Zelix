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

package array

import (
	"fluent/ast"
	"fluent/filecode/types"
	"fluent/ir/pool"
	"fluent/ir/rule/value"
	"fluent/ir/tree"
	"strconv"
	"strings"
)

func MarshalArray(
	global *tree.InstructionTree,
	child *ast.AST,
	fileCodeId int,
	counter *int,
	parent *tree.InstructionTree,
	usedStrings *pool.StringPool,
	exprQueue *[]tree.MarshalPair,
	variables map[string]string,
	expected types.TypeWrapper,
) {
	// Get the array's children
	children := *child.Children

	for _, expr := range children {
		// Check for string literals
		if value.RetrieveVarOrStr(fileCodeId, expr, parent, usedStrings, variables) {
			continue
		}

		// Get a suitable counter
		suitable := *counter
		*counter++
		parent.Representation.WriteString("x")
		parent.Representation.WriteString(strconv.Itoa(suitable))
		parent.Representation.WriteString(" ")

		// Create a new InstructionTree
		instructionTree := tree.InstructionTree{
			Children:       &[]*tree.InstructionTree{},
			Representation: &strings.Builder{},
		}

		*global.Children = append([]*tree.InstructionTree{&instructionTree}, *global.Children...)

		// Add the expression to the expression queue
		*exprQueue = append(*exprQueue, tree.MarshalPair{
			Child:   expr,
			Parent:  &instructionTree,
			Counter: suitable,
			Expected: types.TypeWrapper{
				PointerCount: expected.PointerCount,
				ArrayCount:   expected.ArrayCount - 1,
				Children:     expected.Children,
				BaseType:     expected.BaseType,
				IsPrimitive:  expected.IsPrimitive,
			},
			IsParam: true,
		})
	}
}
