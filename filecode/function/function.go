/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package function

import (
	"fluent/ast"
	"fluent/filecode/trace"
	"fluent/filecode/types"
)

// Function represents a function in the Fluent programming language.
type Function struct {
	// Name is the name of the function.
	Name string
	// Public indicates whether the function is public.
	Public bool
	// Params is a map of parameter names to their types.
	Params map[string]Param
	// ReturnType is the type of the function's return value.
	ReturnType types.TypeWrapper
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
