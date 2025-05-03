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
	"fluent/filecode/types/wrapper"
	"fluent/ir/pool"
	"fluent/ir/tree"
	"fluent/ir/value"
	"strconv"
	"strings"
)

// MarshalArray serializes an AST array node into a string representation
// and updates the instruction tree and expression queue accordingly.
//
// Parameters:
// - global: The global instruction tree.
// - child: The AST node representing the array.
// - fileCodeId: The ID of the file code.
// - counter: A pointer to the counter used for generating unique identifiers.
// - pair: The marshal pair containing the parent instruction tree and expected type.
// - usedStrings: A pool of used string literals.
// - usedNumbers: A pool of used number literals.
// - exprQueue: A queue of marshal pairs for further processing.
func MarshalArray(
	global *tree.InstructionTree,
	child *ast.AST,
	fileCodeId int,
	counter *int,
	pair *tree.MarshalPair,
	usedStrings *pool.StringPool,
	usedNumbers *pool.StringPool,
	exprQueue *[]tree.MarshalPair,
) {
	// Get the array's children
	children := *child.Children

	// Write an arr opcode
	pair.Parent.Representation.WriteString("arr ")

	// Write the array's length
	pair.Parent.Representation.WriteString(strconv.Itoa(len(children)))
	pair.Parent.Representation.WriteString(" ")

	// Prevent collisions
	if (pair.MoveToStack || pair.IsParam) && pair.Counter == *counter {
		*counter++
	}

	for _, expr := range children {
		// Check for string literals
		if value.RetrieveStaticVal(fileCodeId, expr, pair.Parent.Representation, usedStrings) {
			continue
		}

		// Get a suitable counter
		suitable := *counter
		*counter++
		pair.Parent.Representation.WriteString("x")
		pair.Parent.Representation.WriteString(strconv.Itoa(suitable))
		pair.Parent.Representation.WriteString(" ")

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
			Expected: wrapper.TypeWrapper{
				PointerCount: pair.Expected.PointerCount,
				ArrayCount:   pair.Expected.ArrayCount - 1,
				Children:     pair.Expected.Children,
				BaseType:     pair.Expected.BaseType,
				IsPrimitive:  pair.Expected.IsPrimitive,
			},
			IsParam: true,
		})
	}
}
