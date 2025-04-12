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

package expression

import (
	"fluent/ast"
	"fluent/filecode"
	"fluent/filecode/function"
	"fluent/filecode/types/wrapper"
	"fluent/ir/pool"
	"fluent/ir/rule/array"
	"fluent/ir/rule/call"
	"fluent/ir/rule/object"
	"fluent/ir/rule/property"
	"fluent/ir/rule/signed"
	"fluent/ir/tree"
	"fluent/ir/value"
	"fluent/ir/variable"
	"fluent/util"
	"strconv"
	"strings"
)

// MarshalExpression marshals an AST element into a string representation.
// It processes the AST element and its children, converting them into a
// string representation suitable for further processing or output.
//
// Parameters:
// - representation: A pointer to a strings.Builder to store the resulting string representation.
// - trace: A pointer to a filecode.FileCode for tracing information.
// - traceFn: A pointer to a function.Function for function-specific tracing information.
// - fileCodeId: An integer representing the file code ID.
// - isMod: A boolean indicating if the module is modified.
// - traceFileName: A string representing the trace file name.
// - originalPath: A pointer to a string representing the original path.
// - modulePropCounters: A pointer to a map of module property counters.
// - counter: A pointer to an integer counter for tracking elements.
// - element: A pointer to an ast.AST element to be marshaled.
// - variables: A pointer to a map of IR variables.
// - traceCounters: A pointer to a pool.NumPool for trace counters.
// - usedStrings: A pointer to a pool.StringPool for used strings.
// - usedArrays: A pointer to a pool.StringPool for used arrays.
// - usedNumbers: A pointer to a pool.StringPool for used numbers.
// - localCounters: A pointer to a map of local counters.
// - moveToStack: A boolean indicating if the value should be moved to the stack.
// - firstExpected: A pointer to a wrapper.TypeWrapper for the first expected type.
func MarshalExpression(
	representation *strings.Builder,
	trace *filecode.FileCode,
	traceFn *function.Function,
	fileCodeId int,
	isMod bool,
	traceFileName string,
	originalPath *string,
	modulePropCounters *map[string]*util.OrderedMap[string, *string],
	counter *int,
	element *ast.AST,
	variables *map[string]*variable.IRVariable,
	traceCounters *pool.NumPool,
	usedStrings *pool.StringPool,
	usedArrays *pool.StringPool,
	usedNumbers *pool.StringPool,
	localCounters *map[string]*string,
	moveToStack bool,
	firstExpected *wrapper.TypeWrapper,
) {
	result := tree.InstructionTree{
		Children:       &[]*tree.InstructionTree{},
		Representation: &strings.Builder{},
	}

	firstEl := tree.MarshalPair{
		Child:       element,
		Parent:      &result,
		IsInline:    false,
		MoveToStack: moveToStack,
	}

	// Get a suitable counter in case we have to move
	// this value to the stack
	if moveToStack {
		firstEl.Counter = *counter
		firstEl.Expected = *firstExpected
		*counter++
	}

	// Use a queue to process the elements of the AST
	queue := []tree.MarshalPair{
		firstEl,
	}

	for len(queue) > 0 {
		// Get the first element of the queue
		pair := queue[0]
		queue = queue[1:]
		derefVariables := true

		// Get the children of the current element
		children := *pair.Child.Children

		// Move values to the stack for parameters
		if (pair.MoveToStack || pair.IsParam) && !pair.IsInline {
			if !pair.IsParam {
				pair.Parent.Representation.WriteString("alloca x")
			} else {
				pair.Parent.Representation.WriteString("mov x")
			}

			pair.Parent.Representation.WriteString(strconv.Itoa(pair.Counter))

			if !pair.IsParam {
				pair.Parent.Representation.WriteString(" &")
			} else {
				pair.Parent.Representation.WriteString(" ")
			}

			if pair.Expected.IsPrimitive {
				pair.Parent.Representation.WriteString(pair.Expected.Marshal())
			} else {
				oldBaseType := pair.Expected.BaseType
				pair.Expected.BaseType = *(*localCounters)[oldBaseType]
				pair.Parent.Representation.WriteString(pair.Expected.Marshal())
				pair.Expected.BaseType = oldBaseType
			}

			pair.Parent.Representation.WriteString("\nstore x")
			pair.Parent.Representation.WriteString(strconv.Itoa(pair.Counter))
			pair.Parent.Representation.WriteString(" ")
		}

		// Get the remaining expression
		child := children[pair.StartAt]

		// Check for pointers
		if child.Rule == ast.Pointer {
			// Write addr instructions
			pair.Parent.Representation.WriteString("addr ")

			// See if we have a next pointer
			if children[pair.StartAt+1].Rule == ast.Pointer {
				suitable := *counter
				pair.Parent.Representation.WriteString("x")
				pair.Parent.Representation.WriteString(strconv.Itoa(suitable))
				pair.Parent.Representation.WriteString("\n")

				// Increment the counter
				*counter++

				newTree := tree.InstructionTree{
					Children:       &[]*tree.InstructionTree{},
					Representation: &strings.Builder{},
				}

				// Add new instructions to the current tree
				*pair.Parent.Children = append([]*tree.InstructionTree{&newTree}, *pair.Parent.Children...)

				// Add the new tree to the queue
				queue = append(queue, tree.MarshalPair{
					Child:       pair.Child,
					Parent:      pair.Parent,
					IsInline:    pair.IsInline,
					Counter:     suitable,
					MoveToStack: true,

					Expected: wrapper.TypeWrapper{
						PointerCount: pair.Expected.PointerCount - 1,
						ArrayCount:   pair.Expected.ArrayCount,
						Children:     pair.Expected.Children,
						BaseType:     pair.Expected.BaseType,
						IsPrimitive:  pair.Expected.IsPrimitive,
					},
					StartAt: pair.StartAt + 1,
				})

				continue
			}

			derefVariables = false
			child = children[pair.StartAt+1]
		}

		switch child.Rule {
		case ast.FunctionCall:
			call.MarshalFunctionCall(
				&result,
				child,
				traceFileName,
				fileCodeId,
				originalPath,
				isMod,
				trace,
				traceFn,
				counter,
				pair.Parent,
				traceCounters,
				usedStrings,
				usedNumbers,
				&queue,
				localCounters,
			)
		case ast.ObjectCreation:
			object.MarshalObjectCreation(
				&result,
				child,
				traceFileName,
				fileCodeId,
				originalPath,
				isMod,
				trace,
				traceFn,
				modulePropCounters,
				counter,
				&pair,
				traceCounters,
				usedStrings,
				usedArrays,
				usedNumbers,
				&queue,
				localCounters,
			)
		case ast.Identifier:
			// Get the variable
			stored := (*variables)[*child.Value]

			if derefVariables {
				pair.Parent.Representation.WriteString("load ")
			}

			pair.Parent.Representation.WriteString(stored.Addr)
		case ast.StringLiteral:
			// Request an address space for the string literal
			pair.Parent.Representation.WriteString(
				usedStrings.RequestAddress(
					fileCodeId,
					*child.Value,
				),
			)
		case ast.Array:
			array.MarshalArray(
				&result,
				child,
				fileCodeId,
				counter,
				&pair,
				usedStrings,
				usedNumbers,
				&queue,
			)
		case ast.NumberLiteral, ast.DecimalLiteral:
			// Directly write the tree's value
			pair.Parent.Representation.WriteString(
				usedNumbers.RequestAddress(
					fileCodeId,
					*child.Value,
				),
			)
		case ast.BooleanLiteral:
			value.WriteBoolLiteral(child, pair.Parent.Representation)
		case ast.Expression:
			// Add the expression to the queue
			queue = append(queue, tree.MarshalPair{
				Child:       child,
				Parent:      pair.Parent,
				Counter:     pair.Counter,
				Expected:    pair.Expected,
				MoveToStack: pair.MoveToStack,
				IsInline:    true,
			})
		case ast.ArithmeticExpression, ast.BooleanExpression:
			signed.MarshalSignedExpression(
				&result,
				child,
				fileCodeId,
				counter,
				&pair,
				usedStrings,
				usedNumbers,
				&queue,
			)
		case ast.PropertyAccess:
			property.MarshalPropertyAccess(
				&result,
				trace,
				child,
				fileCodeId,
				counter,
				&pair,
				modulePropCounters,
				traceCounters,
				usedStrings,
				usedNumbers,
				&queue,
				localCounters,
				traceFileName,
			)
		default:
		}
	}

	// Append all children to the parent tree
	for _, child := range *result.Children {
		representation.WriteString(child.Representation.String())
		representation.WriteString("\n")
	}

	// Append the expression itself (without the children)
	representation.WriteString(result.Representation.String())
	representation.WriteString("\n")
}
