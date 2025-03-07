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
	"fluent/filecode/types/wrapper"
	"fluent/ir/pool"
	"fluent/ir/rule/expression"
	"fluent/ir/tree"
	"fluent/ir/value"
	"fluent/util"
	"strconv"
	"strings"
)

func MarshalReturn(
	funTree *tree.InstructionTree,
	trace *filecode.FileCode,
	fileCodeId int,
	traceFileName string,
	modulePropCounters *map[string]*util.OrderedMap[string, *string],
	counter *int,
	element *ast.AST,
	variables map[string]string,
	traceCounters *pool.NumPool,
	usedStrings *pool.StringPool,
	usedArrays *pool.StringPool,
	usedNumbers *pool.StringPool,
	nameCounters *map[string]map[string]string,
	localCounters *map[string]string,
	retType *wrapper.TypeWrapper,
) {
	children := *element.Children

	// Check if this return in empty
	if len(children) == 0 {
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
	if value.RetrieveStaticVal(fileCodeId, expr, &retTree, usedStrings, usedNumbers, variables) {
		funTree.Representation.WriteString(retTree.Representation.String())
		funTree.Representation.WriteString("\n")
		return
	}

	retTree.Representation.WriteString("x")
	retTree.Representation.WriteString(strconv.Itoa(*counter))

	// Marshal the expression
	expression.MarshalExpression(
		funTree,
		trace,
		fileCodeId,
		traceFileName,
		modulePropCounters,
		counter,
		expr,
		variables,
		traceCounters,
		usedStrings,
		usedArrays,
		usedNumbers,
		nameCounters,
		localCounters,
		true,
		retType,
	)

	// Write the instruction tree to global tree
	funTree.Representation.WriteString(retTree.Representation.String())
	funTree.Representation.WriteString("\n")
}
