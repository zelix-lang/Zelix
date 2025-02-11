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
	"fluent/filecode/module"
	trace2 "fluent/filecode/trace"
	"fluent/filecode/types"
)

// ConvertDeclaration converts an AST node representing a declaration
// into a module.Declaration. It returns the name of the declaration
// and the module.Declaration struct.
//
// Parameters:
// - node: A pointer to the AST node to be converted.
//
// Returns:
// - A string representing the name of the declaration.
// - A module.Declaration struct containing the details of the declaration.
func ConvertDeclaration(node *ast.AST) (string, module.Declaration) {
	children := *node.Children

	name := *children[1].Value
	constant := *children[0].Value == "const"
	varType := types.ConvertToTypeWrapper(*children[2])
	var valueNode *ast.AST

	// Check if the declaration is incomplete
	if len(children) == 4 {
		valueNode = children[3]
	}

	return name, module.Declaration{
		IsConstant:   constant,
		Value:        valueNode,
		Type:         varType,
		IsIncomplete: valueNode == nil,
		Trace: trace2.Trace{
			Line:   node.Line,
			Column: node.Column,
		},
	}
}
