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
	ast2 "fluent/ast"
	"fluent/filecode/converter/function"
	"fluent/filecode/converter/redefinition"
	function2 "fluent/filecode/function"
	"fluent/filecode/module"
	"fluent/filecode/trace"
)

// ConvertModule converts an AST to a module.Module.
// It processes the AST to extract module properties such as
// public visibility, name, generics, and body.
//
// Parameters:
// - ast: A pointer to the AST to be converted.
// - contents: The contents of the file.
//
// Returns:
// - A module.Module with the extracted properties.
func ConvertModule(ast *ast2.AST, contents string) module.Module {
	// Used to know where to look for the function's name
	startAt := 0

	result := module.Module{
		Functions:    make(map[string]*function2.Function),
		Declarations: make(map[string]module.Declaration),
		Trace: trace.Trace{
			Line:   ast.Line,
			Column: ast.Column,
		},
		Path: *ast.File,
	}

	children := *ast.Children

	// Get the first child
	if children[0].Rule == ast2.Public {
		result.Public = true
		startAt = 1
	}

	// Get the module's name
	result.Name = *children[startAt].Value
	startAt += 1

	// Look for generics
	if children[startAt].Rule == ast2.Templates {
		// Append the generics
		for _, generic := range *children[startAt].Children {
			result.Templates[*generic.Value] = true
		}

		startAt += 1
	}

	// Iterate over the nodes of the block node to find functions and declarations
	for _, node := range *children[startAt].Children {
		rule := node.Rule

		switch rule {
		case ast2.Function:
			// Convert the function to a module.Function
			fn := function.ConvertFunction(node, false)

			// Check for redefinitions
			redefinition.CheckRedefinition(result.Functions, fn.Name, fn, contents, result.Path)

			// Append the function to the module
			result.Functions[fn.Name] = &fn
		case ast2.Declaration, ast2.IncompleteDeclaration:
			// Convert the declaration
			name, dec := ConvertDeclaration(node)

			// Store the declaration
			result.Declarations[name] = dec
		default:
		}
	}

	return result
}
