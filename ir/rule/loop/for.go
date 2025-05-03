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
	"fluent/filecode/function"
	module2 "fluent/filecode/module"
	"fluent/filecode/types/wrapper"
	"fluent/ir/pool"
	"fluent/ir/relocate"
	"fluent/ir/rule/expression"
	"fluent/ir/tree"
	"fluent/ir/value"
	"fluent/ir/variable"
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

// MarshalFor marshals a for-loop construct into the intermediate representation (IR).
// It processes the loop's initialization, condition, and increment expressions,
// and generates the corresponding IR code.
//
// Parameters:
// - queueElement: The current block element being marshaled.
// - trace: The file code trace information.
// - fileCodeId: The ID of the file code.
// - traceFileName: The name of the trace file.
// - isMod: A flag indicating if the module is modified.
// - modulePropCounters: A map of module property counters.
// - counter: A pointer to the current counter value.
// - traceFn: The function trace information.
// - originalPath: The original path of the file.
// - element: The AST element representing the for-loop.
// - variables: A map of IR variables.
// - traceCounters: A pool of numeric counters.
// - appendedBlocks: A pool of block elements.
// - usedStrings: A pool of used strings.
// - usedArrays: A pool of used arrays.
// - localCounters: A map of local counters.
// - blockQueue: A queue of block elements to be marshaled.
func MarshalFor(
	queueElement *tree.BlockMarshalElement,
	trace *filecode.FileCode,
	fileCodeId int,
	traceFileName string,
	isMod bool,
	modulePropCounters *map[*module2.Module]*util.OrderedMap[string, *string],
	counter *int,
	traceFn *function.Function,
	originalPath *string,
	element *ast.AST,
	variables *map[string]*variable.IRVariable,
	traceCounters *pool.NumPool,
	appendedBlocks *pool.BlockPool,
	usedStrings *pool.StringPool,
	usedArrays *pool.StringPool,
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

	// Get a backup variable to increment the loop counter
	suitable := *counter
	backupAddr := fmt.Sprintf("x%d", suitable)
	*counter++

	// Get a suitable counter for the identifier
	suitable = *counter
	identifierAddr := fmt.Sprintf("x%d", suitable)
	(*variables)[*identifier.Value] = &variable.IRVariable{
		Addr: identifierAddr,
		Type: &numWrapper,
	}

	leftTempBuilder := strings.Builder{}
	isBuilderStatic := true
	leftBranchAddr := ""
	// See if we can save memory on the left value
	if value.RetrieveStaticVal(fileCodeId, leftExpr, &leftTempBuilder, usedStrings) {
		*counter++
	} else {
		isBuilderStatic = false
		*counter++

		leftBranchAddr = fmt.Sprintf("x%d", *counter)
		// Marshal the expression directly
		expression.MarshalExpression(
			&leftTempBuilder,
			trace,
			traceFn,
			fileCodeId,
			isMod,
			traceFileName,
			originalPath,
			modulePropCounters,
			counter,
			leftExpr,
			variables,
			traceCounters,
			usedStrings,
			usedArrays,
			localCounters,
			true,
			false,
			&numWrapper,
		)

		queueElement.Representation.WriteString(leftTempBuilder.String())
	}

	tempBuilder := strings.Builder{}
	var rightAddr string
	if value.RetrieveStaticVal(fileCodeId, rightExpr, &tempBuilder, usedStrings) {
		rightAddr = tempBuilder.String()
	} else {
		rightAddr = fmt.Sprintf("x%d", *counter)
		// Marshal the expression directly
		expression.MarshalExpression(
			&tempBuilder,
			trace,
			traceFn,
			fileCodeId,
			isMod,
			traceFileName,
			originalPath,
			modulePropCounters,
			counter,
			rightExpr,
			variables,
			traceCounters,
			usedStrings,
			usedArrays,
			localCounters,
			true,
			false,
			&numWrapper,
		)

		queueElement.Representation.WriteString(tempBuilder.String())
	}

	// Get an address for the break conditional branch
	breakConditionAddr, breakConditionBuilder := appendedBlocks.RequestAddress()

	// Get an address for the loop's block
	blockAddr, blockBuilder := appendedBlocks.RequestAddress()

	// Get an address for the block that changes the value of the counter
	storeAddr, storeBuilder := appendedBlocks.RequestAddress()

	// Move the backup variable
	storeBuilder.WriteString("mov ")
	storeBuilder.WriteString(backupAddr)
	storeBuilder.WriteString(" num add ")
	storeBuilder.WriteString(identifierAddr)
	storeBuilder.WriteString(" 1 ")
	storeBuilder.WriteString("\njump ")
	storeBuilder.WriteString(*breakConditionAddr)
	storeBuilder.WriteString("\n")

	// Decide the value of the identifier
	breakConditionBuilder.WriteString("pick ")
	breakConditionBuilder.WriteString(identifierAddr)

	if queueElement.ParentAddr == nil {
		breakConditionBuilder.WriteString(" entry ")
	} else {
		breakConditionBuilder.WriteString(" ")
		breakConditionBuilder.WriteString(*queueElement.ParentAddr)
		breakConditionBuilder.WriteString(" ")
	}

	if isBuilderStatic {
		breakConditionBuilder.WriteString(leftTempBuilder.String())
	} else {
		breakConditionBuilder.WriteString(leftBranchAddr)
		breakConditionBuilder.WriteString(" ")
	}

	breakConditionBuilder.WriteString(*storeAddr)
	breakConditionBuilder.WriteString(" ")
	breakConditionBuilder.WriteString(backupAddr)
	breakConditionBuilder.WriteString("\n")

	// Get a suitable counter for the break condition
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
	breakConditionBuilder.WriteString(*blockAddr)
	breakConditionBuilder.WriteString("\n")

	// Get a suitable counter for the condition
	suitable = *counter
	*counter++

	// Schedule the block for marshaling
	*blockQueue = append(*blockQueue, &tree.BlockMarshalElement{
		Element:        block,
		Representation: blockBuilder,
		ParentAddr:     storeAddr,
		JumpToParent:   true,
		RemainingAddr:  remainingAddr,
		Id:             appendedBlocks.Counter,
	})

	// Write the appropriate instructions
	queueElement.Representation.WriteString("jump ")
	queueElement.Representation.WriteString(*breakConditionAddr)
	queueElement.Representation.WriteString("\n")
}
