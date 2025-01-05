package mod

import (
	"fluent/code/mod"
	"fluent/ir/wrapper"
	"strconv"
	"strings"
)

// MarshalMod marshals a FluentMod object into the given IR builder
func MarshalMod(ir *wrapper.IrWrapper, builder *strings.Builder) {
	// Get all mods
	mods := ir.GetMods()

	pendingMods := make(map[*mod.FluentMod]int)
	// Iterate over all mods
	for module, counter := range mods {
		// Ignore mods with generic templates, they're constructed dynamically
		if len(module.GetTemplates()) > 0 {
			continue
		}

		// Write the mod's name
		builder.WriteString("mod x")
		builder.WriteString(strconv.Itoa(counter))
		builder.WriteByte(' ')

		// Iterate over all the mod's properties
		props := module.GetVarDeclarations()

		for i, prop := range props {
			// 2nd token is the variable name
			ir.AddModProp(module, prop[1].GetValue(), i)
		}

		// Write a newline
		builder.WriteByte('\n')
	}
}
