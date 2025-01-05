package runtime

import (
	"fluent/ir/wrapper"
	"strings"
)

// MarshalRuntime adds the runtime instructions to the IR string
// based on the given IrWrapper
func MarshalRuntime(ir *wrapper.IrWrapper, builder *strings.Builder) {
	// Iterate over all runtime functions
	for name, functions := range ir.GetRuntimeFunctions() {
		builder.WriteString("runtime ")
		builder.WriteString(name)
		builder.WriteByte(' ')

		// Add the functions to the instruction
		for _, function := range functions {
			builder.WriteString(function)
			builder.WriteByte(' ')
		}

		builder.WriteByte('\n')
	}
}
