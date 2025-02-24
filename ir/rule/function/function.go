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
	"fluent/filecode/function"
	"fluent/ir/tree"
	"fmt"
	"strings"
)

func MarshalFunction(
	fun function.Function,
	traceCounters *map[int]int,
	usedStrings *map[string]string,
	fileTree *tree.InstructionTree,
	localCounters map[string]string,
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

	// Construct the signature of the function
	signature := strings.Builder{}
	signature.WriteString("f ")
	signature.WriteString(localCounters[fun.Name])
	signature.WriteString(" ")

	// Write the parameters
	for _, param := range fun.Params {
		// Calculate the parameter's name
		name := fmt.Sprintf("p%d", counter)
		variables[param.Name] = name

		// Write the parameter's name
		signature.WriteString(name)
		// Write the parameter's type
		signature.WriteString(" ")
		signature.WriteString(param.Type.Marshal())
		signature.WriteString(" ")

		// Increment the counter
		counter++
	}

	// Write a newline to the signature
	signature.WriteString("\n")

	// Create a new InstructionTree for the function
	body := make([]*tree.InstructionTree, 0)
	funTree := tree.InstructionTree{
		Representation: signature.String(),
		Children:       &body,
		IsSignature:    true,
	}

	// Add the function tree node to the global tree
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
		case ast.Expression:

		default:
		}
	}

	// Write an end block to the function's signature
	funTree.Representation += "\nend"
}
