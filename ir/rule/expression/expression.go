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
	"fluent/ir/pool"
	"fluent/ir/rule/array"
	"fluent/ir/rule/call"
	"fluent/ir/rule/object"
	"fluent/ir/rule/signed"
	"fluent/ir/tree"
	"fluent/ir/value"
	"fluent/util"
	"strconv"
	"strings"
)

func MarshalExpression(
	funTree *tree.InstructionTree,
	trace *filecode.FileCode,
	fileCodeId int,
	traceFileName string,
	modulePropCounters *map[string]*util.OrderedMap[string, *string],
	counter *int,
	element *ast.AST,
	variables map[string]string,
	traceCounters *pool.NumPool,
	usedStrings *pool.StringPool,
	usedArrays *pool.StringPool,
	usedNumbers *pool.StringPool,
	nameCounters *map[string]map[string]string,
	localCounters *map[string]string,
) {
	result := tree.InstructionTree{
		Children:       &[]*tree.InstructionTree{},
		Representation: &strings.Builder{},
	}

	// Use a queue to process the elements of the AST
	queue := []tree.MarshalPair{
		{
			Child:    element,
			Parent:   &result,
			IsInline: false,
		},
	}

	for len(queue) > 0 {
		// Get the first element of the queue
		pair := queue[0]
		queue = queue[1:]

		// Get the children of the current element
		children := *pair.Child.Children
		// Used to skip pointers and dereferences
		startAt := 0

		// Move values to the stack for parameters
		if pair.IsParam && !pair.IsInline {
			pair.Parent.Representation.WriteString("mov x")
			pair.Parent.Representation.WriteString(strconv.Itoa(pair.Counter))
			pair.Parent.Representation.WriteString(" ")

			if pair.Expected.IsPrimitive {
				pair.Parent.Representation.WriteString(pair.Expected.Marshal())
			} else {
				oldBaseType := pair.Expected.BaseType
				pair.Expected.BaseType = (*localCounters)[oldBaseType]
				pair.Parent.Representation.WriteString(pair.Expected.Marshal())
				pair.Expected.BaseType = oldBaseType
			}

			pair.Parent.Representation.WriteString(" ")
		}

		// Add pointers and dereferences
		for _, child := range children {
			if child.Rule == ast.Pointer {
				pair.Parent.Representation.WriteString("&")
				startAt++
			} else if child.Rule == ast.Dereference {
				pair.Parent.Representation.WriteString("*")
				startAt++
			} else {
				break
			}
		}

		// Get the remaining expression
		child := children[startAt]

		switch child.Rule {
		case ast.FunctionCall:
			call.MarshalFunctionCall(
				&result,
				child,
				traceFileName,
				fileCodeId,
				trace,
				counter,
				pair.Parent,
				traceCounters,
				nameCounters,
				variables,
				usedStrings,
				usedNumbers,
				&queue,
			)
		case ast.ObjectCreation:
			object.MarshalObjectCreation(
				&result,
				child,
				traceFileName,
				fileCodeId,
				trace,
				modulePropCounters,
				counter,
				&pair,
				traceCounters,
				variables,
				usedStrings,
				usedArrays,
				usedNumbers,
				&queue,
				localCounters,
			)
		case ast.Identifier:
			// Write the variable's name
			pair.Parent.Representation.WriteString(*child.Value)
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
				variables,
			)
		case ast.NumberLiteral, ast.DecimalLiteral:
			// Directly write the tree's value
			pair.Parent.Representation.WriteString(*child.Value)
		case ast.BooleanLiteral:
			value.WriteBoolLiteral(child, pair.Parent)
		case ast.Expression:
			// Add the expression to the queue
			queue = append(queue, tree.MarshalPair{
				Child:    child,
				Parent:   pair.Parent,
				Counter:  pair.Counter,
				Expected: pair.Expected,
				IsParam:  pair.IsParam,
				IsInline: true,
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
				variables,
			)
		default:
		}
	}

	// Append all children to the parent tree
	for _, child := range *result.Children {
		funTree.Representation.WriteString(child.Representation.String())
		funTree.Representation.WriteString("\n")
	}

	// Append the expression itself (without the children)
	funTree.Representation.WriteString(result.Representation.String())
	funTree.Representation.WriteString("\n")
}
