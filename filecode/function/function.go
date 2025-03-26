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

package function

import (
	"fluent/ast"
	"fluent/filecode/trace"
	"fluent/filecode/types/wrapper"
)

// Function represents a function in the Fluent programming language.
type Function struct {
	// Name is the name of the function.
	Name string
	// Public indicates whether the function is public.
	Public bool
	// Params is a map of parameter names to their types.
	Params []Param
	// ReturnType is the type of the function's return value.
	ReturnType wrapper.TypeWrapper
	// Body is the abstract syntax tree (AST) of the function's body.
	Body ast.AST
	// Trace contains trace information for the function.
	Trace trace.Trace
	// Path contains the path of the file where the function is defined.
	Path string
	// Templates is a map of template names (generics).
	Templates map[string]bool
	// IsStd indicates whether the function is a standard library function.
	IsStd bool
}
