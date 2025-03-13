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
	"fluent/ir/rule/conditional"
	"fluent/ir/rule/declaration"
	"fluent/ir/rule/expression"
	"fluent/ir/rule/loop"
	"fluent/ir/rule/ret"
	"fluent/ir/tree"
	"fluent/ir/variable"
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
	isMod bool,
	originalPath *string,
	modulePropCounters *map[string]*util.OrderedMap[string, *string],
	traceCounters *pool.NumPool,
	usedStrings *pool.StringPool,
	usedArrays *pool.StringPool,
	usedNumbers *pool.StringPool,
	fileTree *tree.InstructionTree,
	name *string,
	localCounters *map[string]*string,
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
	variables := make(map[string]*variable.IRVariable)

	// Keep track of the appended blocks
	appendedBlocks := pool.BlockPool{
		Storage: make(map[string]*strings.Builder),
	}

	// Inject the "this" variable if this function belongs to a module
	if injectThis {
		variables["this"] = &variable.IRVariable{
			Addr: "p0",
			// "this" doesn't need a type, as it is not reassigned
			Type: nil,
		}
	}

	// Construct the signature of the function
	signature := strings.Builder{}
	signature.WriteString("f ")

	if fun.ReturnType.IsPrimitive {
		signature.WriteString(fun.ReturnType.Marshal())
	} else {
		signature.WriteString(*(*localCounters)[fun.ReturnType.BaseType])
	}

	signature.WriteString(" ")
	if isMain && fun.Name == "main" {
		signature.WriteString("main")
	} else {
		signature.WriteString(*name)
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
		variables[param.Name] = &variable.IRVariable{
			Addr: name,
			// Parameters don't need a type, as they are not reassigned
			Type: nil,
		}

		// Write the parameter's name
		signature.WriteString(name)
		// Write the parameter's type
		signature.WriteString(" ")

		if param.Type.IsPrimitive {
			signature.WriteString(param.Type.Marshal())
		} else {
			oldBaseType := param.Type.BaseType
			param.Type.BaseType = *(*localCounters)[oldBaseType]
			signature.WriteString(param.Type.Marshal())
			param.Type.BaseType = oldBaseType
		}

		signature.WriteString(" ")
		paramCounter++
	}

	// Write trace parameters
	if !(isMain && name == nil) {
		signature.WriteString("__file str __line str __col str")
	}

	// Write a newline to the signature
	signature.WriteString("\n")

	// Create a new InstructionTree for the function
	body := make([]*tree.InstructionTree, 0)
	funTree := tree.InstructionTree{
		Representation: &signature,
		Children:       &body,
	}

	*fileTree.Children = append(*fileTree.Children, &funTree)

	// Use a queue to marshal blocks in a breadth-first manner
	blockQueue := []*tree.BlockMarshalElement{
		{
			Element:        &fun.Body,
			Representation: funTree.Representation,
		},
	}

	for len(blockQueue) > 0 {
		// Get the first block in the queue
		queueElement := blockQueue[0]
		blockQueue = blockQueue[1:]
		element := queueElement.Element

		// See the element's rule
		rule := element.Rule

		// Used to store whether the current element has a nested block
		hasNested := false

		switch rule {
		case ast.If:
			hasNested = true
			conditional.MarshalIf(
				queueElement,
				trace,
				fileCodeId,
				traceFileName,
				isMod,
				modulePropCounters,
				fun,
				originalPath,
				&counter,
				element,
				&variables,
				traceCounters,
				&appendedBlocks,
				usedStrings,
				usedArrays,
				usedNumbers,
				localCounters,
				&blockQueue,
			)
		case ast.For:
			hasNested = true
			loop.MarshalFor(
				queueElement,
				trace,
				fileCodeId,
				traceFileName,
				isMod,
				modulePropCounters,
				&counter,
				fun,
				originalPath,
				element,
				&variables,
				traceCounters,
				&appendedBlocks,
				usedStrings,
				usedArrays,
				usedNumbers,
				localCounters,
				&blockQueue,
			)
		case ast.While:
			hasNested = true
			loop.MarshalWhile(
				queueElement,
				trace,
				fileCodeId,
				traceFileName,
				isMod,
				modulePropCounters,
				fun,
				originalPath,
				&counter,
				element,
				&variables,
				traceCounters,
				&appendedBlocks,
				usedStrings,
				usedArrays,
				usedNumbers,
				localCounters,
				&blockQueue,
			)
		case ast.Block:
			hasNested = true
			childrenLen := len(*element.Children) - 1
			// Add the block's children to the queue
			for i, el := range *element.Children {
				blockQueue = append(blockQueue, &tree.BlockMarshalElement{
					Element:        el,
					Representation: queueElement.Representation,
					ParentAddr:     queueElement.ParentAddr,
					IsLast:         i == childrenLen,
					JumpToParent:   queueElement.JumpToParent,
					Id:             queueElement.Id,
					RemainingAddr:  queueElement.RemainingAddr,
				})
			}

			// Update the flag if there are no children
			if childrenLen == -1 {
				queueElement.IsLast = true
				hasNested = false
			}
		case ast.Continue:
			funTree.Representation.WriteString("jump ")
			funTree.Representation.WriteString(*queueElement.ParentAddr)
			funTree.Representation.WriteString("\n")
		case ast.Break:
			funTree.Representation.WriteString("jump ")
			funTree.Representation.WriteString(*queueElement.RemainingAddr)
			funTree.Representation.WriteString("\n")
		case ast.Declaration:
			declaration.MarshalDeclaration(
				queueElement,
				trace,
				fileCodeId,
				traceFileName,
				isMod,
				modulePropCounters,
				fun,
				originalPath,
				&counter,
				element,
				&variables,
				traceCounters,
				usedStrings,
				usedArrays,
				usedNumbers,
				localCounters,
			)
		case ast.Expression:
			expression.MarshalExpression(
				queueElement.Representation,
				trace,
				fun,
				fileCodeId,
				isMod,
				traceFileName,
				originalPath,
				modulePropCounters,
				&counter,
				element,
				&variables,
				traceCounters,
				usedStrings,
				usedArrays,
				usedNumbers,
				localCounters,
				false,
				nil,
			)
		case ast.Return:
			ret.MarshalReturn(
				queueElement.Representation,
				trace,
				fileCodeId,
				traceFileName,
				isMod,
				modulePropCounters,
				fun,
				originalPath,
				&counter,
				element,
				&variables,
				traceCounters,
				usedStrings,
				usedArrays,
				usedNumbers,
				localCounters,
				&fun.ReturnType,
			)
		default:
		}

		// Jump to the parent if needed
		if !hasNested && queueElement.IsLast && queueElement.JumpToParent {
			queueElement.Representation.WriteString("jump ")
			queueElement.Representation.WriteString(*queueElement.ParentAddr)
			queueElement.Representation.WriteString("\n")
		} else if !hasNested && queueElement.IsLast && queueElement.RemainingAddr != nil {
			queueElement.Representation.WriteString("jump ")
			queueElement.Representation.WriteString(*queueElement.RemainingAddr)
			queueElement.Representation.WriteString("\n")
		}
	}

	// Write nested blocks
	for address, block := range appendedBlocks.Storage {
		signature.WriteString("block ")
		signature.WriteString(address)
		signature.WriteString("\n")
		signature.WriteString(block.String())
		signature.WriteString("end\n")
	}

	// Add ret_void instructions if the function does not return anything
	if fun.ReturnType.BaseType == "nothing" {
		signature.WriteString("ret_void\n")
	}

	// Write an end block to the function's signature
	signature.WriteString("end")
}
