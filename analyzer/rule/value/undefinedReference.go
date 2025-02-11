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

package value

import (
	error3 "fluent/analyzer/error"
	"fluent/filecode"
	"fluent/filecode/types"
)

// AnalyzeUndefinedReference checks if a given type wrapper has an undefined reference in the provided file code trace.
// If the type is not primitive and the module does not exist in the trace, it logs an undefined reference error,
// prints detailed information, and exits the program.
//
// Parameters:
//   - trace: A pointer to the FileCode structure containing the code trace information.
//   - wrapper: A TypeWrapper structure that holds information about the type being analyzed.
func AnalyzeUndefinedReference(
	trace *filecode.FileCode,
	wrapper types.TypeWrapper,
	generics *map[string]bool,
) error3.Error {
	if !wrapper.IsPrimitive {
		// Check for generics
		if _, ok := (*generics)[wrapper.BaseType]; ok {
			return error3.Error{}
		}

		// See if the mod exists in the trace
		if _, ok := trace.Modules[wrapper.BaseType]; !ok {
			return error3.Error{
				Line:       wrapper.Trace.Line,
				Column:     wrapper.Trace.Column,
				Code:       error3.UndefinedReference,
				Additional: []string{wrapper.BaseType},
			}
		}
	}

	return error3.Error{}
}
