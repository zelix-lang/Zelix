package function

import (
	"fluent/ast"
	"fluent/ir/wrapper"
	"fluent/stack"
	"fluent/token"
	"fluent/tokenUtil/splitter"
	"strings"
)

// MarshalFunctions marshals the functions of the IR
func MarshalFunctions(
	ir *wrapper.IrWrapper,
	fileCode *ast.FileCode,
	builder *strings.Builder,
	counter int,
	variables *stack.Stack,
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

		// Used to skip tokens
		skipToIndex := 0

		// Write the function's body
		body := function.GetBody()
		for i, unit := range body {
			if i < skipToIndex {
				continue
			}

			tokenType := unit.GetType()

			switch tokenType {
			case token.If:
			case token.ElseIf:
			case token.Else:
			case token.While:
			case token.For:
			case token.Return:
			case token.New:
			case token.Identifier:
				statement, _ := splitter.ExtractTokensBefore(
					body[i:],
					token.Semicolon,
					false,
					token.Unknown,
					token.Unknown,
					true,
				)

				MarshalStatement(
					statement,
					builder,
					&counter,
					ir,
					fileCode,
					&skipToIndex,
					variables,
				)
			default:
				continue
			}
		}

		// Add ret_main for the main function
		if realName == "main" {
			builder.WriteString("ret_main\n")
		}

		// Write the end of the function
		builder.WriteString("endf ")
		builder.WriteString(realName)
		builder.WriteByte('\n')
	}
}
