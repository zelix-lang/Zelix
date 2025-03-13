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

package object

import (
	"fluent/ast"
	"fluent/filecode"
	"fluent/ir/pool"
	"fluent/ir/rule/call"
	"fluent/ir/tree"
	"fluent/ir/value"
	"fluent/util"
	"fmt"
	"strconv"
	"strings"
)

func MarshalObjectCreation(
	global *tree.InstructionTree,
	child *ast.AST,
	traceFileName string,
	fileCodeId int,
	trace *filecode.FileCode,
	modulePropCounters *map[string]*util.OrderedMap[string, *string],
	counter *int,
	pair *tree.MarshalPair,
	traceCounters *pool.NumPool,
	variables *map[string]string,
	usedStrings *pool.StringPool,
	usedArrays *pool.StringPool,
	usedNumbers *pool.StringPool,
	exprQueue *[]tree.MarshalPair,
	localCounters *map[string]*string,
) {
	children := *child.Children

	// Get the module's name
	realName := *children[0].Value

	// Use a breadth-first queue to marshal modules
	queue := []tree.ModMarshalPair{
		{
			Name:            realName,
			Parent:          pair.Parent,
			CallConstructor: true,
			IsParam:         pair.IsParam,
			Counter:         pair.Counter,
		},
	}

	for len(queue) > 0 {
		// Get the first element
		element := queue[0]
		queue = queue[1:]

		realName := element.Name
		name := (*localCounters)[realName]
		mod := trace.Modules[realName]

		// Get the prop counter
		propCounter := (*modulePropCounters)[realName]

		// Determine if the module has a constructor
		constructorName, found := propCounter.Get(realName)
		modAddress := ""

		if found {
			// See if this value is being moved to the stack
			if element.IsParam {
				modAddress = fmt.Sprintf("x%d", element.Counter)
			} else {
				suitable := *counter
				modAddress = fmt.Sprintf("x%d", suitable)
				*counter++

				// Move this value to the stack
				element.Parent.Representation.WriteString("mov ")
				element.Parent.Representation.WriteString(modAddress)
				element.Parent.Representation.WriteString(" ")
				element.Parent.Representation.WriteString(*name)
				element.Parent.Representation.WriteString(" ")
			}
		}

		// Write the construct instructions
		element.Parent.Representation.WriteString("co ")
		element.Parent.Representation.WriteString(*name)
		element.Parent.Representation.WriteString(" ")

		// Iterate over the prop counters
		propCounter.Iterate(func(name string, pCounter *string) bool {
			// Get the property by name
			prop, ok := mod.Declarations[name]

			// Break the iteration if the prop was not found
			if !ok {
				return true
			}

			if prop.IsIncomplete {
				var suitable int
				firstVal := *counter
				propType := prop.Type
				baseType := propType.BaseType
				var instructionTree *tree.InstructionTree

				// Change the base type for non-primitive types
				if !propType.IsPrimitive {
					baseType = *(*localCounters)[propType.BaseType]
				}

				if propType.PointerCount > 0 {
					arrayRepresentation := strings.Repeat("[]", propType.ArrayCount)

					for i := propType.PointerCount; i > 0; i-- {
						pointerRepresentation := strings.Repeat("&", i)

						// Get a suitable counter for this expression
						suitable = *counter
						suitableStr := strconv.Itoa(suitable)
						*counter++

						if i == 1 {
							element.Parent.Representation.WriteString("x")
							element.Parent.Representation.WriteString(strconv.Itoa(firstVal))
							element.Parent.Representation.WriteString(" ")
						}

						// Create a local tree to represent the data
						localTree := tree.InstructionTree{
							Children:       &[]*tree.InstructionTree{},
							Representation: &strings.Builder{},
						}

						localTree.Representation.WriteString("mov x")
						localTree.Representation.WriteString(suitableStr)
						localTree.Representation.WriteString(" ")
						localTree.Representation.WriteString(pointerRepresentation)
						localTree.Representation.WriteString(baseType)
						localTree.Representation.WriteString(arrayRepresentation)
						localTree.Representation.WriteString(" addr ")

						*global.Children = append(
							[]*tree.InstructionTree{&localTree},
							*global.Children...,
						)

						if i != 1 {
							// Take the address of the next value
							localTree.Representation.WriteString("x")
							localTree.Representation.WriteString(strconv.Itoa(*counter))
							localTree.Representation.WriteString(" ")
						}

						instructionTree = &localTree
					}
				} else {
					instructionTree = element.Parent
				}

				if propType.ArrayCount > 0 {
					// Request an address for this type
					addr := usedArrays.RequestAddress(
						fileCodeId,
						baseType,
					)

					instructionTree.Representation.WriteString(addr)
					return false
				}

				if propType.IsPrimitive {
					switch propType.BaseType {
					case "num", "bool":
						instructionTree.Representation.WriteString("__FALSE ")
					case "dec":
						// Request an address for 0.0
						addr := usedNumbers.RequestAddress(fileCodeId, "0.0")
						instructionTree.Representation.WriteString(addr)
						instructionTree.Representation.WriteString(" ")
					case "str":
						// Request an address for an empty string
						addr := usedStrings.RequestAddress(fileCodeId, "")
						instructionTree.Representation.WriteString(addr)
						instructionTree.Representation.WriteString(" ")
					}
					return false
				}

				// Add the expression tree to the queue
				queue = append(queue, tree.ModMarshalPair{
					Name:            propType.BaseType,
					Parent:          instructionTree,
					CallConstructor: false,
					IsParam:         true,
					Counter:         suitable,
				})

			} else {
				// See if we can save memory on this value
				val := prop.Value
				if value.RetrieveStaticVal(fileCodeId, val, element.Parent.Representation, usedStrings, usedNumbers, variables) {
					return false
				}

				// Get a suitable counter for this value
				suitable := *counter
				*counter++

				// Write the new memory space
				element.Parent.Representation.WriteString("x")
				element.Parent.Representation.WriteString(strconv.Itoa(suitable))
				element.Parent.Representation.WriteString(" ")

				// Create a local tree to represent the data
				localTree := tree.InstructionTree{
					Children:       &[]*tree.InstructionTree{},
					Representation: &strings.Builder{},
				}

				// Add the tree to the stack
				*global.Children = append(
					[]*tree.InstructionTree{&localTree},
					*global.Children...,
				)

				// Schedule the expression
				*exprQueue = append(*exprQueue, tree.MarshalPair{
					Child:    val,
					Parent:   &localTree,
					Counter:  suitable,
					Expected: prop.Type,
					IsParam:  true,
				})
			}

			return false
		})

		if element.CallConstructor && found {
			// Since nodes are written before the representation itself
			// we have to write further instructions to the parent's
			// representation to avoid the call to the constructor
			// appearing before the instruction that constructs
			// the module
			element.Parent.Representation.WriteString("\n")

			// Insert a call to the constructor
			element.Parent.Representation.WriteString("c ")
			element.Parent.Representation.WriteString(*constructorName)
			element.Parent.Representation.WriteString(" ")
			element.Parent.Representation.WriteString(modAddress)
			element.Parent.Representation.WriteString(" ")

			lineAddress := traceCounters.RequestAddress(fileCodeId, child.Line)
			colAddress := traceCounters.RequestAddress(fileCodeId, child.Column)

			// Check if we have parameters
			if len(children) == 1 {
				element.Parent.Representation.WriteString(traceFileName)
				element.Parent.Representation.WriteString(" ")
				element.Parent.Representation.WriteString(lineAddress)
				element.Parent.Representation.WriteString(" ")
				element.Parent.Representation.WriteString(colAddress)
				element.Parent.Representation.WriteString(" ")
				return
			}

			// Get the parameters node
			paramsNode := children[1]
			params := *paramsNode.Children

			// Get the constructor
			constructor := mod.Functions[realName]

			// Marshal the parameters
			call.MarshalParams(
				constructor,
				params,
				counter,
				global,
				fileCodeId,
				element.Parent,
				variables,
				usedStrings,
				usedNumbers,
				exprQueue,
				lineAddress,
				colAddress,
				traceFileName,
			)
		}
	}

}
