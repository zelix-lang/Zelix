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

func MarshalModule(
	mod *module.Module,
	trace *filecode.FileCode,
	modulePropCounters *map[string]*util.OrderedMap[*string, *string],
	localCounters *map[string]string,
	fileTree *tree.InstructionTree,
	traceFileName string,
	fileCodeId int,
	traceCounters *pool.NumPool,
	usedStrings *pool.StringPool,
	usedNumbers *pool.StringPool,
	nameCounters *map[string]map[string]string,
) {
	// Create a new InstructionTree for the module
	modTree := tree.InstructionTree{
		Children:       &[]*tree.InstructionTree{},
		Representation: &strings.Builder{},
	}

	// Add the mod tree to the file tree
	*fileTree.Children = append(*fileTree.Children, &modTree)

	modTree.Representation.WriteString("mod ")
	modTree.Representation.WriteString((*localCounters)[mod.Name])
	modTree.Representation.WriteString(" ")

	// Get the prop counters
	propCounters := (*modulePropCounters)[mod.Name]

	// Iterate over all the module's properties
	propCounters.Iterate(func(name *string, computedName *string) bool {
		// Get the property
		prop, ok := mod.Declarations[*name]

		if ok {
			modTree.Representation.WriteString(prop.Type.Marshal())
			modTree.Representation.WriteString(" ")
			return false
		}

		// Otherwise, marshal the method as a function
		fun := mod.Functions[*name]
		function.MarshalFunction(
			fun,
			trace,
			(*localCounters)[mod.Name],
			true,
			traceFileName,
			fileCodeId,
			false,
			modulePropCounters,
			traceCounters,
			usedStrings,
			usedNumbers,
			fileTree,
			nameCounters,
			*computedName,
		)

		return false
	})
}
