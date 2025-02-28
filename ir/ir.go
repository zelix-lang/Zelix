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

package ir

import (
	"fluent/filecode"
	"fluent/ir/pool"
	"fluent/ir/rule/function"
	"fluent/ir/tree"
	"fmt"
	"strings"
)

// BuildIr constructs the Intermediate Representation (IR) for the given file code.
// It processes the file code, appends imported modules and functions, and writes
// the functions to the IR.
//
// Parameters:
// - fileCode: The file code to build the IR for.
// - entry: A map of file codes representing the imported files.
// - fileId: The ID of the file being processed.
// - traceCounters: A pointer to a map of trace counters.
// - usedStrings: A pointer to a map of used strings.
//
// Returns:
// - A string representing the constructed IR.
func BuildIr(
	fileCode filecode.FileCode,
	entry map[string]filecode.FileCode,
	fileId int,
	isMain bool,
	traceCounters *pool.NumPool,
	usedStrings *pool.StringPool,
	usedNumbers *pool.StringPool,
	nameCounters *map[string]map[string]string,
	localCounters map[string]string,
) string {
	// Use a strings.Builder to properly handle the IR building
	builder := strings.Builder{}

	// Add trace instructions
	traceFileName := fmt.Sprintf("__trace_f%d__", fileId)
	builder.WriteString("ref ")
	builder.WriteString(traceFileName)
	builder.WriteString(" str ")
	builder.WriteString(fileCode.Path)
	builder.WriteString("\n")

	// Create a global InstructionTree for the file
	fileTree := tree.InstructionTree{
		Children: &[]*tree.InstructionTree{},
	}

	// Add all imported functions and modules
	for _, imp := range fileCode.Imports {
		impFile := entry[imp]

		for _, fun := range impFile.Functions {
			if fun.Public {
				fileCode.Functions[fun.Name] = fun
			}
		}

		for _, mod := range impFile.Modules {
			if mod.Public {
				fileCode.Modules[mod.Name] = mod
			}
		}
	}

	// Marshal all functions
	for _, fun := range fileCode.Functions {
		// Skip imported functions
		if fun.Path != fileCode.Path {
			continue
		}

		// Skip functions with generics
		if len(fun.Templates) > 0 {
			continue
		}

		function.MarshalFunction(
			fun,
			&fileCode,
			traceFileName,
			fileId,
			isMain,
			traceCounters,
			usedStrings,
			usedNumbers,
			&fileTree,
			nameCounters,
			localCounters,
		)
	}

	// Write the instructions to the builder
	for _, child := range *fileTree.Children {
		builder.WriteString(child.Representation.String())
		builder.WriteString("\n")
	}

	// Remove this FileCode's ID from the string pool
	usedStrings.RemoveId(fileId)
	return builder.String()
}
