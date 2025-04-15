/*
   The Fluent Programming Language
   -----------------------------------------------------
   This code is released under the GNU GPL v3 license.
   For more information, please visit:
   https://www.gnu.org/licenses/gpl-3.0.html
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package module

import (
	"fluent/filecode"
	"fluent/filecode/module"
	"fluent/ir/pool"
	"fluent/ir/rule/function"
	"fluent/ir/tree"
	"fluent/util"
	"strings"
)

// MarshalModule marshals a module into an InstructionTree.
// Parameters:
// - mod: The module to be marshaled.
// - originalPath: The original path of the module.
// - trace: The file code trace.
// - modulePropCounters: A map of module property counters.
// - localCounters: A map of local counters.
// - fileTree: The instruction tree to which the module will be added.
// - traceFileName: The name of the trace file.
// - fileCodeId: The file code ID.
// - traceCounters: A pool of numeric counters for tracing.
// - usedStrings: A pool of used strings.
// - usedArrays: A pool of used arrays.
// - usedNumbers: A pool of used numbers.
func MarshalModule(
	mod *module.Module,
	originalPath *string,
	trace *filecode.FileCode,
	modulePropCounters *map[*module.Module]*util.OrderedMap[string, *string],
	localCounters *map[string]*string,
	fileTree *tree.InstructionTree,
	traceFileName string,
	fileCodeId int,
	traceCounters *pool.NumPool,
	usedStrings *pool.StringPool,
	usedArrays *pool.StringPool,
	usedNumbers *pool.StringPool,
) {
	// Create a new InstructionTree for the module
	modTree := tree.InstructionTree{
		Children:       &[]*tree.InstructionTree{},
		Representation: &strings.Builder{},
	}

	// Add the mod tree to the file tree
	*fileTree.Children = append(*fileTree.Children, &modTree)

	modTree.Representation.WriteString("mod ")
	modTree.Representation.WriteString(*(*localCounters)[mod.Name])
	modTree.Representation.WriteString(" ")

	// Get the prop counters
	propCounters := (*modulePropCounters)[mod]

	// Iterate over all the module's properties
	propCounters.Iterate(func(name string, computedName *string) bool {
		// Get the property
		prop, ok := mod.Declarations[name]

		if ok {
			modTree.Representation.WriteString("&")
			if prop.Type.IsPrimitive {
				modTree.Representation.WriteString(prop.Type.Marshal())
			} else {
				oldBaseType := prop.Type.BaseType
				prop.Type.BaseType = *(*localCounters)[oldBaseType]
				modTree.Representation.WriteString(prop.Type.Marshal())
				prop.Type.BaseType = oldBaseType
			}

			modTree.Representation.WriteString(" ")
			return false
		}

		// Otherwise, marshal the method as a function
		fun := mod.Functions[name]
		function.MarshalFunction(
			fun,
			trace,
			*(*localCounters)[mod.Name],
			true,
			traceFileName,
			fileCodeId,
			false,
			true,
			originalPath,
			modulePropCounters,
			traceCounters,
			usedStrings,
			usedArrays,
			usedNumbers,
			fileTree,
			computedName,
			localCounters,
		)

		return false
	})
}
