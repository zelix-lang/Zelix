package function

import (
	"fluent/ast"
	"fluent/ir/wrapper"
	"strings"
)

// MarshalFunctions marshals the functions of the IR
func MarshalFunctions(
	ir *wrapper.IrWrapper,
	fileCode *ast.FileCode,
	builder *strings.Builder,
) {
	// Get all functions
	functions := ir.GetFunctions()

	// Iterate over all functions
	for function, name := range functions {
		// Retrieve the function's real name
		realName := function.GetName()

		builder.WriteString("defn ")

		if realName != "main" {
			realName = name
		}

		// Change the name
		builder.WriteString(realName)

		builder.WriteByte(' ')
		// Write the return type
		returnType := function.GetReturnType()
		builder.WriteString(returnType.Marshal())

		// Write the function's arguments
		builder.WriteByte(' ')

		arguments := function.GetParameters()
		for i, arg := range arguments {
			argType := arg.GetType()
			builder.WriteString(argType.Marshal())
			builder.WriteByte(' ')
			builder.WriteString(arg.GetName())

			// Write a space if this is not the last argument
			if i < len(arguments)-1 {
				builder.WriteByte(' ')
			}
		}

		builder.WriteByte('\n')

		// Write the function's body
		// todo!()

		// Write the end of the function
		builder.WriteString("endf ")
		builder.WriteString(realName)
	}
}
