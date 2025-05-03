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

package signed

import (
	"fluent/ast"
	"fluent/filecode/types/wrapper"
	"fluent/ir/pool"
	"fluent/ir/tree"
	"fluent/ir/value"
	"strconv"
	"strings"
)

// Create global TypeWrappers for the signed expressions
var booleanWrapper = wrapper.TypeWrapper{
	BaseType:    "bool",
	IsPrimitive: true,
	Children:    &[]*wrapper.TypeWrapper{},
}

// writeSignOpcode writes the appropriate opcode to the parent InstructionTree
// based on the provided sign. It returns true if the sign represents a boolean
// operation, otherwise false.
//
// Parameters:
// - sign: The sign of the operation (e.g., "+", "-", "==").
// - parent: The parent InstructionTree to which the opcode will be written.
//
// Returns:
// - bool: True if the sign represents a boolean operation, otherwise false.
func writeSignOpcode(sign string, parent *tree.InstructionTree) bool {
	switch sign {
	case "+":
		parent.Representation.WriteString("add ")
		return false
	case "-":
		parent.Representation.WriteString("sub ")
		return false
	case "*":
		parent.Representation.WriteString("mul ")
		return false
	case "/":
		parent.Representation.WriteString("div ")
		return false
	case "==":
		parent.Representation.WriteString("eq ")
		return true
	case ">":
		parent.Representation.WriteString("gt ")
		return true
	case "<":
		parent.Representation.WriteString("lt ")
		return true
	case "<=":
		parent.Representation.WriteString("le ")
		return true
	case ">=":
		parent.Representation.WriteString("ge ")
		return true
	case "!=":
		parent.Representation.WriteString("ne ")
		return true
	case "||":
		parent.Representation.WriteString("or ")
		return true
	case "&&":
		parent.Representation.WriteString("and ")
		return true
	}
	return false
}

// processCandidate processes a candidate AST node and schedules it for further processing.
//
// Parameters:
// - global: The global InstructionTree.
// - candidate: The candidate AST node to process.
// - fileCodeId: The ID of the file code.
// - isBool: A boolean indicating if the candidate is a boolean expression.
// - counter: A pointer to the counter used for generating unique identifiers.
// - pair: The MarshalPair containing the parent InstructionTree and expected type.
// - preferredParent: The preferred parent InstructionTree for the candidate.
// - usedStrings: The pool of used strings.
// - exprQueue: The queue of expressions to be processed.
// - variables: The map of variables used in the IR.
func processCandidate(
	global *tree.InstructionTree,
	candidate *ast.AST,
	fileCodeId int,
	isBool bool,
	counter *int,
	pair *tree.MarshalPair,
	preferredParent *tree.InstructionTree,
	usedStrings *pool.StringPool,
	exprQueue *[]tree.MarshalPair,
) {
	// See if we can save memory
	if value.RetrieveStaticVal(fileCodeId, candidate, preferredParent.Representation, usedStrings) {
		return
	}

	// Get a suitable counter
	suitable := *counter
	*counter++

	preferredParent.Representation.WriteString("x")
	preferredParent.Representation.WriteString(strconv.Itoa(suitable))
	preferredParent.Representation.WriteString(" ")

	// Create a new InstructionTree for the candidate
	candidateTree := tree.InstructionTree{
		Children:       &[]*tree.InstructionTree{},
		Representation: &strings.Builder{},
	}

	*global.Children = append([]*tree.InstructionTree{&candidateTree}, *global.Children...)

	// Schedule the candidate
	var expected wrapper.TypeWrapper
	if isBool {
		expected = booleanWrapper
	} else {
		expected = pair.Expected
	}

	*exprQueue = append(*exprQueue, tree.MarshalPair{
		Child:       candidate,
		Parent:      &candidateTree,
		IsInline:    pair.IsInline,
		Counter:     suitable,
		MoveToStack: true,
		Expected:    expected,
	})
}

