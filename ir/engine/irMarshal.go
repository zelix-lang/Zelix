package engine

import (
	"fluent/ir/engine/runtime"
	"fluent/ir/wrapper"
	"strings"
)

func MarshalIrWrapper(ir *wrapper.IrWrapper) string {
	// Use a builder because string concatenation is slow
	builder := strings.Builder{}

	runtime.MarshalRuntime(ir, &builder)

	return builder.String()
}
