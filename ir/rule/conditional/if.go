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

package conditional

import (
	"fluent/ast"
	"fluent/filecode"
	"fluent/filecode/types/wrapper"
	"fluent/ir/pool"
	"fluent/ir/rule/expression"
	"fluent/ir/tree"
	"fluent/ir/value"
	"fluent/util"
	"strconv"
	"strings"
)

// A dummy boolean type wrapper
var booleanTypeWrapper = wrapper.TypeWrapper{
	BaseType:    "bool",
	IsPrimitive: true,
	Children:    &[]*wrapper.TypeWrapper{},
}

func marshalCondition(
	representation *strings.Builder,
	trace *filecode.FileCode,
	fileCodeId int,
	traceFileName string,
	modulePropCounters *map[string]*util.OrderedMap[string, *string],
	counter *int,
	element *ast.AST,
	children []*ast.AST,
	variables map[string]string,
	traceCounters *pool.NumPool,
	appendedBlocks *pool.BlockPool,
	usedStrings *pool.StringPool,
	usedArrays *pool.StringPool,
	usedNumbers *pool.StringPool,
	nameCounters *map[string]map[string]string,
	localCounters *map[string]string,
	blockQueue *[]tree.BlockMarshalElement,
	isLast bool,
) *strings.Builder {
	var block *ast.AST
	var condition *ast.AST
	var nextBuilder *strings.Builder
	var blockAddress *string
	var blockBuilder *strings.Builder

	if element.Rule == ast.Else {
		block = children[0]
	} else {
		condition = children[0]
		block = children[1]
	}

	// Marshal the condition if needed
	if condition != nil {
		// Request an address for the block
		blockAddress, blockBuilder = appendedBlocks.RequestAddress()

		// Create a temporary builder
		tempBuilder := strings.Builder{}

		// See if we can save memory on the condition
		if value.RetrieveStaticVal(fileCodeId, condition, &tempBuilder, usedStrings, usedNumbers, variables) {
			// Write the instruction
			representation.WriteString("if ")
			representation.WriteString(tempBuilder.String())
		} else {
			suitable := *counter

			// Marshal the expression
			expression.MarshalExpression(
				representation,
				trace,
				fileCodeId,
				traceFileName,
				modulePropCounters,
				counter,
				condition,
				variables,
				traceCounters,
				usedStrings,
				usedArrays,
				usedNumbers,
				nameCounters,
				localCounters,
				true,
				&booleanTypeWrapper,
			)

			// Write the instructions
			representation.WriteString("if ")
			representation.WriteString("x")
			representation.WriteString(strconv.Itoa(suitable))
			representation.WriteString(" ")
		}

		// Write the block's address
		representation.WriteString(*blockAddress)
		representation.WriteString(" ")

		if !isLast {
			// Request a block for the next element
			nextAddr, nextRepresentation := appendedBlocks.RequestAddress()
			representation.WriteString(*nextAddr)
			representation.WriteString("\n")
			nextBuilder = nextRepresentation
		} else {
			// Write an end block
			representation.WriteString("__block_end__\n")
		}
	} else {
		// Don't request any other block and reuse the existing one
		blockBuilder = representation
	}

	// Schedule the block
	*blockQueue = append(*blockQueue, tree.BlockMarshalElement{
		Element:        block,
		Representation: blockBuilder,
	})

	return nextBuilder
}

func MarshalIf(
	representation *strings.Builder,
	trace *filecode.FileCode,
	fileCodeId int,
	traceFileName string,
	modulePropCounters *map[string]*util.OrderedMap[string, *string],
	counter *int,
	element *ast.AST,
	variables map[string]string,
	traceCounters *pool.NumPool,
	appendedBlocks *pool.BlockPool,
	usedStrings *pool.StringPool,
	usedArrays *pool.StringPool,
	usedNumbers *pool.StringPool,
	nameCounters *map[string]map[string]string,
	localCounters *map[string]string,
	blockQueue *[]tree.BlockMarshalElement,
) {
	// Get the expression's children
	children := *element.Children

	// Determine if this expression has an else/elseif block
	childrenLen := len(children) - 1
	lastRepresentation := representation

	// Marshal all other conditions
	for i := 0; i <= childrenLen; i++ {
		// Determine if this child is the last one
		isLast := i == childrenLen

		var child *ast.AST
		if i == 0 {
			child = element
		} else {
			child = children[i]
		}

		// Skip the first condition's block
		if i == 1 {
			continue
		}

		// Check if we have any elements remaining
		if i == 0 {
			isLast = childrenLen == 1
		}

		exprChildren := *child.Children

		// Marshal this condition
		newRepresentation := marshalCondition(
			lastRepresentation,
			trace,
			fileCodeId,
			traceFileName,
			modulePropCounters,
			counter,
			child,
			exprChildren,
			variables,
			traceCounters,
			appendedBlocks,
			usedStrings,
			usedArrays,
			usedNumbers,
			nameCounters,
			localCounters,
			blockQueue,
			isLast,
		)

		// Update the last representation
		lastRepresentation = newRepresentation
	}
}
