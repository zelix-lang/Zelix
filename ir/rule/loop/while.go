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
	"fluent/ir/pool"
	"fluent/ir/rule/conditional"
	"fluent/ir/rule/expression"
	"fluent/ir/tree"
	"fluent/ir/value"
	"fluent/util"
	"fmt"
	"strings"
)

func MarshalWhile(
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
	blockQueue *[]*tree.BlockMarshalElement,
) {
	// Get the children
	children := *element.Children

	// Get the conditional
	condition := children[0]

	// Get the block
	block := children[1]

	// Request an address for the conditional block
	conditionalAddr, conditionalBuilder := appendedBlocks.RequestAddress()

	// Request an address for the block
	blockAddr, blockBuilder := appendedBlocks.RequestAddress()

	// Write the appropriate instructions
	representation.WriteString("jump ")
	representation.WriteString(*conditionalAddr)
	representation.WriteString("\n")

	// Create a temporary builder to marshal the condition
	tempBuilder := strings.Builder{}
	var conditionAddr string

	// See if we can save memory on the condition
	if value.RetrieveStaticVal(fileCodeId, condition, &tempBuilder, usedStrings, usedNumbers, variables) {
		conditionAddr = tempBuilder.String()
	} else {
		conditionAddr = fmt.Sprintf("x%d ", *counter)

		// Marshal the expression directly
		expression.MarshalExpression(
			&tempBuilder,
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
			&conditional.BooleanTypeWrapper,
		)

		conditionalBuilder.WriteString(tempBuilder.String())
	}

	// Write the conditional
	conditionalBuilder.WriteString("if ")
	conditionalBuilder.WriteString(conditionAddr)
	conditionalBuilder.WriteString(*blockAddr)
	conditionalBuilder.WriteString(" __block_end__\n")

	// Schedule the block
	*blockQueue = append(*blockQueue, &tree.BlockMarshalElement{
		Element:        block,
		Representation: blockBuilder,
		ParentAddr:     conditionalAddr,
		JumpToParent:   true,
	})

}
