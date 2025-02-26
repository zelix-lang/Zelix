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

package call

import (
	"fluent/ast"
	"fluent/filecode"
	"fluent/ir/pool"
	"fluent/ir/tree"
	"fmt"
	"strconv"
	"strings"
)

func MarshalFunctionCall(
	global *tree.InstructionTree,
	child *ast.AST,
	traceFileName string,
	fileCodeId int,
	trace *filecode.FileCode,
	traceMagicCounter *int,
	counter *pool.CounterPool,
	parent *tree.InstructionTree,
	traceCounters *map[int]string,
	nameCounters *map[string]map[string]string,
	usedStrings *pool.StringPool,
	exprQueue *[]tree.MarshalPair,
) {
	lineCounter, ok := (*traceCounters)[child.Line]

	if !ok {
		formattedCounter := fmt.Sprintf("__trace_magic_%d", *traceMagicCounter)
		(*traceCounters)[child.Line] = formattedCounter
		lineCounter = formattedCounter
		*traceMagicCounter++
	}

	colCounter, ok := (*traceCounters)[child.Column]

	if !ok {
		formattedCounter := fmt.Sprintf("__trace_magic_%d", *traceMagicCounter)
		(*traceCounters)[child.Column] = formattedCounter
		colCounter = formattedCounter
		*traceMagicCounter++
	}

	// Write the call instruction to the parent
	parent.Representation.WriteString("c ")

	// Get the call's children
	children := *child.Children
	funName := *children[0].Value

	// Attempt to determine the function's counter
	fun := trace.Functions[funName]
	// Get the counter
	funCounter, ok := (*nameCounters)[fun.Path][funName]

	if !ok {
		// External impl available, write the name directly
		parent.Representation.WriteString(funName)
	} else {
		parent.Representation.WriteString(funCounter)
	}
	parent.Representation.WriteString(" ")

	// Determine if the function call has parameters
	hasParams := len(children) > 1

	if hasParams {
		// Get the call's parameters
		params := *children[1].Children
		i := 0

		// Add parameters
		for _, param := range fun.Params {
			// Get the param node
			paramNode := params[i]

			// Get the expression inside the param node
			expr := (*paramNode.Children)[0]

			// Get the expression's children
			exprChildren := *expr.Children

			// Check if we can reuse a string
			if len(exprChildren) == 1 && exprChildren[0].Rule == ast.StringLiteral {
				strLiteral := exprChildren[0]
				parent.Representation.WriteString(
					usedStrings.RequestAddress(
						fileCodeId,
						*strLiteral.Value,
					),
				)

				parent.Representation.WriteString(" ")
				continue
			}

			// Generate a suitable counter
			suitable := counter.RequestSuitable()

			parent.Representation.WriteString("x")
			parent.Representation.WriteString(strconv.Itoa(suitable))
			parent.Representation.WriteString(" ")

			// Create a new InstructionTree
			instructionTree := tree.InstructionTree{
				Children:       &[]*tree.InstructionTree{},
				Representation: &strings.Builder{},
			}

			*global.Children = append([]*tree.InstructionTree{&instructionTree}, *global.Children...)

			// Add the expression to the queue
			*exprQueue = append(*exprQueue, tree.MarshalPair{
				Child:    expr,
				Parent:   &instructionTree,
				Counter:  suitable,
				Expected: param.Type,
				IsParam:  true,
			})

			i++
		}
	}

	// Add trace params
	parent.Representation.WriteString(traceFileName)
	parent.Representation.WriteString(" ")
	parent.Representation.WriteString(lineCounter)
	parent.Representation.WriteString(" ")
	parent.Representation.WriteString(colCounter)
	parent.Representation.WriteString(" ")
}
