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
	module2 "fluent/filecode/module"
	"fluent/ir/pool"
	"fluent/ir/rule/function"
	"fluent/ir/rule/module"
	"fluent/ir/tree"
	"fluent/util"
	"fmt"
	"strings"
)

// BuildIr constructs the Intermediate Representation (IR) for a given file.
// It processes the file's imports, functions, and modules, and marshals them
// into a string representation.
//
// Parameters:
// - fileCode: The FileCode object representing the file to be processed.
// - entry: A map of file names to FileCode objects for all files in the project.
// - fileId: An integer identifier for the file.
// - isMain: A boolean indicating if the file is the main entry point.
// - originalPath: A pointer to the original path of the file.
// - traceCounters: A pool of counters for tracing purposes.
// - traceStrings: A map of trace strings.
// - usedStrings: A pool of used strings.
// - usedArrays: A pool of used arrays.
// - usedNumbers: A pool of used numbers.
// - modulePropCounters: A map of module property counters.
// - localCounters: A map of local counters.
//
// Returns:
// - A string representing the IR of the file.
func BuildIr(
	fileCode filecode.FileCode,
	entry map[string]filecode.FileCode,
	fileId int,
	isMain bool,
	originalPath *string,
	traceCounters *pool.NumPool,
	traceStrings *map[string]*string,
	usedStrings *pool.StringPool,
	usedArrays *pool.StringPool,
	usedNumbers *pool.StringPool,
	modulePropCounters *map[*module2.Module]*util.OrderedMap[string, *string],
	localCounters map[string]*string,
) string {
	// Use a strings.Builder to properly handle the IR building
	builder := strings.Builder{}

	// Add trace instructions
	traceFileName := fmt.Sprintf("__trace_f%d__", fileId)
	(*traceStrings)[traceFileName] = &fileCode.Path

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

	// Marshal all modules
	for _, mod := range fileCode.Modules {
		// Skip imported modules
		if mod.Path != fileCode.Path {
			continue
		}

		// Skip functions with generics
		if len(mod.Templates) > 0 {
			continue
		}

		module.MarshalModule(
			mod,
			originalPath,
			&fileCode,
			modulePropCounters,
			&localCounters,
			&fileTree,
			traceFileName,
			fileId,
			traceCounters,
			usedStrings,
			usedArrays,
			usedNumbers,
		)
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
			"",
			false,
			traceFileName,
			fileId,
			isMain,
			false,
			originalPath,
			modulePropCounters,
			traceCounters,
			usedStrings,
			usedArrays,
			usedNumbers,
			&fileTree,
			localCounters[fun.Name],
			&localCounters,
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
