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
	"fluent/filecode/function"
	"fluent/ir/pool"
	"fluent/ir/tree"
	"fluent/ir/value"
	"strconv"
	"strings"
)

// MarshalParams marshals the parameters of a function call into the parent instruction tree.
// It processes each parameter, retrieves static values if needed, generates a suitable counter,
// and adds the expression to the queue. It also adds trace parameters to the parent representation.
//
// Parameters:
// - fun: The function whose parameters are being marshaled.
// - params: The AST nodes representing the parameters.
// - counter: A pointer to an integer counter used for generating unique identifiers.
// - global: The global instruction tree.
// - fileCodeId: The ID of the file code.
// - parent: The parent instruction tree.
// - variables: A map of variable names to IRVariable pointers.
// - usedStrings: A pool of used strings.
// - usedNumbers: A pool of used numbers.
// - exprQueue: A queue of expressions to be marshaled.
// - lineCounter: The line counter for trace information.
// - colCounter: The column counter for trace information.
// - traceFileName: The name of the trace file.
func MarshalParams(
	fun *function.Function,
	params []*ast.AST,
	counter *int,
	global *tree.InstructionTree,
	fileCodeId int,
	parent *tree.InstructionTree,
	usedStrings *pool.StringPool,
	usedNumbers *pool.StringPool,
	exprQueue *[]tree.MarshalPair,
	lineCounter string,
	colCounter string,
	traceFileName string,
) {
	i := 0

	// Add parameters
	for _, param := range fun.Params {
		// Get the param node
		paramNode := params[i]

		// Get the expression inside the param node
		expr := (*paramNode.Children)[0]

		// Retrieve the string literal if needed
		if value.RetrieveStaticVal(fileCodeId, expr, parent.Representation, usedStrings) {
			continue
		}

		// Generate a suitable counter
		*counter++
		suitable := *counter

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

	// Add trace params
	parent.Representation.WriteString(traceFileName)
	parent.Representation.WriteString(" ")
	parent.Representation.WriteString(lineCounter)
	parent.Representation.WriteString(" ")
	parent.Representation.WriteString(colCounter)
	parent.Representation.WriteString(" ")
}

// RequestTrace determines the trace information for a function call.
// It checks if the main function is the caller and requests trace information accordingly.
//
// Parameters:
// - traceFn: The function for which trace information is being requested.
// - originalPath: The original path of the function.
// - traceCounters: A pool of trace counters.
// - traceFileName: The name of the trace file.
// - fileCodeId: The ID of the file code.
// - line: The line number in the source code.
// - col: The column number in the source code.
// - isMod: A boolean indicating if the function is a module.
//
// Returns:
// - The line counter as a string.
// - The column counter as a string.
// - The trace file name as a string.
func RequestTrace(
	traceFn *function.Function,
	originalPath *string,
	traceCounters *pool.NumPool,
	traceFileName string,
	fileCodeId int,
	line int,
	col int,
	isMod bool,
) (string, string, string) {
	// Determine if the main function is the caller
	callerIsMain := !isMod && traceFn.Path == *originalPath && traceFn.Name == "main"

	// Determine trace information
	if callerIsMain {
		// Request trace information
		lineCounter := traceCounters.RequestAddress(fileCodeId, line)
		colCounter := traceCounters.RequestAddress(fileCodeId, col)

		return lineCounter, colCounter, traceFileName
	}

	return "__line", "__col", "__file"
}

// MarshalFunctionCall marshals a function call into the parent instruction tree.
// It writes the call instruction, determines if the main function is being called,
// adds trace information, and processes the function call parameters if present.
//
// Parameters:
// - global: The global instruction tree.
// - child: The AST node representing the function call.
// - traceFileName: The name of the trace file.
// - fileCodeId: The ID of the file code.
// - originalPath: The original path of the function.
// - isMod: A boolean indicating if the function is a module.
// - trace: The file code trace information.
// - traceFn: The function for which trace information is being requested.
// - counter: A pointer to an integer counter used for generating unique identifiers.
// - parent: The parent instruction tree.
// - traceCounters: A pool of trace counters.
// - usedStrings: A pool of used strings.
// - usedNumbers: A pool of used numbers.
// - exprQueue: A queue of expressions to be marshaled.
// - localCounters: A map of local counters for functions.
func MarshalFunctionCall(
	global *tree.InstructionTree,
	child *ast.AST,
	traceFileName string,
	fileCodeId int,
	originalPath *string,
	isMod bool,
	trace *filecode.FileCode,
	traceFn *function.Function,
	counter *int,
	parent *tree.InstructionTree,
	traceCounters *pool.NumPool,
	usedStrings *pool.StringPool,
	usedNumbers *pool.StringPool,
	exprQueue *[]tree.MarshalPair,
	localCounters *map[string]*string,
) {
	// Write the call instruction to the parent
	parent.Representation.WriteString("c ")

	// Get the call's children
	children := *child.Children
	funName := *children[0].Value

	// Determine if the main function is being called
	var isMain bool
	// Attempt to determine the function's counter
	fun := trace.Functions[funName]
	// Get the counter
	funCounter, ok := (*localCounters)[funName]

	if !ok {
		isMain = funName == "main"
		// External impl available, write the name directly
		parent.Representation.WriteString(funName)
	} else {
		parent.Representation.WriteString(*funCounter)
	}
	parent.Representation.WriteString(" ")

	// Add trace information
	var lineCounter string
	var colCounter string
	var traceFile string

	// Determine trace information
	if !isMain {
		lineCounter, colCounter, traceFile = RequestTrace(
			traceFn,
			originalPath,
			traceCounters,
			traceFileName,
			fileCodeId,
			child.Line,
			child.Column,
			isMod,
		)
	}

	// Determine if the function call has parameters
	hasParams := len(children) > 1

	if !hasParams {
		parent.Representation.WriteString(traceFile)
		parent.Representation.WriteString(" ")
		parent.Representation.WriteString(lineCounter)
		parent.Representation.WriteString(" ")
		parent.Representation.WriteString(colCounter)
		parent.Representation.WriteString(" ")
	} else {
		// Get the call's parameters
		params := *children[1].Children
		MarshalParams(
			fun,
			params,
			counter,
			global,
			fileCodeId,
			parent,
			usedStrings,
			usedNumbers,
			exprQueue,
			lineCounter,
			colCounter,
			traceFile,
		)
	}
}
