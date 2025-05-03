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
	"fluent/ir/pool"
	"fluent/ir/relocate"
	"fluent/ir/rule/conditional"
	"fluent/ir/rule/expression"
	"fluent/ir/tree"
	"fluent/ir/value"
	"fluent/ir/variable"
	"fluent/util"
	"fmt"
	"strings"
)

// MarshalWhile marshals a while loop into the intermediate representation.
// It takes various parameters to handle the marshalling process, including
// the queue element, trace information, file code ID, module properties,
// and pools for managing variables, strings, arrays, and numbers.
//
// Parameters:
// - queueElement: The current block marshal element in the queue.
// - trace: The file code trace information.
// - fileCodeId: The ID of the file code.
// - traceFileName: The name of the trace file.
// - isMod: A boolean indicating if the module is modified.
// - modulePropCounters: A map of module property counters.
// - traceFn: The function trace information.
// - originalPath: The original path of the file.
// - counter: A pointer to the counter for generating unique IDs.
// - element: The AST element representing the while loop.
// - variables: A map of IR variables.
// - traceCounters: A pool for managing trace counters.
// - appendedBlocks: A pool for managing appended blocks.
// - usedStrings: A pool for managing used strings.
// - usedArrays: A pool for managing used arrays.
// - usedNumbers: A pool for managing used numbers.
// - localCounters: A map of local counters.
// - blockQueue: A queue of block marshal elements.
func MarshalWhile(
	queueElement *tree.BlockMarshalElement,
	trace *filecode.FileCode,
	fileCodeId int,
	traceFileName string,
	isMod bool,
	modulePropCounters *map[*module2.Module]*util.OrderedMap[string, *string],
	traceFn *function.Function,
	originalPath *string,
	counter *int,
	element *ast.AST,
	variables *map[string]*variable.IRVariable,
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

	// Get the conditional
	condition := children[0]

	// Get the block
	block := children[1]

	// Request an address for the conditional block
	conditionalAddr, conditionalBuilder := appendedBlocks.RequestAddress()

	// Request an address for the block
	blockAddr, blockBuilder := appendedBlocks.RequestAddress()

	// Relocate the rest of the code
	remainingAddr := relocate.Remaining(appendedBlocks, blockQueue, queueElement)

	// Write the appropriate instructions
	queueElement.Representation.WriteString("jump ")
	queueElement.Representation.WriteString(*conditionalAddr)
	queueElement.Representation.WriteString("\n")

	// Create a temporary builder to marshal the condition
	tempBuilder := strings.Builder{}
	var conditionAddr string

	// See if we can save memory on the condition
	if value.RetrieveStaticVal(fileCodeId, condition, &tempBuilder, usedStrings) {
		conditionAddr = tempBuilder.String()
	} else {
		conditionAddr = fmt.Sprintf("x%d ", *counter)

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
			condition,
			variables,
			traceCounters,
			usedStrings,
			usedArrays,
			usedNumbers,
			localCounters,
			true,
			true,
			&conditional.BooleanTypeWrapper,
		)

		conditionalBuilder.WriteString(tempBuilder.String())
	}

	// Write the conditional
	conditionalBuilder.WriteString("if ")
	conditionalBuilder.WriteString(conditionAddr)
	conditionalBuilder.WriteString(*blockAddr)
	conditionalBuilder.WriteString(" ")
	conditionalBuilder.WriteString(*remainingAddr)
	conditionalBuilder.WriteString("\n")

	// Schedule the block
	*blockQueue = append(*blockQueue, &tree.BlockMarshalElement{
		Element:        block,
		Representation: blockBuilder,
		ParentAddr:     conditionalAddr,
		JumpToParent:   true,
		RemainingAddr:  remainingAddr,
		Id:             appendedBlocks.Counter,
	})

}
