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
	"fluent/ir/queue"
	"strconv"
	"strings"
)

func MarshalFunctionCall(
	element queue.PendingIRMarshal,
	input *ast.AST,
	trace *filecode.FileCode,
	builder *strings.Builder,
	counter *int,
	exprQueue *[]queue.PendingIRMarshal,
	traceFileVarName string,
	lineTraceValName int,
	columnTraceValName int,
) {
	// Get the function call's children
	children := *input.Children

	// See if the call's parameters have been processed
	if !element.HasProcessedParams {
		// Determine if the function has parameters
		hasParams := len(children) > 1

		if hasParams {
			paramsNode := children[1]
			originalCounter := *counter

			// Schedule all parameters for marshaling
			for _, param := range *paramsNode.Children {
				*counter++

				children := *param.Children
				*exprQueue = append(*exprQueue, queue.PendingIRMarshal{
					Input:   children[0],
					IsParam: true,
					Counter: *counter,
				})
			}

			// Schedule the function call again
			*exprQueue = append(*exprQueue, queue.PendingIRMarshal{
				Input:              element.Input,
				HasProcessedParams: true,
				Counter:            originalCounter,
			})

			return
		}
	}

	// Get the function's name
	funName := *children[0].Value

	// Get the function from the FileCode
	fun := trace.Functions[funName]

	builder.WriteString("c ")
	builder.WriteString(funName)
	builder.WriteString(" ")

	// Add the counter of the parameters
	for i := 0; i < len(fun.Params); i++ {
		builder.WriteString("x")
		builder.WriteString(strconv.Itoa(i + element.Counter + 1))
		builder.WriteString(" ")
	}

	// Add the trace
	builder.WriteString("__trace_x")
	builder.WriteString(strconv.Itoa(lineTraceValName))
	builder.WriteString(" ")

	builder.WriteString("__trace_x")
	builder.WriteString(strconv.Itoa(columnTraceValName))
	builder.WriteString(" ")

	builder.WriteString(traceFileVarName)
	builder.WriteString("\n")
}
