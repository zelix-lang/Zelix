package engine

import (
	"fluent/ast"
	"fluent/ir/engine/function"
	"fluent/ir/engine/runtime"
	"fluent/ir/wrapper"
	"strings"
)

func MarshalIrWrapper(
	ir *wrapper.IrWrapper,
	fileCode *ast.FileCode,
	counter int,
) string {
	// Use a builder because string concatenation is slow
	builder := strings.Builder{}

	runtime.MarshalRuntime(ir, &builder)
	function.MarshalFunctions(ir, fileCode, &builder, counter)

	return builder.String()
}
