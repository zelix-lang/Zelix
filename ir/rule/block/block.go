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

package block

import (
	"fluent/ast"
	"fluent/filecode"
	expression2 "fluent/ir/rule/expression"
	"strings"
)

func MarshalBlock(
	block *ast.AST,
	trace *filecode.FileCode,
	builder *strings.Builder,
	counter int,
	traceFileVarName string,
	traceCounters *map[int]int,
	traceCounter *int,
) {
	// Use a queue to marshal nested blocks
	queue := make([]*ast.AST, 0)
	queue = append(queue, block)

	for len(queue) > 0 {
		// Get the first element
		block := queue[0]
		queue = queue[1:]

		for _, expression := range *block.Children {
			// Handle nested blocks
			if expression.Rule == ast.Block {
				queue = append(queue, expression)
				continue
			}

			expression2.MarshalExpression(
				expression,
				trace,
				builder,
				&counter,
				traceFileVarName,
				traceCounters,
				traceCounter,
			)
		}
	}
}
