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

package declaration

import (
	"fluent/ast"
	"fluent/filecode"
	"fluent/filecode/function"
	"fluent/filecode/types"
	"fluent/ir/pool"
	expression2 "fluent/ir/rule/expression"
	"fluent/ir/tree"
	"fluent/util"
	"fmt"
)

func MarshalDeclaration(
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
	variables *map[string]string,
	traceCounters *pool.NumPool,
	usedStrings *pool.StringPool,
	usedArrays *pool.StringPool,
	usedNumbers *pool.StringPool,
	localCounters *map[string]*string,
) {
	// Get the children
	children := *element.Children

	// Get the variable name
	name := *children[1].Value

	// Convert the type node to a TypeWrapper
	varType := types.ConvertToTypeWrapper(*children[2])

	// Get the expression
	expression := children[3]

	// Save the variable
	(*variables)[name] = fmt.Sprintf("x%d", *counter)

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
		expression,
		variables,
		traceCounters,
		usedStrings,
		usedArrays,
		usedNumbers,
		localCounters,
		true,
		&varType,
	)
}
