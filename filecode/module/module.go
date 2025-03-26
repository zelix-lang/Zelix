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
	"fluent/ast"
	"fluent/filecode/function"
	"fluent/filecode/trace"
	"fluent/filecode/types/wrapper"
)

// Declaration represents a declaration in the Fluent programming language.
type Declaration struct {
	// IsConstant indicates whether the declaration is a constant.
	IsConstant bool
	// Value is the AST node representing the value of the declaration.
	Value *ast.AST
	// Type is the type of the declaration.
	Type wrapper.TypeWrapper
	// IsIncomplete indicates whether the declaration is incomplete.
	IsIncomplete bool
	// Trace contains trace information for the declaration.
	Trace trace.Trace
}

// Module represents a module in the Fluent programming language.
type Module struct {
	// Name is the name of the module.
	Name string
	// Public indicates whether the module is public.
	Public bool
	// Functions is a map of function names to their definitions.
	Functions map[string]*function.Function
	// Declarations is a list of declarations in the module.
	Declarations map[string]Declaration
	// Templates is a list of generic types in the module.
	Templates map[string]bool
	// Trace contains trace information for the module.
	Trace trace.Trace
	// Path is the file path of the module.
	Path string
}
