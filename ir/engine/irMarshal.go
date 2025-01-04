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
) string {
	// Use a builder because string concatenation is slow
	builder := strings.Builder{}

	runtime.MarshalRuntime(ir, &builder)
	function.MarshalFunctions(ir, fileCode, &builder)

	return builder.String()
}
