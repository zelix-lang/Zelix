package function

import (
	"fluent/ast"
	"fluent/code/types"
	wrapper2 "fluent/code/wrapper"
	builder2 "fluent/ir/engine/builder"
	"fluent/ir/wrapper"
	"fluent/stack"
	"fluent/token"
	"fluent/tokenUtil/splitter"
	"strconv"
	"strings"
)

// MarshalStatement marshals the given statement into the IR
func MarshalStatement(
	statement []token.Token,
	builder *strings.Builder,
	counter *int,
	ir *wrapper.IrWrapper,
	fileCode *ast.FileCode,
	skipToIndex *int,
	variables *stack.Stack,
) wrapper2.FluentObject {
	//startAt := 0
	lastValue := wrapper2.NewFluentObject(
		wrapper2.ForceNewTypeWrapper(
			"nothing",
			make([]wrapper2.TypeWrapper, 0),
			types.NothingType,
		),
		nil,
	)

	// Use the first token to figure out what kind of statement this is
	firstToken := statement[0]

	// See if it's a defined function
	fun, funFound, _ := ast.LocateFunction(
		*fileCode.GetFunctions(),
		firstToken.GetFile(),
		firstToken.GetValue(),
	)

	if funFound {
		argsRaw, _ := splitter.ExtractTokensBefore(
			// Skip the function name
			statement[1:],
			token.CloseParen,
			true,
			token.OpenParen,
			token.CloseParen,
			true,
		)

		// Split by commas
		argsSplit := splitter.SplitTokens(
			// Skip the open parenthesis
			argsRaw[1:],
			token.Comma,
			token.OpenParen,
			token.CloseParen,
		)

		// Store the arguments in a map to construct the function call
		argsMap := make(map[int]int)

		// Iterate over the arguments and construct them
		for i, arg := range argsSplit {
			*counter++
			argsMap[i] = *counter
			object := builder2.BuildObject(arg)

			// Add the object to the IR
			builder.WriteString("mov x")
			builder.WriteString(strconv.Itoa(*counter))
			builder.WriteString(" ")

			// Transpile the type
			typeWrapper := object.GetType()
			switch typeWrapper.GetType() {
			case types.StringType:
				builder.WriteString("str ")
				// Write the value
				builder.WriteString(object.GetValue().(string))
			case types.IntType:
				builder.WriteString("i32 ")
				// Write the value
				builder.WriteString(strconv.Itoa(object.GetValue().(int)))
			case types.DecimalType:
				builder.WriteString("f32 ")
				// Write the value
				builder.WriteString(strconv.FormatFloat(object.GetValue().(float64), 'f', -1, 32))
			default:
				break
			}

			// Write a newline
			builder.WriteByte('\n')
		}

		// Parse the arguments of the function
		// Check if this function is a runtime function
		runtimePath, isRuntime := ir.GetRuntimeFunction(fun)

		if isRuntime {
			// Add s_c instruction (standard call)
			builder.WriteString("s_c ")
			builder.WriteString(runtimePath)
			builder.WriteByte(' ')
			builder.WriteString(firstToken.GetValue())
			builder.WriteByte(' ')
		} else {
			// Retrieve the computed function
			function := ir.GetFunction(fun)

			// Write the function name
			builder.WriteString("call ")
			builder.WriteString(function)
			builder.WriteByte(' ')
		}

		// Add all arguments
		for i := range argsSplit {
			// Retrieve the counter for this argument
			counterNum := argsMap[i]

			// Add the argument
			builder.WriteString("x")
			builder.WriteString(strconv.Itoa(counterNum))

			if i < len(argsSplit)-1 {
				builder.WriteByte(' ')
			}
		}

		builder.WriteByte('\n')
		// Destroy the arguments from the registry
		for i := range argsSplit {
			// Retrieve the counter for this argument
			counterNum := argsMap[i]

			// Add the argument
			builder.WriteString("end x")
			builder.WriteString(strconv.Itoa(counterNum))
			builder.WriteByte('\n')
		}
	} else {

	}

	return lastValue
}
