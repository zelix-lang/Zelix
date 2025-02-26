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
	"fmt"
	"strings"
)

func MarshalFunction(
	fun function.Function,
	trace *filecode.FileCode,
	traceFileName string,
	fileCodeId int,
	isMain bool,
	traceMagicCounter *int,
	traceCounters *map[int]string,
	usedStrings *pool.StringPool,
	poolExclusions *map[int]bool,
	fileTree *tree.InstructionTree,
	nameCounters *map[string]map[string]string,
	localCounters map[string]string,
) {
	// Keep a counter for all variables in the function
	// this is done to prevent name collisions with
	// internal variables injected by the IR generator
	// Example:
	// let a: str = "Hello, World!"
	// Converted to:
	// let x0: str = "Hello, World!"
	counter := pool.CounterPool{
		Exclusions: *poolExclusions,
	}

	// Keep in a map the variables used in the function
	// to retrieve their counter
	variables := make(map[string]string)

	// Construct the signature of the function
	signature := strings.Builder{}
	signature.WriteString("f ")
	signature.WriteString(fun.ReturnType.Marshal())
	signature.WriteString(" ")
	if isMain && fun.Name == "main" {
		signature.WriteString("main")
	} else {
		signature.WriteString(localCounters[fun.Name])
	}
	signature.WriteString(" ")

	paramCounter := 0
	// Write the parameters
	for _, param := range fun.Params {
		// Calculate the parameter's name
		name := fmt.Sprintf("p%d", paramCounter)
		variables[param.Name] = name

		// Write the parameter's name
		signature.WriteString(name)
		// Write the parameter's type
		signature.WriteString(" ")
		signature.WriteString(param.Type.Marshal())
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
				&counter,
				element,
				traceMagicCounter,
				variables,
				traceCounters,
				usedStrings,
				nameCounters,
			)
		default:
		}
	}

	// Write an end block to the function's signature
	signature.WriteString("end")
}
