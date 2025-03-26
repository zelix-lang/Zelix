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

package tree

// ModMarshalPair represents a pair of module marshaling information.
type ModMarshalPair struct {
	// Name is the name of the module.
	Name string
	// Parent is a pointer to the parent instruction tree.
	Parent *InstructionTree
	// CallConstructor indicates whether the constructor should be called.
	CallConstructor bool
	// IsParam indicates whether this is a parameter.
	IsParam bool
	// Counter is a counter for the module.
	Counter int
}
