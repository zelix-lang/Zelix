/*
   The Fluent Programming Language
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

func (a AST) Marshal(spaces int) string {
	var result string
	result += a.Rule.String() + " "
	if a.Value != nil {
		result += *a.Value
	}
	if a.Children != nil {
		result += " {"
		for _, child := range *a.Children {
			result += "\n"
			for i := 0; i < spaces; i++ {
				result += " "
			}
			result += child.Marshal(spaces + 2)
		}
		result += "\n"
		for i := 0; i < spaces-2; i++ {
			result += " "
		}
		result += "}"
	}

	return result
}
