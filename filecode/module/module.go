/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package module

import (
	"fluent/ast"
	"fluent/filecode/function"
	"fluent/filecode/trace"
)

// Module represents a module in the Fluent programming language.
type Module struct {
	// Name is the name of the module.
	Name string
	// Public indicates whether the module is public.
	Public bool
	// Functions is a map of function names to their definitions.
	Functions map[string]function.Function
	// Declarations is a list of declarations in the module.
	Declarations []ast.AST
	// IncompleteDeclarations is a list of incomplete declarations in the module.
	IncompleteDeclarations []ast.AST
	// Generics is a list of generic types in the module.
	Generics []string
	// Trace contains trace information for the module.
	Trace trace.Trace
	// Path is the file path of the module.
	Path string
}
