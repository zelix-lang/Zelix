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
	"fluent/ir/variable"
	"fluent/util"
	"fmt"
	"strconv"
)

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
	var inferredType *wrapper.TypeWrapper
	var lastModule *module.Module

	// Get the left expression's children
	leftChildren := *left.Children

	// Get the candidate
	candidate := leftChildren[0]
	inferredType = candidate.InferredType
	lastModule = trace.Modules[inferredType.BaseType]

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
			inferredType = &declaration.Type

			// Update the last module if needed
			if !declaration.Type.IsPrimitive {
				lastModule = trace.Modules[declaration.Type.BaseType]
			}
		case ast.FunctionCall:
			// Get the expression's children
			callChildren := *child.Children

			// Get the function's name
			name := *callChildren[0].Value

			// Retrieve the function
			fun := lastModule.Functions[name]

			// Modify the inferred type accordingly
			inferredType = &fun.ReturnType

			// Update the last module if needed
			if !fun.ReturnType.IsPrimitive {
				lastModule = trace.Modules[fun.ReturnType.BaseType]
			}
		default:
		}
	}

	// Calculate the counter the right expression is going to have
	rightExprCounter := *counter

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
		inferredType,
	)

	// Calculate the counter the expression is going to have
	exprCounter := *counter

	// Marshal the left expression directly
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
		leftExpr,
		variables,
		traceCounters,
		usedStrings,
		usedArrays,
		usedNumbers,
		localCounters,
		true,
		inferredType,
	)

	// Add store instructions
	queueElement.Representation.WriteString("store x")
	queueElement.Representation.WriteString(strconv.Itoa(exprCounter))
	queueElement.Representation.WriteString(" x")
	queueElement.Representation.WriteString(strconv.Itoa(rightExprCounter))
	queueElement.Representation.WriteString("\n")
}
