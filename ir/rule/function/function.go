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
	"fluent/filecode"
	function2 "fluent/filecode/function"
	"fluent/ir/rule/block"
	"strconv"
	"strings"
)

func MarshalFunction(
	builder *strings.Builder,
	name string,
	fun function2.Function,
	trace *filecode.FileCode,
	traceFileVarName string,
	traceCounters *map[int]int,
	traceCounter *int,
) {
	builder.WriteString("f ")
	builder.WriteString(name)
	builder.WriteString(" ")

	// Use a counter to make sure a name never repeats
	counter := 0

	// Write all parameters
	for _, param := range fun.Params {
		builder.WriteString(param.Name)
		builder.WriteString(" ")

		// Write all pointers
		for range param.Type.PointerCount {
			builder.WriteString("&")
		}

		builder.WriteString(param.Type.BaseType)
		builder.WriteString(" ")
		counter++
	}

	// Add compiler magic for tracing functions
	if !(fun.Name == "panic" && fun.IsStd) {
		builder.WriteString("__line")
		builder.WriteString(strconv.Itoa(counter))
		builder.WriteString(" num ")
		counter++

		builder.WriteString("__column")
		builder.WriteString(strconv.Itoa(counter))
		builder.WriteString(" num ")
		counter++

		builder.WriteString("__file")
		builder.WriteString(strconv.Itoa(counter))
		builder.WriteString(" num")
		counter++
	}

	builder.WriteString("\n")

	// Write the function's block
	block.MarshalBlock(
		&fun.Body,
		trace,
		builder,
		counter,
		traceFileVarName,
		traceCounters,
		traceCounter,
	)

	// Write the end block
	builder.WriteString("ef\n")
}
