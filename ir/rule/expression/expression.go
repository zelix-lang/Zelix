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

package expression

import (
	"fluent/ast"
	"fluent/filecode"
	queue2 "fluent/ir/queue"
	"fluent/ir/rule/call"
	"strconv"
	"strings"
)

func findOrInsertTraceCounter(
	traceCounters *map[int]int,
	val int,
	builder *strings.Builder,
	traceCounter *int,
) int {
	// Make sure to have a suitable value stored for the trace values
	storedVal, valAppended := (*traceCounters)[val]

	// Append the line to the counter if it not found
	if !valAppended {
		*traceCounter++
		storedVal = *traceCounter
		(*traceCounters)[val] = *traceCounter

		// Write ref instructions
		builder.WriteString("ref __trace_x")
		builder.WriteString(strconv.Itoa(storedVal))
		builder.WriteString(" num ")
		builder.WriteString(strconv.Itoa(val))
		builder.WriteString("\n")
	}

	return storedVal
}

func MarshalExpression(
	input *ast.AST,
	trace *filecode.FileCode,
	builder *strings.Builder,
	counter *int,
	traceFileVarName string,
	traceCounters *map[int]int,
	traceCounter *int,
) {
	// Use a queue to marshal nested expressions
	queue := make([]queue2.PendingIRMarshal, 0)
	queue = append(queue, queue2.PendingIRMarshal{
		Input: input,
	})

	for len(queue) > 0 {
		// Get the first element
		element := queue[0]
		queue = queue[1:]

		parent := element.Input
		children := *parent.Children
		input := children[0]

		// Make sure to have a suitable value for the trace of this expression
		line := findOrInsertTraceCounter(traceCounters, input.Line, builder, traceCounter)
		column := findOrInsertTraceCounter(traceCounters, input.Column, builder, traceCounter)

		// Get the item's rule
		rule := input.Rule

		switch rule {
		case ast.FunctionCall:
			call.MarshalFunctionCall(
				element,
				input,
				trace,
				builder,
				counter,
				&queue,
				traceFileVarName,
				line,
				column,
			)
		default:
		}
	}
}
