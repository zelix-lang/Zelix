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
	"fluent/ir/rule/function"
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
// - traceCounter: A pointer to the current trace counter.
//
// Returns:
// - A string representing the constructed IR.
func BuildIr(
	fileCode filecode.FileCode,
	entry map[string]filecode.FileCode,
	fileId int,
	traceCounters *map[int]int,
	traceCounter *int,
) string {
	// Use a strings.Builder to properly handle the IR building
	builder := strings.Builder{}

	traceFileVarName := fmt.Sprintf("__fc_%d_trace_line", fileId)

	// Move the stack a string for the trace of this file
	builder.WriteString("ref ")
	builder.WriteString(traceFileVarName)
	builder.WriteString(" str ")
	builder.WriteString(fileCode.Path)
	builder.WriteString("\n")

	// Append all the imported modules and functions to the file
	for _, importPath := range fileCode.Imports {
		importedFile := entry[importPath]

		// Append the imported functions
		for _, fun := range importedFile.Functions {

			if !fun.Public {
				continue
			}

			fileCode.Functions[fun.Name] = fun
		}

		// Append the imported modules
		for _, mod := range importedFile.Modules {
			if !mod.Public {
				continue
			}

			fileCode.Modules[mod.Name] = mod
		}
	}

	// Write the functions
	for name, fun := range fileCode.Functions {
		if fun.Path != fileCode.Path {
			continue
		}

		// Functions with generics are not built right now
		if len(fun.Templates) > 0 {
			continue
		}

		function.MarshalFunction(
			&builder,
			name,
			fun,
			&fileCode,
			traceFileVarName,
			traceCounters,
			traceCounter,
		)
	}

	return builder.String()
}
