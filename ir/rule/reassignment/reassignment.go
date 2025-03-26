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
	"fluent/filecode/module"
	"fluent/filecode/types/wrapper"
	"fluent/ir/pool"
	expression2 "fluent/ir/rule/expression"
	"fluent/ir/tree"
	"fluent/ir/value"
	"fluent/ir/variable"
	"fluent/util"
	"fmt"
	"strconv"
	"strings"
)

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
	modulePropCounters *map[string]*util.OrderedMap[string, *string],
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

		// Calculate the counter for this expression
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

		// Give this variable another address in order
		// to adhere to SSA form
		storedVar.Addr = fmt.Sprintf("x%d", suitable)

		return
	}

	// Otherwise, we have a property reassignment
	// Calculate the inferred type of the entire
	// expression
	var inferredType []*wrapper.TypeWrapper
	var lastModule *module.Module
	var candidateAddr string
	var candidateVarName *string

	// Get the left expression's children
	leftChildren := *left.Children

	// Get the candidate
	candidate := leftChildren[0]

	// Marshal the candidate if needed
	if (*candidate.Children)[0].Rule == ast.Identifier {
		candidateVarName = (*candidate.Children)[0].Value
		candidateAddr = (*variables)[*candidateVarName].Addr
	} else {
		candidateAddr = fmt.Sprintf("x%d", *counter)
		// Marshal the candidate directly
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
			rightExpr,
			variables,
			traceCounters,
			usedStrings,
			usedArrays,
			usedNumbers,
			localCounters,
			true,
			candidate.InferredType,
		)
	}

	inferredType = []*wrapper.TypeWrapper{candidate.InferredType}
	lastModule = trace.Modules[candidate.InferredType.BaseType]

	// Iterate over all elements of the AST to
	// calculate the type this expression evaluates to
	for i := 1; i < len(leftChildren); i++ {
		childExpr := leftChildren[i]
		childExprChildren := *childExpr.Children
		child := childExprChildren[0]

		// Update the inferred type accordingly
		switch child.Rule {
		case ast.Identifier:
			// Find the declaration
			declaration := lastModule.Declarations[*child.Value]
			inferredType = append(inferredType, &declaration.Type)

			// Update the last module if needed
			if !declaration.Type.IsPrimitive {
				lastModule = trace.Modules[declaration.Type.BaseType]
			}
		default:
		}
	}

	// Calculate the counter the right expression is going to have
	var rightExprAddr string

	// See if we can save memory on the right expression
	tempBuilder := strings.Builder{}
	if value.RetrieveStaticVal(fileCodeId, rightExpr, &tempBuilder, usedStrings, usedNumbers, variables) {
		str := tempBuilder.String()
		// Remove the last space
		rightExprAddr = str[:len(str)-1]
	} else {
		rightExprAddr = fmt.Sprintf("x%d ", *counter)
		// Marshal the right expression directly
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
			rightExpr,
			variables,
			traceCounters,
			usedStrings,
			usedArrays,
			usedNumbers,
			localCounters,
			true,
			inferredType[len(inferredType)-1],
		)
	}

	// Marshal all expression
	for i := len(leftChildren) - 1; i >= 1; i-- {
		childExpr := leftChildren[i]
		childType := inferredType[i]
		child := (*childExpr.Children)[0]

		// Find the mod
		prev := inferredType[i-1]
		// Get the declaration's index
		idx, _ := (*modulePropCounters)[prev.BaseType].Get(*child.Value)

		// Get a suitable counter to move this expression
		suitable := *counter
		*counter++

		// Move the element to the stack
		queueElement.Representation.WriteString("mov x")
		queueElement.Representation.WriteString(strconv.Itoa(suitable))
		queueElement.Representation.WriteString(" ")

		if childType.IsPrimitive {
			queueElement.Representation.WriteString(childType.Marshal())
		} else {
			queueElement.Representation.WriteString(*(*localCounters)[childType.BaseType])
		}
		queueElement.Representation.WriteString(" mod_copy ")
		queueElement.Representation.WriteString(candidateAddr)
		queueElement.Representation.WriteString(" ")
		queueElement.Representation.WriteString(*idx)
		queueElement.Representation.WriteString(" ")
		queueElement.Representation.WriteString(rightExprAddr)
		queueElement.Representation.WriteString("\n")

		rightExprAddr = fmt.Sprintf("x%d", suitable)
		candidateAddr = rightExprAddr
	}

	// Update the candidate in order to adhere to SSA if needed
	if candidateVarName != nil {
		(*variables)[*candidateVarName].Addr = candidateAddr
	}
}
