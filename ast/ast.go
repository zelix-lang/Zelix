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

package ast

// AST represents an abstract syntax tree node.
type AST struct {
	// rule is the parser rule associated with this AST node.
	Rule Rule
	// children are the child nodes of this AST node.
	Children *[]*AST
	// value is the string value of this AST node.
	Value *string
	// line is the line number in the source file where this AST node is located.
	Line int
	// column is the column number in the source file where this AST node is located.
	Column int
	// file is the name of the source file where this AST node is located.
	File *string
}
