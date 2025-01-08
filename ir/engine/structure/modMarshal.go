package structure

import (
	"fluent/ast"
	"fluent/code/mod"
	"fluent/ir/wrapper"
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
) {
	// Build variables that have mods with generics
	props := module.GetVarDeclarations()

	for _, prop := range props {
		MarshalGenericsInVar(
			prop,
			ir,
			builder,
			counter,
			fileCode,
		)
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
	builder.WriteRune('\n')

	// Marshal properties
	for _, prop := range props {
		MarshalVariable(
			prop,
			ir,
			builder,
			counter,
			fileCode,
		)
	}

	// TODO!

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

		MarshalSingleMod(module, ir, builder, counter, fileCode)
	}
}
