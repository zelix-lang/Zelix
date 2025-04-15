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

package reassignment

import (
	"fluent/ast"
	"fluent/filecode"
	"fluent/filecode/function"
	module2 "fluent/filecode/module"
	"fluent/filecode/types/wrapper"
	"fluent/ir/pool"
	expression2 "fluent/ir/rule/expression"
	"fluent/ir/tree"
	"fluent/ir/value"
	"fluent/ir/variable"
	"fluent/util"
	"fmt"
	"strings"
)

func marshalExpr(
	queueElement *tree.BlockMarshalElement,
	trace *filecode.FileCode,
	fileCodeId int,
	traceFileName string,
	isMod bool,
	modulePropCounters *map[*module2.Module]*util.OrderedMap[string, *string],
	traceFn *function.Function,
	originalPath *string,
	counter *int,
	expr *ast.AST,
	variables *map[string]*variable.IRVariable,
	traceCounters *pool.NumPool,
	usedStrings *pool.StringPool,
	usedArrays *pool.StringPool,
	usedNumbers *pool.StringPool,
	localCounters *map[string]*string,
	isParam bool,
	expectedType *wrapper.TypeWrapper,
) string {
	// See if we can save memory if the expression is a static value
	dummyBuilder := strings.Builder{}
	if value.RetrieveStaticVal(fileCodeId, expr, &dummyBuilder, usedStrings, usedNumbers) {
		return dummyBuilder.String()
	}

	suitable := *counter

	// Marshal the expression directly
	expression2.MarshalExpression(
		queueElement.Representation,
		trace,
		traceFn,
		fileCodeId,
		isMod,
		traceFileName,
		originalPath,
		modulePropCounters,
		counter,
		expr,
		variables,
		traceCounters,
		usedStrings,
		usedArrays,
		usedNumbers,
		localCounters,
		true,
		isParam,
		expectedType,
	)

	return fmt.Sprintf("x%d ", suitable)
}

// MarshalReassignment marshals a reassignment operation in the intermediate representation (IR).
// It handles both variable reassignments and property reassignments.
//
// Parameters:
// - queueElement: The element in the block queue to be marshaled.
// - trace: The file code trace.
// - fileCodeId: The ID of the file code.
// - traceFileName: The name of the trace file.
// - isMod: A boolean indicating if the operation is a module.
// - modulePropCounters: A map of module property counters.
// - traceFn: The function trace.
// - originalPath: The original path of the file.
// - counter: A pointer to the counter for generating unique addresses.
// - element: The AST element representing the reassignment.
// - variables: A map of IR variables.
// - traceCounters: A pool of numeric counters.
// - usedStrings: A pool of used strings.
// - usedArrays: A pool of used arrays.
// - usedNumbers: A pool of used numbers.
// - localCounters: A map of local counters.
func MarshalReassignment(
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
	usedStrings *pool.StringPool,
	usedArrays *pool.StringPool,
	usedNumbers *pool.StringPool,
	localCounters *map[string]*string,
) {
	// Get the children
	children := *element.Children

	// Get the left expression
	leftExpr := children[0]

	// Get the right expression
	rightExpr := children[1]

	// Get the node from the left expression in order
	// to determine what to do
	left := (*leftExpr.Children)[0]

	// Check for variable reassignments
	if left.Rule == ast.Identifier {
		// Get the variable's name
		name := *left.Value

		// Get the variable
		storedVar := (*variables)[name]

		// Marshal the right expression
		right := marshalExpr(
			queueElement,
			trace,
			fileCodeId,
			traceFileName,
			isMod,
			modulePropCounters,
			traceFn,
			originalPath,
			counter,
			rightExpr,
			variables,
			traceCounters,
			usedStrings,
			usedArrays,
			usedNumbers,
			localCounters,
			true,
			storedVar.Type,
		)

		// Write store instructions
		queueElement.Representation.WriteString("store ")
		queueElement.Representation.WriteString(storedVar.Addr)
		queueElement.Representation.WriteString(" ")
		queueElement.Representation.WriteString(right)
		queueElement.Representation.WriteString("\n")

		return
	}

	// Marshal the right expression
	rightAddr := marshalExpr(
		queueElement,
		trace,
		fileCodeId,
		traceFileName,
		isMod,
		modulePropCounters,
		traceFn,
		originalPath,
		counter,
		rightExpr,
		variables,
		traceCounters,
		usedStrings,
		usedArrays,
		usedNumbers,
		localCounters,
		true,
		rightExpr.InferredType,
	)

	// Increment pointer count
	rightExpr.InferredType.PointerCount += 2

	// Store the new property in the pointer
	leftAddr := marshalExpr(
		queueElement,
		trace,
		fileCodeId,
		traceFileName,
		isMod,
		modulePropCounters,
		traceFn,
		originalPath,
		counter,
		leftExpr,
		variables,
		traceCounters,
		usedStrings,
		usedArrays,
		usedNumbers,
		localCounters,
		true,
		rightExpr.InferredType,
	)

	// Write store instructions
	queueElement.Representation.WriteString("store ")
	queueElement.Representation.WriteString(leftAddr)
	queueElement.Representation.WriteString(rightAddr)
	queueElement.Representation.WriteString("\n")
}