// MarshalSignedExpression marshals a signed expression into the intermediate representation (IR).
//
// Parameters:
// - global: The global InstructionTree.
// - child: The AST node representing the signed expression.
// - fileCodeId: The ID of the file code.
// - counter: A pointer to the counter used for generating unique identifiers.
// - pair: The MarshalPair containing the parent InstructionTree and expected type.
// - usedStrings: The pool of used strings.
// - exprQueue: The queue of expressions to be processed.
func MarshalSignedExpression(
	global *tree.InstructionTree,
	child *ast.AST,
	fileCodeId int,
	counter *int,
	pair *tree.MarshalPair,
	usedStrings *pool.StringPool,
	exprQueue *[]tree.MarshalPair,
) {
	children := *child.Children
	var expr *ast.AST

	if len(children) == 1 {
		expr = children[0]

		if expr.Rule == ast.Expression {
			children = *expr.Children
			expr = children[0]
			children = *expr.Children
		}
	} else {
		expr = child
	}

	// Determine the expression's sign
	generalSign := *children[1].Value

	// Write the appropriate opcode depending on the sign
	isBoolean := writeSignOpcode(generalSign, pair.Parent)

	// Process the first pair outside the queue
	processCandidate(global, children[0], fileCodeId, isBoolean, counter, pair, pair.Parent, usedStrings, exprQueue)

	// See if we can save memory in the 2nd operand
	if len(children) == 3 {
		if value.RetrieveStaticVal(fileCodeId, children[2], pair.Parent.Representation, usedStrings) {
			return
		}
	}

	// Process the expression in a breadth-first manner
	queue := children[2:]

	// Get a suitable pointer for the expression
	suitable := *counter
	*counter++
	lastParent := pair.Parent

	for len(queue) > 0 {
		suitable = *counter
		*counter++

		// Create a new InstructionTree for this expression
		exprTree := tree.InstructionTree{
			Children:       &[]*tree.InstructionTree{},
			Representation: &strings.Builder{},
		}

		if len(queue) == 1 {
			// See if we can save memory in the operand
			if value.RetrieveStaticVal(fileCodeId, queue[0], lastParent.Representation, usedStrings) {
				break
			}

			// Write the new address to the last parent
			lastParent.Representation.WriteString("x")
			lastParent.Representation.WriteString(strconv.Itoa(suitable))
			lastParent.Representation.WriteString(" ")

			// Schedule the expression
			expr := queue[0]
			var expected wrapper.TypeWrapper
			if expr.Rule == ast.Expression && (*expr.Children)[0].Rule == ast.BooleanExpression {
				expected = booleanWrapper
			} else {
				expected = pair.Expected
			}

			*exprQueue = append(*exprQueue, tree.MarshalPair{
				Child:       expr,
				Parent:      &exprTree,
				IsInline:    pair.IsInline,
				Counter:     suitable,
				MoveToStack: true,
				Expected:    expected,
			})

			*global.Children = append([]*tree.InstructionTree{&exprTree}, *global.Children...)
			break
		}

		*global.Children = append([]*tree.InstructionTree{&exprTree}, *global.Children...)
		// Write the new address to the last parent
		lastParent.Representation.WriteString("x")
		lastParent.Representation.WriteString(strconv.Itoa(suitable))
		lastParent.Representation.WriteString(" ")

		// Write mov instructions to the tree
		exprTree.Representation.WriteString("mov x")
		exprTree.Representation.WriteString(strconv.Itoa(suitable))
		exprTree.Representation.WriteString(pair.Expected.Marshal())
		exprTree.Representation.WriteString(" ")

		// Get the candidate
		candidate := queue[0]

		// Determine the expression's sign
		sign := *queue[1].Value

		// Remove the 1st operand and the sign from the queue
		queue = queue[2:]

		// Write the appropriate opcode depending on the sign
		isBoolean := writeSignOpcode(sign, &exprTree)

		// Process the first pair outside the queue
		processCandidate(global, candidate, fileCodeId, isBoolean, counter, pair, &exprTree, usedStrings, exprQueue)

		// Update the last parent
		lastParent = &exprTree
	}
}
