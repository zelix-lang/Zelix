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

package ret

import (
	"fluent/ast"
	"fluent/filecode"
	"fluent/filecode/function"
	"fluent/filecode/types/wrapper"
	"fluent/ir/pool"
	"fluent/ir/rule/expression"
	"fluent/ir/tree"
	"fluent/ir/value"
	"fluent/ir/variable"
	"fluent/util"
	"strconv"
	"strings"
)

func MarshalReturn(
	representation *strings.Builder,
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
	retType *wrapper.TypeWrapper,
) {
	children := *element.Children

	// Check if this return in empty
	if len(children) == 0 {
		representation.WriteString("ret_void\n")
		return
	}

	// Get the returned expression
	expr := children[0]

	// Create a new instruction tree for this return
	retTree := tree.InstructionTree{
		Children:       nil,
		Representation: &strings.Builder{},
	}

	retTree.Representation.WriteString("ret ")

	// See if we can save memory in the expression
	if value.RetrieveStaticVal(fileCodeId, expr, retTree.Representation, usedStrings, usedNumbers, variables) {
		representation.WriteString(retTree.Representation.String())
		representation.WriteString("\n")
		return
	}

	retTree.Representation.WriteString("x")
	retTree.Representation.WriteString(strconv.Itoa(*counter))

	// Marshal the expression
	expression.MarshalExpression(
		representation,
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
		retType,
	)

	// Write the instruction tree to global tree
	representation.WriteString(retTree.Representation.String())
	representation.WriteString("\n")
}
