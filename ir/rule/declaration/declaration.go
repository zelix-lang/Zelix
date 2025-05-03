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
	module2 "fluent/filecode/module"
	"fluent/filecode/types"
	"fluent/ir/pool"
	expression2 "fluent/ir/rule/expression"
	"fluent/ir/tree"
	"fluent/ir/variable"
	"fluent/util"
	"fmt"
)

// MarshalDeclaration marshals a declaration into the provided queue element.
// It processes the AST element, extracts the variable name, type, and expression,
// and saves the variable in the provided variables map. It then marshals the expression.
//
// Parameters:
// - queueElement: The queue element to marshal the declaration into.
// - trace: The file code trace.
// - fileCodeId: The ID of the file code.
// - traceFileName: The name of the trace file.
// - isMod: A boolean indicating if it is a module.
// - modulePropCounters: A map of module property counters.
// - traceFn: The function trace.
// - originalPath: The original path of the file.
// - counter: A pointer to the counter for generating unique addresses.
// - element: The AST element representing the declaration.
// - variables: A map of variables to update with the new variable.
// - traceCounters: A pool of number counters for tracing.
// - usedStrings: A pool of used strings.
// - localCounters: A map of local counters.
func MarshalDeclaration(
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
	(*variables)[name] = &variable.IRVariable{
		Addr: fmt.Sprintf("x%d", *counter),
		Type: &varType,
	}

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
		localCounters,
		true,
		false,
		&varType,
	)
}
