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

package function

import (
	"fluent/ast"
	"fluent/filecode"
	"fluent/filecode/function"
	"fluent/ir/pool"
	"fluent/ir/rule/expression"
	"fluent/ir/tree"
	"fluent/util"
	"fmt"
	"strings"
)

func MarshalFunction(
	fun *function.Function,
	trace *filecode.FileCode,
	modType string,
	injectThis bool,
	traceFileName string,
	fileCodeId int,
	isMain bool,
	modulePropCounters *map[string]*util.OrderedMap[string, *string],
	traceCounters *pool.NumPool,
	usedStrings *pool.StringPool,
	usedNumbers *pool.StringPool,
	fileTree *tree.InstructionTree,
	nameCounters *map[string]map[string]string,
	name string,
	localCounters *map[string]string,
) {
	// Keep a counter for all variables in the function
	// this is done to prevent name collisions with
	// internal variables injected by the IR generator
	// Example:
	// let a: str = "Hello, World!"
	// Converted to:
	// let x0: str = "Hello, World!"
	counter := 0

	// Keep in a map the variables used in the function
	// to retrieve their counter
	variables := make(map[string]string)

	// Inject the "this" variable if this function belongs to a module
	if injectThis {
		variables["this"] = "p0"
	}

	// Construct the signature of the function
	signature := strings.Builder{}
	signature.WriteString("f ")
	signature.WriteString(fun.ReturnType.Marshal())
	signature.WriteString(" ")
	if isMain && fun.Name == "main" {
		signature.WriteString("main")
	} else {
		signature.WriteString(name)
	}
	signature.WriteString(" ")

	paramCounter := 0

	// Inject the "this" parameter if needed
	if injectThis {
		signature.WriteString("p0 ")
		signature.WriteString(modType)
		signature.WriteString(" ")
		paramCounter++
	}

	// Write the parameters
	for _, param := range fun.Params {
		// Calculate the parameter's name
		name := fmt.Sprintf("p%d", paramCounter)
		variables[param.Name] = name

		// Write the parameter's name
		signature.WriteString(name)
		// Write the parameter's type
		signature.WriteString(" ")

		if param.Type.IsPrimitive {
			signature.WriteString(param.Type.Marshal())
		} else {
			oldBaseType := param.Type.BaseType
			param.Type.BaseType = (*localCounters)[oldBaseType]
			signature.WriteString(param.Type.Marshal())
			param.Type.BaseType = oldBaseType
		}

		signature.WriteString(" ")
		paramCounter++
	}

	// Write trace parameters
	signature.WriteString("__file str __line str __col str")

	// Write a newline to the signature
	signature.WriteString("\n")

	// Create a new InstructionTree for the function
	body := make([]*tree.InstructionTree, 0)
	funTree := tree.InstructionTree{
		Representation: &signature,
		Children:       &body,
	}

	*fileTree.Children = append(*fileTree.Children, &funTree)

	// Use a queue to marshal blocks
	// in a breadth-first manner
	blockQueue := []*ast.AST{&fun.Body}

	for len(blockQueue) > 0 {
		// Get the first block in the queue
		element := blockQueue[0]
		blockQueue = blockQueue[1:]

		// See the element's rule
		rule := element.Rule

		switch rule {
		case ast.Block:
			// Add the block's children to the queue
			blockQueue = append(blockQueue, *element.Children...)
		case ast.Continue, ast.Break:
			// Directly write the tree's value
			funTree.Representation.WriteString(*element.Value)
			funTree.Representation.WriteString("\n")
		case ast.Expression:
			expression.MarshalExpression(
				&funTree,
				trace,
				fileCodeId,
				traceFileName,
				modulePropCounters,
				&counter,
				element,
				variables,
				traceCounters,
				usedStrings,
				usedNumbers,
				nameCounters,
				localCounters,
			)
		default:
		}
	}

	// Write an end block to the function's signature
	signature.WriteString("end")
}
