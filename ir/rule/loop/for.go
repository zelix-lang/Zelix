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

package loop

import (
	"fluent/ast"
	"fluent/filecode"
	"fluent/filecode/types/wrapper"
	"fluent/ir/pool"
	"fluent/ir/relocate"
	"fluent/ir/rule/expression"
	"fluent/ir/tree"
	"fluent/ir/value"
	"fluent/util"
	"fmt"
	"strconv"
	"strings"
)

var numWrapper = wrapper.TypeWrapper{
	BaseType:    "num",
	IsPrimitive: true,
	Children:    &[]*wrapper.TypeWrapper{},
}

func MarshalFor(
	queueElement *tree.BlockMarshalElement,
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
	localCounters *map[string]*string,
	blockQueue *[]*tree.BlockMarshalElement,
) {
	// Get the children
	children := *element.Children

	// Get the left and right expressions
	leftExpr := children[0]
	rightExpr := children[1]

	// Get the identifier
	identifier := children[2]

	// Get the block
	block := children[3]

	// Relocate the rest of the code
	remainingAddr := relocate.Remaining(appendedBlocks, blockQueue, queueElement)

	// Get a suitable counter for the identifier
	suitable := *counter
	identifierAddr := fmt.Sprintf("x%d", suitable)
	variables[*identifier.Value] = identifierAddr

	tempBuilder := strings.Builder{}
	// See if we can save memory on the left value
	if value.RetrieveStaticVal(fileCodeId, leftExpr, &tempBuilder, usedStrings, usedNumbers, variables) {
		// Move the left value to the stack
		queueElement.Representation.WriteString("mov ")
		queueElement.Representation.WriteString(identifierAddr)
		queueElement.Representation.WriteString(" num ")
		queueElement.Representation.WriteString(tempBuilder.String())
		queueElement.Representation.WriteString("\n")
		*counter++
	} else {
		// Marshal the expression directly
		expression.MarshalExpression(
			&tempBuilder,
			trace,
			fileCodeId,
			traceFileName,
			modulePropCounters,
			counter,
			leftExpr,
			variables,
			traceCounters,
			usedStrings,
			usedArrays,
			usedNumbers,
			localCounters,
			true,
			&numWrapper,
		)

		queueElement.Representation.WriteString(tempBuilder.String())
	}

	// See if we can save memory on the right value
	tempBuilder.Reset()
	var rightAddr string
	if value.RetrieveStaticVal(fileCodeId, rightExpr, &tempBuilder, usedStrings, usedNumbers, variables) {
		rightAddr = tempBuilder.String()
	} else {
		rightAddr = fmt.Sprintf("x%d", *counter)
		// Marshal the expression directly
		expression.MarshalExpression(
			&tempBuilder,
			trace,
			fileCodeId,
			traceFileName,
			modulePropCounters,
			counter,
			rightExpr,
			variables,
			traceCounters,
			usedStrings,
			usedArrays,
			usedNumbers,
			localCounters,
			true,
			&numWrapper,
		)

		queueElement.Representation.WriteString(tempBuilder.String())
	}

	// Get an address for the break conditional branch
	breakConditionAddr, breakConditionBuilder := appendedBlocks.RequestAddress()

	// Get an address for the conditional branch
	conditionAddr, conditionBuilder := appendedBlocks.RequestAddress()

	// Get an address for the loop's block
	blockAddr, blockBuilder := appendedBlocks.RequestAddress()

	// Get an address for the block that reassigns the variable
	storeBlockAddr, storeBlockBuilder := appendedBlocks.RequestAddress()

	// Write the store block
	storeBlockBuilder.WriteString("store ")
	storeBlockBuilder.WriteString(identifierAddr)
	storeBlockBuilder.WriteString(" add ")
	storeBlockBuilder.WriteString(identifierAddr)
	storeBlockBuilder.WriteString(" ")
	storeBlockBuilder.WriteString(usedNumbers.RequestAddress(fileCodeId, "1"))
	storeBlockBuilder.WriteString("\njump ")
	storeBlockBuilder.WriteString(*breakConditionAddr)
	storeBlockBuilder.WriteString("\n")

	// Get a suitable counter for break condition
	suitable = *counter
	*counter++
	breakConditionBuilder.WriteString("mov x")
	breakConditionBuilder.WriteString(strconv.Itoa(suitable))
	breakConditionBuilder.WriteString(" bool eq ")
	breakConditionBuilder.WriteString(identifierAddr)
	breakConditionBuilder.WriteString(" ")
	breakConditionBuilder.WriteString(rightAddr)
	breakConditionBuilder.WriteString("\n")

	// Write the break condition
	breakConditionBuilder.WriteString("if x")
	breakConditionBuilder.WriteString(strconv.Itoa(suitable))
	breakConditionBuilder.WriteString(" ")
	breakConditionBuilder.WriteString(*remainingAddr)
	breakConditionBuilder.WriteString(" ")
	breakConditionBuilder.WriteString(*conditionAddr)
	breakConditionBuilder.WriteString("\n")

	// Get a suitable counter for the condition
	suitable = *counter
	*counter++
	conditionBuilder.WriteString("mov x")
	conditionBuilder.WriteString(strconv.Itoa(suitable))
	conditionBuilder.WriteString(" bool lt ")
	conditionBuilder.WriteString(identifierAddr)
	conditionBuilder.WriteString(" ")
	conditionBuilder.WriteString(rightAddr)
	conditionBuilder.WriteString("\n")

	// Write the condition
	conditionBuilder.WriteString("if x")
	conditionBuilder.WriteString(strconv.Itoa(suitable))
	conditionBuilder.WriteString(" ")
	conditionBuilder.WriteString(*blockAddr)
	conditionBuilder.WriteString(" ")
	conditionBuilder.WriteString(*storeBlockAddr)
	conditionBuilder.WriteString("\n")

	// Schedule the block for marshaling
	*blockQueue = append(*blockQueue, &tree.BlockMarshalElement{
		Element:        block,
		Representation: blockBuilder,
		ParentAddr:     breakConditionAddr,
		JumpToParent:   true,
		RemainingAddr:  remainingAddr,
		Id:             appendedBlocks.Counter,
	})

	// Write the appropriate instructions
	queueElement.Representation.WriteString("jump ")
	queueElement.Representation.WriteString(*breakConditionAddr)
	queueElement.Representation.WriteString("\n")
}
