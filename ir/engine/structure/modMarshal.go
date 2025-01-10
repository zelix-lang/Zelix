package structure

import (
	"fluent/ast"
	"fluent/code/mod"
	"fluent/code/types"
	wrapper2 "fluent/code/wrapper"
	function2 "fluent/ir/engine/function"
	"fluent/ir/wrapper"
	"fluent/stack"
	"strconv"
	"strings"
)

// MarshalSingleMod marshals a single mod into the given strings.Builder
func MarshalSingleMod(
	module *mod.FluentMod,
	ir *wrapper.IrWrapper,
	builder *strings.Builder,
	counter *int,
	fileCode *ast.FileCode,
	comment string,
) {
	// Build variables that have mods with generics
	props := module.GetVarDeclarations()
	propWrappers := make([]wrapper2.TypeWrapper, 0)

	for _, prop := range props {
		resultantWrapper := MarshalGenericsInVar(
			prop,
			ir,
			builder,
			counter,
			fileCode,
		)

		propWrappers = append(propWrappers, resultantWrapper)
	}

	// Find how many times this mod has been built
	modCounter, _ := ir.GetMod(module)

	if modCounter == -1 {
		*counter++
		ir.AddMod(module, *counter)
		modCounter = *counter
	}

	// Write the mod with its counter
	builder.WriteString("mod x")
	builder.WriteString(strconv.Itoa(modCounter))
	builder.WriteString(" ")

	// Marshal properties
	for propCounter, prop := range props {
		ir.AddModProp(module, prop[1].GetValue(), propCounter)

		// Find the wrapper for this property
		propWrapper := propWrappers[propCounter]

		MarshalType(propWrapper, builder)
		builder.WriteString(" ")

		// Write the value
		isMod := propWrapper.GetType() == types.ModType

		if isMod {
			var modToFind *mod.FluentMod
			genericMod, isGeneric := ir.GetGenericMod(propWrapper.Marshal())

			if isGeneric {
				modToFind = genericMod
			} else {
				modToFind, _, _ = mod.FindMod(
					fileCode.GetModules(),
					propWrapper.GetBaseType(),
					prop[1].GetFile(),
				)
			}

			propModCounter, _ := ir.GetMod(modToFind)
			builder.WriteString("x")
			builder.WriteString(strconv.Itoa(propModCounter))
		} else {
			builder.WriteString(prop[1].GetValue())
		}

		if propCounter < len(props)-1 {
			builder.WriteString(" ")
		}
	}

	builder.WriteString(" ; ")
	builder.WriteString(comment)
	builder.WriteRune('\n')

	// Marshal all the mod's functions
	for _, function := range *module.GetMethods() {
		// Create a new stack for the function
		variables := stack.NewStack()
		variables.CreateScope()
		variables.Append(
			"this",
			wrapper2.NewFluentObject(
				module.BuildDummyWrapper(),
				module,
			),
			true,
		)

		function2.MarshalSingleFunction(
			ir,
			function,
			ir.GetFunction(function),
			fileCode,
			builder,
			counter,
			variables,
		)
	}

	builder.WriteString("endm\n")
}

// MarshalMods marshals all the  mods inside the given IrWrapper into the given strings.Builder
func MarshalMods(
	ir *wrapper.IrWrapper,
	builder *strings.Builder,
	counter *int,
	fileCode *ast.FileCode,
) {
	// Get all mods
	mods := ir.GetMods()

	// Iterate over all mods
	for module := range mods {
		// Ignore mods with generic templates, they're constructed dynamically
		if len(module.GetTemplates()) > 0 {
			continue
		}

		dummyWrapper := module.BuildDummyWrapper()
		MarshalSingleMod(module, ir, builder, counter, fileCode, dummyWrapper.Marshal())
	}
}
