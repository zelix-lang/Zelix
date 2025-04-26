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

package property

import (
	"fluent/ast"
	"fluent/filecode"
	"fluent/filecode/module"
	"fluent/filecode/types/wrapper"
	"fluent/ir/pool"
	"fluent/ir/rule/call"
	"fluent/ir/tree"
	"fluent/ir/value"
	"fluent/ir/variable"
	"fluent/util"
	"fmt"
	"strconv"
	"strings"
)

var emptyString = ""

// MarshalPropertyAccess marshals property access expressions into the intermediate representation (IR).
// It processes the children of the given AST node and generates the corresponding IR instructions.
//
// Parameters:
// - global: The global instruction tree.
// - trace: The file code trace information.
// - child: The AST node representing the property access.
// - fileCodeId: The ID of the file code.
// - counter: A pointer to the counter for generating unique identifiers.
// - pair: The marshal pair containing the parent instruction tree.
// - modulePropCounters: A map of module property counters.
// - traceCounters: A pool of trace counters.
// - usedStrings: A pool of used strings.
// - usedNumbers: A pool of used numbers.
// - exprQueue: A queue of marshal pairs for expressions.
// - localCounters: A map of local counters.
// - traceFileName: The name of the trace file.
func MarshalPropertyAccess(
	global *tree.InstructionTree,
	trace *filecode.FileCode,
	child *ast.AST,
	fileCodeId int,
	counter *int,
	pair *tree.MarshalPair,
	variables *map[string]*variable.IRVariable,
	modulePropCounters *map[*module.Module]*util.OrderedMap[string, *string],
	traceCounters *pool.NumPool,
	usedStrings *pool.StringPool,
	usedNumbers *pool.StringPool,
	exprQueue *[]tree.MarshalPair,
	localCounters *map[string]*string,
	traceFileName string,
) {
	// Get the property's children
	children := *child.Children
	childrenLen := len(children) - 1

	// Save the necessary information for marshaling
	var lastMod *module.Module
	lastCandidateAddress := &emptyString
	lastExpression := tree.InstructionTree{
		Children:       &[]*tree.InstructionTree{},
		Representation: &strings.Builder{},
	}

	// Iterate over the children
	for i, child := range children {
		exprChildren := *child.Children
		expr := exprChildren[0]

		// Check if we are in the first iteration
		if i == 0 {
			candidate := tree.InstructionTree{
				Children:       &[]*tree.InstructionTree{},
				Representation: &strings.Builder{},
			}

			// Retrieve the inferred type
			var inferredType *wrapper.TypeWrapper
			if expr.Rule == ast.Expression {
				inferredType = expr.InferredType
			} else {
				inferredType = child.InferredType
			}

			lastMod = trace.Modules[inferredType.BaseType]

			// Retrieve variables if possible instead of cloning values
			if child.Children != nil && len(*child.Children) > 0 && (*child.Children)[0].Rule == ast.Identifier {
				// Get the variable
				storedVar := (*variables)[*(*child.Children)[0].Value]
				*lastCandidateAddress = fmt.Sprintf("%s ", storedVar.Addr)
			} else if value.RetrieveStaticVal(fileCodeId, child, candidate.Representation, usedStrings, usedNumbers) {
				// Check if we can save memory on the candidate
				*lastCandidateAddress = candidate.Representation.String()
			} else {
				// Get a suitable counter for this expression
				suitable := *counter
				*counter++

				// Append the candidate to the queue
				*exprQueue = append(*exprQueue, tree.MarshalPair{
					Child:       child,
					Parent:      &candidate,
					IsInline:    false,
					Counter:     suitable,
					MoveToStack: true,
					Expected:    *inferredType,
				})

				*global.Children = append([]*tree.InstructionTree{&candidate}, *global.Children...)
				*lastCandidateAddress = fmt.Sprintf("x%d ", suitable)
			}

			continue
		}

		// Determine the IR instruction
		var instruction string
		var representation *strings.Builder
		isCall := false
		switch expr.Rule {
		case ast.Identifier:
			instruction = "prop "
		case ast.FunctionCall:
			isCall = true
			instruction = "c "
		default:
		}

		isLast := i == childrenLen
		if isLast {
			representation = pair.Parent.Representation
		} else {
			representation = lastExpression.Representation

			// Clone the last expression to avoid memory issues
			lastExpressionClone := tree.InstructionTree{
				Children:       lastExpression.Children,
				Representation: lastExpression.Representation,
			}

			*global.Children = append([]*tree.InstructionTree{&lastExpressionClone}, *global.Children...)

			// Reset the last expression
			lastExpression = tree.InstructionTree{
				Children:       &[]*tree.InstructionTree{},
				Representation: &strings.Builder{},
			}
		}

		if lastMod == nil {
			panic("Cannot compile this IR - cannot find a module")
		}

		var suitable int
		if !isLast {
			// Get a suitable counter for this expression
			suitable = *counter
			*counter++

			representation.WriteString("mov x")
			representation.WriteString(strconv.Itoa(suitable))
			representation.WriteString(" &&")
		}

		// Write the appropriate format
		if isCall {
			props := (*modulePropCounters)[lastMod]
			// Get the function's name
			exprChildren = *expr.Children
			name := *exprChildren[0].Value

			if !isLast {
				// Get the method's type wrapper
				prop := lastMod.Functions[name]
				representation.WriteString(*(*localCounters)[prop.ReturnType.BaseType])
				representation.WriteString("\nstore x")
				representation.WriteString(strconv.Itoa(suitable))
				representation.WriteString(" ")
			}

			representation.WriteString(instruction)
			methodName, _ := props.Get(name)
			representation.WriteString(*methodName)
			representation.WriteString(" ")
			representation.WriteString(*lastCandidateAddress)
			representation.WriteString(" ")

			// Determine if this expression has parameters
			hasParams := len(exprChildren) > 1
			method := lastMod.Functions[name]

			if hasParams {
				// Get counters for the trace information
				lineCounter := traceCounters.RequestAddress(fileCodeId, child.Line)
				colCounter := traceCounters.RequestAddress(fileCodeId, child.Column)

				call.MarshalParams(
					method,
					*exprChildren[1].Children,
					counter,
					global,
					fileCodeId,
					pair.Parent,
					usedStrings,
					usedNumbers,
					exprQueue,
					lineCounter,
					colCounter,
					traceFileName,
				)
			}

			if !isLast {
				// Get the prop's type wrapper
				lastMod = trace.Modules[method.ReturnType.BaseType]
			}
		} else {
			if !isLast {
				// Get the prop's type wrapper
				prop := lastMod.Declarations[*expr.Value]
				representation.WriteString(*(*localCounters)[prop.Type.BaseType])
				representation.WriteString(" ")
			}

			representation.WriteString(instruction)
			representation.WriteString(*lastCandidateAddress)

			// Get the module's ordered map
			props := (*modulePropCounters)[lastMod]
			propCounter, _ := props.Get(*expr.Value)
			representation.WriteString(*propCounter)

			if !isLast {
				// Get the prop's type wrapper
				prop := lastMod.Declarations[*expr.Value]
				lastMod = trace.Modules[prop.Type.BaseType]
			}
		}

		if !isLast {
			*lastCandidateAddress = fmt.Sprintf("x%d ", suitable)
		}
	}
}
