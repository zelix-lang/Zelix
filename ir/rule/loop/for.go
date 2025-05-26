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
	"fluent/ir/variable"
	"fluent/util"
	"fmt"
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

	// Create a new block for the condition break
	breakConditionAddr, breakConditionBuilder := appendedBlocks.RequestAddress()

	// Create a new block for the loop body
	blockAddr, blockBuilder := appendedBlocks.RequestAddress()

	// Create a new block for updating the loop counter
	storeAddr, storeAddrBuilder := appendedBlocks.RequestAddress()

	// Marshal the left and right expressions directly
	identifierAddr := fmt.Sprintf("x%d", *counter)
	(*variables)[*identifier.Value] = &variable.IRVariable{
		Addr: identifierAddr,
		Type: &numWrapper,
	}

	// Use the left expression as our loop variable
	expression.MarshalExpression(
		queueElement.Representation,
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
		localCounters,
		true,
		false,
		&numWrapper,
	)

	rightExprAddr := fmt.Sprintf("x%d", *counter)
	expression.MarshalExpression(
		queueElement.Representation,
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
		localCounters,
		true,
		true,
		&numWrapper,
	)

	// Create a backup variable to load the identifier
	backupAddr := fmt.Sprintf("x%d", *counter)
	breakConditionBuilder.WriteString("mov ")
	breakConditionBuilder.WriteString(backupAddr)
	breakConditionBuilder.WriteString(" take ")
	breakConditionBuilder.WriteString(identifierAddr)
	breakConditionBuilder.WriteString("\n")
	*counter++

	checkConditionAddr := fmt.Sprintf("x%d", *counter)
	breakConditionBuilder.WriteString("mov ")
	breakConditionBuilder.WriteString(checkConditionAddr)
	breakConditionBuilder.WriteString(" bool eq ")
	breakConditionBuilder.WriteString(backupAddr)
	breakConditionBuilder.WriteString(" ")
	breakConditionBuilder.WriteString(rightExprAddr)
	breakConditionBuilder.WriteString("\n")
	*counter++

	// Write the break block
	breakConditionBuilder.WriteString("if ")
	breakConditionBuilder.WriteString(checkConditionAddr)
	breakConditionBuilder.WriteString(" ")
	breakConditionBuilder.WriteString(*remainingAddr)
	breakConditionBuilder.WriteString(" ")
	breakConditionBuilder.WriteString(*blockAddr)
	breakConditionBuilder.WriteString("\n")

	// Create another backup variable to load the identifier
	backupAddr = fmt.Sprintf("x%d", *counter)
	storeAddrBuilder.WriteString("mov ")
	storeAddrBuilder.WriteString(backupAddr)
	storeAddrBuilder.WriteString(" take ")
	storeAddrBuilder.WriteString(identifierAddr)
	storeAddrBuilder.WriteString("\n")
	*counter++

	// Write the store block
	storeAddrBuilder.WriteString("store ")
	storeAddrBuilder.WriteString(identifierAddr)
	storeAddrBuilder.WriteString(" add ")
	storeAddrBuilder.WriteString(backupAddr)
	storeAddrBuilder.WriteString(" __fluentc_const_one\n")
	storeAddrBuilder.WriteString("jump ")
	storeAddrBuilder.WriteString(*breakConditionAddr)
	storeAddrBuilder.WriteString("\n")

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
